package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

// viewConfig is the immutable per-view input shared by every worker.
type viewConfig struct {
	Username      string
	CapmonsterKey string
	// DeviceNum is the numeric device enum guns.lol expects: desktop=0,
	// mobile=1, tablet=2.
	DeviceNum int
	Referrer  string
}

// reporter receives step transitions from runView so the same pipeline can
// drive either the single-run spinner or a parallel board row.
type reporter interface {
	// begin marks the start of a named step.
	begin(step string)
	// complete marks the current step done, with an optional trailing detail.
	complete(detail string)
}

// spinnerReporter renders one animated spinner per step for the single-run
// flow. On error the failing step's spinner is left running so fatalf can
// resolve it to a ✗.
type spinnerReporter struct {
	sp    *spinner
	label string
}

func (r *spinnerReporter) begin(step string) {
	r.sp = startSpinner(step)
	r.label = step
}

func (r *spinnerReporter) complete(detail string) {
	if r.sp != nil {
		r.sp.succeedDetail(r.label, detail)
		r.sp = nil
	}
}

// rowReporter maps step transitions onto a worker's row in the parallel board.
type rowReporter struct {
	b    *board
	slot int
}

func (r *rowReporter) begin(step string)      { r.b.setPhase(r.slot, step) }
func (r *rowReporter) complete(detail string) {}

// solveResult carries the two artifacts a submission needs plus their solve
// durations for reporting.
type solveResult struct {
	proof  string
	token  string
	powDur time.Duration
	tsDur  time.Duration
}

// solvePowAndTurnstile runs the WASM proof-of-work and the Turnstile solve
// concurrently. The PoW is pure computation and Turnstile is proxyless, so
// neither touches the session's HTTP client or cookies.
func (s *session) solvePowAndTurnstile(ctx context.Context, cfg viewConfig, wd *WorkerData, mod *powModule) (*solveResult, error) {
	type wasmOut struct {
		res *WasmResult
		err error
		dur time.Duration
	}
	type turnOut struct {
		token string
		err   error
		dur   time.Duration
	}

	wasmCh := make(chan wasmOut, 1)
	turnCh := make(chan turnOut, 1)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		start := time.Now()
		res, err := SolveWithWasm(ctx, mod, wd.Challenge, wd.Difficulty, strconv.FormatInt(wd.Timestamp, 10), wd.Nonce, wd.Seal)
		wasmCh <- wasmOut{res, err, time.Since(start)}
	}()

	go func() {
		defer wg.Done()
		start := time.Now()
		token, err := SolveTurnstile(ctx, cfg.CapmonsterKey, "https://guns.lol/"+cfg.Username, wd.Action, wd.CData)
		turnCh <- turnOut{token, err, time.Since(start)}
	}()

	wg.Wait()

	w := <-wasmCh
	if w.err != nil {
		return nil, fmt.Errorf("solve pow: %w", w.err)
	}
	t := <-turnCh
	if t.err != nil {
		return nil, fmt.Errorf("solve turnstile: %w", t.err)
	}
	return &solveResult{proof: w.res.Oo, token: t.token, powDur: w.dur, tsDur: t.dur}, nil
}

// runView drives one full record-a-view flow — fetch challenge, fetch the
// rotating PoW module, solve PoW + Turnstile concurrently, submit — reporting
// each step through r.
func (s *session) runView(ctx context.Context, cfg viewConfig, r reporter) error {
	r.begin("Fetching worker data")
	wd, err := s.FetchWorkerData(ctx, cfg.Username)
	if err != nil {
		return fmt.Errorf("fetch worker data: %w", err)
	}
	if wd == nil {
		return errors.New("no worker data returned")
	}
	r.complete(fmt.Sprintf("difficulty %d · challenge %s", wd.Difficulty, wd.ID))

	r.begin("Fetching PoW module")
	mod, err := s.FetchPowModule(ctx, wd.ID, wd.WorkerURL)
	if err != nil {
		return fmt.Errorf("fetch pow module: %w", err)
	}
	r.complete("")

	r.begin("Solving PoW & Turnstile")
	sr, err := s.solvePowAndTurnstile(ctx, cfg, wd, mod)
	if err != nil {
		return err
	}
	r.complete(fmt.Sprintf("PoW %s · Turnstile %s", sr.powDur.Round(time.Millisecond), sr.tsDur.Round(time.Millisecond)))

	r.begin("Submitting solution")
	err = s.SubmitSolution(ctx, SolutionPayload{
		Username:  cfg.Username,
		Version:   wd.Version,
		ID:        wd.ID,
		Timestamp: wd.Timestamp,
		Nonce:     wd.Nonce,
		Seal:      wd.Seal,
		Challenge: wd.Challenge,
		Proof:     sr.proof,
		Token:     sr.token,
		Referrer:  cfg.Referrer,
		DeviceNum: cfg.DeviceNum,
	}, cfg.CapmonsterKey)
	if err != nil {
		return fmt.Errorf("submit: %w", err)
	}
	r.complete("")
	return nil
}
