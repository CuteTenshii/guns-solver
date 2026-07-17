package main

import (
	"context"
	"encoding/json"
	"testing"
)

// TestSolveWithWasmKnownSample validates the migrated new-wasm ABI against the
// captured request.http sample: the next.js worker inputs must reproduce the
// exact _oo value seen in that request body.
func TestSolveWithWasmKnownSample(t *testing.T) {
	const (
		o09      = "7f9b9373edec946e8db9cba3759707c19c0e859e4c8398b9a72f0eb72f33fd32"
		orgTs    = "1784225895"
		nonce    = "nfuZwLGN02rjt3Ve9AdsWjcJvijTyzNZ"
		twoXa    = "skAFNBcCOj4DBAECAKgmpGgLJGGIOGQwZDNlZWUzNzQ5ODNhN2MxODk3NGExZTRkMTY4NTFiY2IwYWUyYTY3Y2RlYTdjMGM3MTNkYmI4NTMuuVQK0dDgEA"
		expectOo = "YgU5OTM1ZQAAAAAmST5X7i4dUg"
	)

	res, err := SolveWithWasm(context.Background(), o09, 5, orgTs, nonce, twoXa)
	if err != nil {
		t.Fatalf("SolveWithWasm: %v", err)
	}
	t.Logf("_oo=%s seal=%s", res.Oo, res.Seal)
	if res.Oo != expectOo {
		t.Fatalf("_oo mismatch:\n got  %q\n want %q", res.Oo, expectOo)
	}
}

// TestGppChallengeMatchesSample locks the "_gpp_ch" object (obfuscated keys,
// field order, value types) to the captured request.http POST body.
func TestGppChallengeMatchesSample(t *testing.T) {
	ch := gppChallenge{
		Version:   gppFormatVersion,
		Timestamp: 1784225895,
		O09:       "7f9b9373edec946e8db9cba3759707c19c0e859e4c8398b9a72f0eb72f33fd32",
		Nonce:     "nfuZwLGN02rjt3Ve9AdsWjcJvijTyzNZ",
		Salt:      "skAFNBcCOj4DBAECAKgmpGgLJGGIOGQwZDNlZWUzNzQ5ODNhN2MxODk3NGExZTRkMTY4NTFiY2IwYWUyYTY3Y2RlYTdjMGM3MTNkYmI4NTMuuVQK0dDgEA",
		Oo:        "YgU5OTM1ZQAAAAAmST5X7i4dUg",
	}
	const want = `{"_c14f":2,"_e92a":1784225895,"_b73d":"7f9b9373edec946e8db9cba3759707c19c0e859e4c8398b9a72f0eb72f33fd32","_a18c":"nfuZwLGN02rjt3Ve9AdsWjcJvijTyzNZ","_f52b":"skAFNBcCOj4DBAECAKgmpGgLJGGIOGQwZDNlZWUzNzQ5ODNhN2MxODk3NGExZTRkMTY4NTFiY2IwYWUyYTY3Y2RlYTdjMGM3MTNkYmI4NTMuuVQK0dDgEA","_7bc1":"YgU5OTM1ZQAAAAAmST5X7i4dUg"}`

	got, err := json.Marshal(ch)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(got) != want {
		t.Fatalf("_gpp_ch mismatch:\n got  %s\n want %s", got, want)
	}
}
