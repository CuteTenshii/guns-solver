package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const gunsOrigin = "https://guns.lol"

// powCacheRoot is where fetched view-PoW modules are cached on disk, keyed by
// challenge id. It is a var so tests can redirect it.
var powCacheRoot = filepath.Join(os.TempDir(), "guns-solver-pow")

// cacheKeyRe restricts a challenge id to filesystem-safe characters before it
// is used as a cache filename (the id is server-supplied).
var cacheKeyRe = regexp.MustCompile(`^[A-Za-z0-9_-]{1,128}$`)

// powModule is a fetched, rotation-specific PoW WASM together with the export
// names its glue drives. guns.lol rotates the challenge worker, glue, binary,
// and the hashed export symbols periodically, so these must be fetched per
// challenge rather than embedded: a stale binary validates the seal against an
// old key and the constructor rejects it with error code 1003.
type powModule struct {
	wasm        []byte
	ctorExport  string // GunsSolver constructor (9-arg wasm export)
	solveExport string // solve method (2-arg wasm export taking this.__wbg_ptr)
}

var (
	// xorCallRe matches the worker's `_helper([n,n,...],key)` de-obfuscation
	// calls, capturing the byte array and the XOR key.
	xorCallRe = regexp.MustCompile(`\[([0-9,]+)\],(\d+)\)`)
	// glueWasmURLRe extracts the wasm filename from the glue's init function.
	glueWasmURLRe = regexp.MustCompile(`new URL\("([^"]+)"`)
	// glueWasmCallRe matches every `wasm._export(args)` call in the glue.
	glueWasmCallRe = regexp.MustCompile(`wasm\.(_[0-9a-f]+)\(([^)]*)\)`)
)

// FetchPowModule resolves the current PoW WASM for a challenge, using an
// on-disk cache keyed by the challenge id. The module (binary and its export
// names) is immutable under a given id, and the id only changes when guns.lol
// rotates the challenge, so a cache hit is always current. On a miss it fetches
// via fetchPowModule and populates the cache best-effort.
func (s *session) FetchPowModule(ctx context.Context, id, workerPath string) (*powModule, error) {
	if m, ok := loadCachedModule(id); ok {
		return m, nil
	}
	m, err := s.fetchPowModule(ctx, workerPath)
	if err != nil {
		return nil, err
	}
	storeCachedModule(id, m)
	return m, nil
}

// fetchPowModule downloads the PoW module for a challenge. workerPath is the
// `_gpp_ch.u` value (e.g. /_challenge/pow/<id>/<worker>); the worker, glue, and
// binary all live in that directory. It fetches the worker, de-obfuscates the
// glue module name, parses the glue for the wasm URL and the rotated
// constructor/solve export symbols, then fetches the binary.
func (s *session) fetchPowModule(ctx context.Context, workerPath string) (*powModule, error) {
	slash := strings.LastIndex(workerPath, "/")
	if slash < 0 {
		return nil, fmt.Errorf("invalid worker path %q", workerPath)
	}
	baseDir := workerPath[:slash] // /_challenge/pow/<id>

	workerJS, err := s.powGet(ctx, workerPath)
	if err != nil {
		return nil, fmt.Errorf("fetch worker: %w", err)
	}
	gluePath := decodeGlueImport(string(workerJS))
	if gluePath == "" {
		return nil, errors.New("could not locate glue import in worker")
	}
	// gluePath is like "./_abc"; resolve it against the challenge directory.
	glueURL := baseDir + "/" + strings.TrimPrefix(gluePath, "./")

	glueJS, err := s.powGet(ctx, glueURL)
	if err != nil {
		return nil, fmt.Errorf("fetch glue: %w", err)
	}
	wasmName, ctorExport, solveExport, err := parseGlue(string(glueJS))
	if err != nil {
		return nil, err
	}

	wasmBytes, err := s.powGet(ctx, baseDir+"/"+wasmName)
	if err != nil {
		return nil, fmt.Errorf("fetch wasm: %w", err)
	}
	return &powModule{wasm: wasmBytes, ctorExport: ctorExport, solveExport: solveExport}, nil
}

// decodeGlueImport reverses the worker's XOR obfuscation to recover the glue
// module path — the decoded string that starts with "./".
func decodeGlueImport(workerJS string) string {
	for _, m := range xorCallRe.FindAllStringSubmatch(workerJS, -1) {
		key, err := strconv.Atoi(m[2])
		if err != nil {
			continue
		}
		parts := strings.Split(m[1], ",")
		b := make([]byte, 0, len(parts))
		ok := true
		for _, p := range parts {
			n, err := strconv.Atoi(p)
			if err != nil {
				ok = false
				break
			}
			b = append(b, byte(n^key))
		}
		if ok {
			if s := string(b); strings.HasPrefix(s, "./") {
				return s
			}
		}
	}
	return ""
}

// parseGlue extracts the wasm filename and the rotated constructor/solve export
// symbols from a wasm-bindgen glue module. The constructor is the wasm export
// called with nine arguments; solve is the one whose second argument is
// this.__wbg_ptr (distinguishing it from the 2-arg free export).
func parseGlue(glueJS string) (wasmName, ctorExport, solveExport string, err error) {
	u := glueWasmURLRe.FindStringSubmatch(glueJS)
	if u == nil {
		return "", "", "", errors.New("glue: wasm URL not found")
	}
	wasmName = u[1]

	for _, m := range glueWasmCallRe.FindAllStringSubmatch(glueJS, -1) {
		args := strings.Split(m[2], ",")
		switch {
		case len(args) == 9:
			ctorExport = m[1]
		case len(args) == 2 && strings.Contains(args[1], "__wbg_ptr"):
			solveExport = m[1]
		}
	}
	if ctorExport == "" || solveExport == "" {
		return "", "", "", fmt.Errorf("glue: could not resolve exports (ctor=%q solve=%q)", ctorExport, solveExport)
	}
	return wasmName, ctorExport, solveExport, nil
}

// powMeta is the sidecar stored next to a cached wasm binary, recording the
// rotated export names parsed from that rotation's glue.
type powMeta struct {
	Ctor  string `json:"ctor"`
	Solve string `json:"solve"`
}

// loadCachedModule returns the cached module for id, if present and complete.
func loadCachedModule(id string) (*powModule, bool) {
	if !cacheKeyRe.MatchString(id) {
		return nil, false
	}
	wasm, err := os.ReadFile(filepath.Join(powCacheRoot, id+".wasm"))
	if err != nil {
		return nil, false
	}
	metaBytes, err := os.ReadFile(filepath.Join(powCacheRoot, id+".json"))
	if err != nil {
		return nil, false
	}
	var meta powMeta
	if err := json.Unmarshal(metaBytes, &meta); err != nil || meta.Ctor == "" || meta.Solve == "" {
		return nil, false
	}
	return &powModule{wasm: wasm, ctorExport: meta.Ctor, solveExport: meta.Solve}, true
}

// storeCachedModule persists m under id, best-effort (cache failures are not
// fatal). The wasm is written atomically so a crash cannot leave a truncated
// binary that would later fail to instantiate.
func storeCachedModule(id string, m *powModule) {
	if !cacheKeyRe.MatchString(id) {
		return
	}
	if err := os.MkdirAll(powCacheRoot, 0700); err != nil {
		return
	}
	wasmPath := filepath.Join(powCacheRoot, id+".wasm")
	// A unique temp suffix keeps concurrent writers (parallel workers racing to
	// cache the same rotation) from clobbering each other's in-progress file.
	tmp := wasmPath + "." + randomSessionID() + ".tmp"
	if err := os.WriteFile(tmp, m.wasm, 0600); err != nil {
		return
	}
	if err := os.Rename(tmp, wasmPath); err != nil {
		os.Remove(tmp)
		return
	}
	metaBytes, err := json.Marshal(powMeta{Ctor: m.ctorExport, Solve: m.solveExport})
	if err != nil {
		return
	}
	os.WriteFile(filepath.Join(powCacheRoot, id+".json"), metaBytes, 0600)
}

// powGet fetches a challenge asset (worker/glue/wasm) from guns.lol with the
// session's client, User-Agent, and clearance cookies. path is an absolute
// origin-relative path beginning with "/".
func (s *session) powGet(ctx context.Context, path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", gunsOrigin+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	s.addClearanceCookies(req)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	s.captureCfClearance(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: status %d", path, resp.StatusCode)
	}
	return body, nil
}
