package main

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"
)

// TestLivePowSolve runs the full fetch-current-module-and-solve path against a
// live guns.lol challenge. It is network-gated: set GUNS_LIVE=1 to run. This is
// the regression test for the rotation 1003 bug — it confirms the fetched
// module matches the current challenge's seal.
func TestLivePowSolve(t *testing.T) {
	if os.Getenv("GUNS_LIVE") != "1" {
		t.Skip("set GUNS_LIVE=1 to run the live PoW integration test")
	}
	username := os.Getenv("GUNS_USER")
	if username == "" {
		username = "tenshii"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	wd, err := FetchWorkerData(ctx, username)
	if err != nil {
		t.Fatalf("FetchWorkerData: %v", err)
	}
	t.Logf("challenge id=%s worker=%s difficulty=%d", wd.ID, wd.WorkerURL, wd.Difficulty)

	mod, err := FetchPowModule(ctx, wd.ID, wd.WorkerURL)
	if err != nil {
		t.Fatalf("FetchPowModule: %v", err)
	}
	t.Logf("ctor=%s solve=%s wasm=%d bytes", mod.ctorExport, mod.solveExport, len(mod.wasm))

	res, err := SolveWithWasm(ctx, mod, wd.Challenge, wd.Difficulty,
		strconv.FormatInt(wd.Timestamp, 10), wd.Nonce, wd.Seal)
	if err != nil {
		t.Fatalf("SolveWithWasm: %v", err)
	}
	t.Logf("proof=%s", res.Oo)
	if res.Oo == "" {
		t.Fatal("solver returned an empty _oo proof")
	}
}
