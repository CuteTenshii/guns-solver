package main

import (
	"context"
	"strings"
	"testing"
)

// TestResolveProxySession checks the {session} placeholder is replaced with a
// fresh token each call, and that URLs without it are left untouched.
func TestResolveProxySession(t *testing.T) {
	const tmpl = "http://user:pass_session-{session}_lifetime-30s@geo.iproyal.com:12345"
	a := resolveProxySession(tmpl)
	b := resolveProxySession(tmpl)
	if strings.Contains(a, proxySessionPlaceholder) {
		t.Fatalf("placeholder not substituted: %q", a)
	}
	if a == b {
		t.Fatal("expected a distinct session token per call")
	}

	const plain = "http://user:pass@host:1234"
	if got := resolveProxySession(plain); got != plain {
		t.Fatalf("proxy without placeholder was altered: %q", got)
	}
}

// TestSolveWithWasmKnownSample validates the wazero ABI against the embedded
// clearance binary with a captured known-answer sample: the same inputs must
// reproduce the exact _oo. This also exercises the powModule plumbing that the
// fetched view-PoW modules share.
func TestSolveWithWasmKnownSample(t *testing.T) {
	const (
		o09      = "7f9b9373edec946e8db9cba3759707c19c0e859e4c8398b9a72f0eb72f33fd32"
		orgTs    = "1784225895"
		nonce    = "nfuZwLGN02rjt3Ve9AdsWjcJvijTyzNZ"
		twoXa    = "skAFNBcCOj4DBAECAKgmpGgLJGGIOGQwZDNlZWUzNzQ5ODNhN2MxODk3NGExZTRkMTY4NTFiY2IwYWUyYTY3Y2RlYTdjMGM3MTNkYmI4NTMuuVQK0dDgEA"
		expectOo = "YgU5OTM1ZQAAAAAmST5X7i4dUg"
	)

	res, err := SolveWithWasm(context.Background(), clearancePowModule(), o09, 5, orgTs, nonce, twoXa)
	if err != nil {
		t.Fatalf("SolveWithWasm: %v", err)
	}
	t.Logf("_oo=%s seal=%s", res.Oo, res.Seal)
	if res.Oo != expectOo {
		t.Fatalf("_oo mismatch:\n got  %q\n want %q", res.Oo, expectOo)
	}
}

// TestDecodeGlueImport verifies the worker's XOR obfuscation is reversed to the
// glue module path (the decoded string starting with "./"), ignoring decoys.
func TestDecodeGlueImport(t *testing.T) {
	// [106,...]^5 = "onmessage" (decoy), [43,...]^5 = "./_abc" (glue path).
	const worker = `self[_h([106,107,104,96,118,118,100,98,96],5)]=async function(e){const n=await import(_h([43,42,90,100,103,102],5));}`
	if got := decodeGlueImport(worker); got != "./_abc" {
		t.Fatalf("decodeGlueImport = %q, want %q", got, "./_abc")
	}
}

// TestParseGlue locks the extraction of the wasm filename and the rotated
// constructor (9-arg) and solve (this.__wbg_ptr) export symbols from glue JS.
func TestParseGlue(t *testing.T) {
	const glue = `class X{constructor(e,t,n,s,o){const m=wasm._aaa111(i,a,t,r,c,l,d,u,h);return this.__wbg_ptr=m}_bbb222(){wasm.__wbindgen_add_to_stack_pointer(-16);wasm._ccc333(s,this.__wbg_ptr)}free(){wasm._ddd444(e,0)}}` +
		`function __wbg_init(e){e===void 0&&(e=new URL("_eee555",import.meta.url))}`

	wasmName, ctor, solve, err := parseGlue(glue)
	if err != nil {
		t.Fatalf("parseGlue: %v", err)
	}
	if wasmName != "_eee555" || ctor != "_aaa111" || solve != "_ccc333" {
		t.Fatalf("parseGlue = (%q, %q, %q), want (_eee555, _aaa111, _ccc333)", wasmName, ctor, solve)
	}
}

// TestPowModuleCache checks the disk cache round-trips a module keyed by
// challenge id and refuses filesystem-unsafe keys.
func TestPowModuleCache(t *testing.T) {
	old := powCacheRoot
	powCacheRoot = t.TempDir()
	defer func() { powCacheRoot = old }()

	if _, ok := loadCachedModule("kH9YIpdB_Pj_2hyVg9PdEQ"); ok {
		t.Fatal("unexpected cache hit on empty cache")
	}

	want := &powModule{wasm: []byte("\x00asm\x01\x00\x00\x00"), ctorExport: "_f1fdb45177b3e", solveExport: "_891cc7bafaa129cd766"}
	storeCachedModule("kH9YIpdB_Pj_2hyVg9PdEQ", want)

	got, ok := loadCachedModule("kH9YIpdB_Pj_2hyVg9PdEQ")
	if !ok {
		t.Fatal("expected cache hit after store")
	}
	if string(got.wasm) != string(want.wasm) || got.ctorExport != want.ctorExport || got.solveExport != want.solveExport {
		t.Fatalf("cache mismatch: %+v", got)
	}

	// A path-traversal key must never be stored or loaded.
	storeCachedModule("../evil", want)
	if _, ok := loadCachedModule("../evil"); ok {
		t.Fatal("unsafe cache key should never hit")
	}
}

// TestBuildViewPayloadShape locks the positional array of the rewritten view
// beacon: [token, [v, e, t, n, s, c, proof], username, deviceNum, referrer].
// Field order and nesting are load-bearing, so they are asserted verbatim.
func TestBuildViewPayloadShape(t *testing.T) {
	p := SolutionPayload{
		Username:  "tenshii",
		Version:   3,
		ID:        "pUYewMGcOW2C5cSCBaiUiQ",
		Timestamp: 1784905110,
		Nonce:     "ms1Wd6PfYqbA9ntoS7rjEceSvNIhgyuN",
		Seal:      "akAFEQssIioDAgABBLyE9lx48jEUZjE0ZmE4OGYzMTg5ZTQ2ODY2OWMzOGU3MjQ4YzhiYTgyZTA5YTQ2NDJmOTM1YmYxNTMyODhkZWIyZjHpDQAbc8kfcw",
		Challenge: "62e9a3fa26732e6b8de941ed4eedbddebb6990a328e077d236999e6a3c276c50",
		Proof:     "YgU5OTM1ZQAAAAAmST5X7i4dUg",
		Token:     "TT.dummy-turnstile-token",
		Referrer:  "",
		DeviceNum: 0,
	}
	const want = `["TT.dummy-turnstile-token",[3,"pUYewMGcOW2C5cSCBaiUiQ",1784905110,"ms1Wd6PfYqbA9ntoS7rjEceSvNIhgyuN","akAFEQssIioDAgABBLyE9lx48jEUZjE0ZmE4OGYzMTg5ZTQ2ODY2OWMzOGU3MjQ4YzhiYTgyZTA5YTQ2NDJmOTM1YmYxNTMyODhkZWIyZjHpDQAbc8kfcw","62e9a3fa26732e6b8de941ed4eedbddebb6990a328e077d236999e6a3c276c50","YgU5OTM1ZQAAAAAmST5X7i4dUg"],"tenshii",0,""]`

	got, err := buildViewPayload(p)
	if err != nil {
		t.Fatalf("buildViewPayload: %v", err)
	}
	if string(got) != want {
		t.Fatalf("view payload mismatch:\n got  %s\n want %s", got, want)
	}
}

// TestParseGppChallenge verifies the _gpp_ch object is recovered from a
// profile page's escaped Next.js flight payload.
func TestParseGppChallenge(t *testing.T) {
	const body = `self.__next_f.push([1,"5:[...\"uid\":1617758,\"_gpp_ch\":{\"v\":3,\"e\":\"pUYewMGcOW2C5cSCBaiUiQ\",\"u\":\"/_challenge/pow/pUYewMGcOW2C5cSCBaiUiQ/_a9dee1fa9e512fc3736ee20\",\"t\":1784905110,\"n\":\"ms1Wd6PfYqbA9ntoS7rjEceSvNIhgyuN\",\"s\":\"akAFEQ\",\"c\":\"62e9a3fa26732e6b8de941ed4eedbddebb6990a328e077d236999e6a3c276c50\",\"d\":5,\"cd\":\"62e9a3fa26732e6b8de941ed4eedbddebb6990a328e077d236999e6a3c276c50\",\"a\":\"guns_view\"},\"success\":true}]"])`

	wd, err := parseGppChallenge(body)
	if err != nil {
		t.Fatalf("parseGppChallenge: %v", err)
	}
	if wd.Version != 3 || wd.ID != "pUYewMGcOW2C5cSCBaiUiQ" || wd.Timestamp != 1784905110 ||
		wd.Nonce != "ms1Wd6PfYqbA9ntoS7rjEceSvNIhgyuN" || wd.Difficulty != 5 || wd.Action != "guns_view" {
		t.Fatalf("unexpected challenge: %+v", wd)
	}
	if wd.WorkerURL != "/_challenge/pow/pUYewMGcOW2C5cSCBaiUiQ/_a9dee1fa9e512fc3736ee20" {
		t.Fatalf("unexpected worker URL: %q", wd.WorkerURL)
	}
}
