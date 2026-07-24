package main

import (
	"bytes"
	"context"
	crand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// The guns_clearance and cf_clearance cookies are cached between runs in the
// same temp dir as the fetched PoW modules (powCacheRoot).
func clearanceFile() string   { return filepath.Join(powCacheRoot, "clearance.txt") }
func cfClearanceFile() string { return filepath.Join(powCacheRoot, "cf_clearance.txt") }

func init() {
	if data, err := os.ReadFile(clearanceFile()); err == nil {
		gunsClearance = strings.TrimSpace(string(data))
	}
	if data, err := os.ReadFile(cfClearanceFile()); err == nil {
		cfClearance = strings.TrimSpace(string(data))
	}
}

func saveClearance() {
	if err := os.MkdirAll(powCacheRoot, 0700); err == nil {
		os.WriteFile(clearanceFile(), []byte(gunsClearance), 0600)
	}
}

func saveCfClearance() {
	if err := os.MkdirAll(powCacheRoot, 0700); err == nil {
		os.WriteFile(cfClearanceFile(), []byte(cfClearance), 0600)
	}
}

// addClearanceCookies attaches the guns.lol and Cloudflare clearance cookies to
// req when they are known. cf_clearance is required by the Cloudflare Managed
// Challenge in front of guns.lol; without it the API returns 403.
func addClearanceCookies(req *http.Request) {
	if gunsClearance != "" {
		req.AddCookie(&http.Cookie{Name: "guns_clearance", Value: gunsClearance})
	}
	if cfClearance != "" {
		req.AddCookie(&http.Cookie{Name: "cf_clearance", Value: cfClearance})
	}
}

// captureCfClearance persists a cf_clearance cookie if the response sets one, so
// a rotated value survives across requests and runs.
func captureCfClearance(resp *http.Response) {
	for _, c := range resp.Cookies() {
		if c.Name == "cf_clearance" && c.Value != "" && c.Value != cfClearance {
			cfClearance = c.Value
			saveCfClearance()
		}
	}
}

var (
	gunsClearance = ""
	cfClearance   = ""
	// proxyURL is the upstream proxy every request egresses through; it is also
	// passed to CapMonster when minting a cf_clearance cookie so the cookie's
	// bound IP matches the tool's.
	proxyURL   = ""
	httpClient = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	// gppChRegex captures the `_gpp_ch` challenge object out of the profile
	// page's Next.js flight payload. Interior quotes are backslash-escaped there
	// (the object lives inside a JS string literal), so the captured group is
	// unescaped before it parses as JSON. The object is flat, so a brace-free
	// body ([^{}]*) matches its full extent.
	gppChRegex = regexp.MustCompile(`\\"_gpp_ch\\":(\{[^{}]*\})`)
	// flightUnescaper turns the escaped flight-payload substring back into JSON.
	flightUnescaper = strings.NewReplacer(`\"`, `"`, `\\`, `\`, `\/`, `/`)
	// Challenge page regexes (double-quoted JS object fields in _gs_sets), used
	// by the legacy guns_clearance interstitial.
	challengeNonceRegex = regexp.MustCompile(`[{,]_n:"([^"]+)"`)
	challengeO09Regex   = regexp.MustCompile(`[{,]o09:"([^"]+)"`)
	challenge2xaRegex   = regexp.MustCompile(`[{,]_2xa:"([^"]+)"`)
	challengeOrgTsRegex = regexp.MustCompile(`[{,]_org_ts:"([^"]+)"`)
	challengeDRegex     = regexp.MustCompile(`,d:"([^"]+)"`)
	challengeSRegex     = regexp.MustCompile(`[{,]__s:"([^"]+)"`)
)

// proxySessionPlaceholder, when present in the proxy URL, is replaced with a
// fresh random token each run. Sticky-session proxy providers key the session
// (and thus the exit IP) off the credentials — e.g. IPRoyal's
// `password_session-<id>_lifetime-30s` — so substituting a new token per
// execution hands out a new IP, which is what botting views requires.
const proxySessionPlaceholder = "{session}"

// SetProxy routes every request (and CapMonster's egress, via proxyURL) through
// rawURL. A {session} placeholder is resolved once here so the tool and
// CapMonster share the same sticky session for the whole run.
func SetProxy(rawURL string) error {
	resolved := resolveProxySession(rawURL)
	u, err := url.Parse(resolved)
	if err != nil {
		return err
	}
	httpClient.Transport = &http.Transport{
		Proxy: http.ProxyURL(u),
	}
	proxyURL = resolved
	return nil
}

// resolveProxySession replaces every proxySessionPlaceholder in rawURL with one
// freshly generated session token, leaving URLs without the placeholder
// unchanged.
func resolveProxySession(rawURL string) string {
	if !strings.Contains(rawURL, proxySessionPlaceholder) {
		return rawURL
	}
	return strings.ReplaceAll(rawURL, proxySessionPlaceholder, randomSessionID())
}

// randomSessionID returns a 16-hex-character token for use as a proxy session
// id.
func randomSessionID() string {
	b := make([]byte, 8)
	if _, err := crand.Read(b); err != nil {
		panic(fmt.Sprintf("generate proxy session id: %v", err))
	}
	return hex.EncodeToString(b)
}

// WorkerData is the guns.lol proof-of-work challenge scraped from a profile
// page's Next.js flight payload (the `_gpp_ch` object). The obfuscated single-
// letter page keys are expanded here; each field comment names its on-page key.
type WorkerData struct {
	Version    int    // v — payload schema version
	ID         string // e — challenge id (first path segment of WorkerURL)
	WorkerURL  string // u — same-origin path of the PoW worker module
	Timestamp  int64  // t — challenge issue time (unix seconds)
	Nonce      string // n — per-challenge nonce
	Seal       string // s — opaque server seal, replayed in the solution
	Challenge  string // c — 64-char hex challenge input to the solver
	Difficulty int    // d — required proof difficulty
	CData      string // cd — Turnstile data-cdata
	Action     string // a — Turnstile data-action
}

// gppChallengeData mirrors the on-page `_gpp_ch` JSON with its raw short keys.
type gppChallengeData struct {
	V  int    `json:"v"`
	E  string `json:"e"`
	U  string `json:"u"`
	T  int64  `json:"t"`
	N  string `json:"n"`
	S  string `json:"s"`
	C  string `json:"c"`
	D  int    `json:"d"`
	Cd string `json:"cd"`
	A  string `json:"a"`
}

// parseGppChallenge extracts and decodes the `_gpp_ch` challenge object from a
// profile page body.
func parseGppChallenge(body string) (*WorkerData, error) {
	m := gppChRegex.FindStringSubmatch(body)
	if m == nil {
		return nil, errors.New("failed to locate _gpp_ch challenge in page")
	}
	var ch gppChallengeData
	if err := json.Unmarshal([]byte(flightUnescaper.Replace(m[1])), &ch); err != nil {
		return nil, fmt.Errorf("decode _gpp_ch: %w", err)
	}
	return &WorkerData{
		Version:    ch.V,
		ID:         ch.E,
		WorkerURL:  ch.U,
		Timestamp:  ch.T,
		Nonce:      ch.N,
		Seal:       ch.S,
		Challenge:  ch.C,
		Difficulty: ch.D,
		CData:      ch.Cd,
		Action:     ch.A,
	}, nil
}

func FetchWorkerData(ctx context.Context, username string) (*WorkerData, error) {
	req, err := http.NewRequest("GET", "https://guns.lol/"+username, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	req.AddCookie(&http.Cookie{Name: "GUNS_LOCALE", Value: "en"})
	req.AddCookie(&http.Cookie{Name: "GUNS_PATH_LOCALE", Value: "en"})
	addClearanceCookies(req)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	captureCfClearance(resp)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	body := string(bodyBytes)

	if resp.StatusCode == http.StatusUnauthorized {
		if gunsClearance != "" {
			gunsClearance = ""
			os.Remove(clearanceFile())
		}
		warnln("Got 401 — solving guns_clearance interstitial")
		if err = solveChallenge(ctx, body); err != nil {
			return nil, fmt.Errorf("challenge: %w", err)
		}
		return FetchWorkerData(ctx, username)
	}

	if resp.StatusCode == http.StatusTemporaryRedirect {
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "guns_clearance" {
				gunsClearance = cookie.Value
				saveClearance()
				return FetchWorkerData(ctx, username)
			}
		}
		if gunsClearance != "" {
			location := resp.Header.Get("Location")
			return FetchWorkerData(ctx, location[1:]) // remove leading slash
		}
		return nil, errors.New("307 redirect without clearance cookie")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d %s", resp.StatusCode, resp.Status)
	}

	return parseGppChallenge(body)
}

// solveChallenge handles the guns_clearance 401 interstitial (solve a PoW →
// POST /_challenge/verify → guns_clearance cookie). This is a separate,
// stable subsystem from the rotating view PoW: it uses the `_gs_sets` challenge
// format and the embedded clearance binary.
func solveChallenge(ctx context.Context, body string) error {
	nonce := challengeNonceRegex.FindStringSubmatch(body)
	o09 := challengeO09Regex.FindStringSubmatch(body)
	twoXa := challenge2xaRegex.FindStringSubmatch(body)
	orgTs := challengeOrgTsRegex.FindStringSubmatch(body)
	d := challengeDRegex.FindStringSubmatch(body)
	s := challengeSRegex.FindStringSubmatch(body)

	if nonce == nil || o09 == nil || twoXa == nil || orgTs == nil || d == nil || s == nil {
		return errors.New("failed to parse challenge data from page")
	}

	difficulty, err := strconv.Atoi(d[1])
	if err != nil {
		return fmt.Errorf("invalid difficulty %q: %w", d[1], err)
	}

	infoln("clearance params: difficulty %d nonce %s", difficulty, nonce[1])
	res, err := SolveWithWasm(ctx, clearancePowModule(), o09[1], difficulty, orgTs[1], nonce[1], twoXa[1])
	if err != nil {
		return fmt.Errorf("wasm solve: %w", err)
	}
	infoln("clearance solved: _oo %s", truncateMiddle(res.Oo, 24))

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("_o", res.Oo)
	w.WriteField("_s", s[1])
	w.WriteField("_u", nonce[1])
	w.WriteField("_i", twoXa[1])
	w.WriteField("_x", o09[1])
	w.WriteField("_t", orgTs[1])
	w.Close()

	req, err := http.NewRequest("POST", "https://guns.lol/_challenge/verify", &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("User-Agent", userAgent)
	addClearanceCookies(req)

	vresp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer vresp.Body.Close()
	captureCfClearance(vresp)
	io.ReadAll(vresp.Body) // drain

	if vresp.StatusCode != http.StatusOK {
		return fmt.Errorf("verify returned %d", vresp.StatusCode)
	}

	for _, cookie := range vresp.Cookies() {
		if cookie.Name == "guns_clearance" {
			gunsClearance = cookie.Value
			saveClearance()
			return nil
		}
	}
	return errors.New("no clearance cookie in verify response")
}
