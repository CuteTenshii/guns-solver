package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

// userAgent is the User-Agent sent on every guns.lol request, and the UA passed
// to CapMonster when minting a cf_clearance cookie. cf_clearance is bound to the
// UA that solved the Cloudflare challenge, so the same value must be used
// everywhere for the whole flow to stay consistent.
var userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.0.0 Safari/537.36"

func fatalf(format string, a ...any) {
	if activeSpinner != nil {
		activeSpinner.fail()
	}
	fmt.Fprintf(os.Stderr, "\n  %s %s\n\n", red(bold("✗")), fmt.Sprintf(format, a...))
	os.Exit(1)
}

func main() {
	username := flag.String("username", "", "Profile username")
	capmonsterKey := flag.String("capmonster-key", "", "CapMonster API key for solving Turnstile and minting cf_clearance")
	proxy := flag.String("proxy", "", "Proxy URL for guns.lol requests (e.g. http://user:pass@host:port). A {session} placeholder is replaced with a random token each run for rotating sticky-session proxies")
	linkID := flag.String("link-id", "", "Link UUID to record a click event instead of a profile view")
	flag.Parse()

	// Proxy is applied before any request (and passed to CapMonster when minting
	// cf_clearance) so the whole flow shares one egress IP. SetProxy also
	// resolves any {session} placeholder and records the result in proxyURL.
	if *proxy != "" {
		if err := SetProxy(*proxy); err != nil {
			fatalf("Invalid proxy URL: %s", err)
		}
	}

	banner()

	if *linkID != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
		defer cancel()
		sp := startSpinner("Acquiring Cloudflare clearance")
		if _, err := FetchWorkerData(ctx, *username); err != nil {
			fatalf("Error acquiring Cloudflare clearance: %s", err)
		}
		sp.succeed("Acquired Cloudflare clearance")

		sp = startSpinner("Submitting link click")
		if err := SubmitLinkClick(*username, *linkID); err != nil {
			fatalf("Error submitting link click: %s", err)
		}
		sp.succeed("Submitted link click")

		doneln("Link click recorded for %s", cyan(*linkID))
		return
	}

	if *username == "" {
		fatalf("Missing required flag %s (example: guns-solver -username <username>)", bold("-username"))
	}
	if *capmonsterKey == "" {
		fatalf("Missing required flag %s (needed to solve Cloudflare Turnstile)", bold("-capmonster-key"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	infoln("target %s", cyan("guns.lol/"+*username))

	sp := startSpinner("Fetching worker data")
	workerData, err := FetchWorkerData(ctx, *username)
	if err != nil {
		fatalf("Error fetching worker data: %s", err)
	}
	if workerData == nil {
		fatalf("Error fetching worker data: no data returned")
	}
	data := *workerData
	sp.succeedDetail("Fetched worker data", fmt.Sprintf("difficulty %d · challenge %s", data.Difficulty, data.ID))

	sp = startSpinner("Fetching current PoW module")
	powMod, err := FetchPowModule(ctx, data.ID, data.WorkerURL)
	if err != nil {
		fatalf("Error fetching PoW module: %s", err)
	}
	sp.succeed("Fetched current PoW module")

	sp = startSpinner("Solving PoW and Turnstile")

	type wasmResult struct {
		res      *WasmResult
		err      error
		duration time.Duration
	}
	type turnstileResult struct {
		token    string
		err      error
		duration time.Duration
	}

	wasmCh := make(chan wasmResult, 1)
	turnstileCh := make(chan turnstileResult, 1)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		res, err := SolveWithWasm(ctx, powMod, data.Challenge, data.Difficulty, strconv.FormatInt(data.Timestamp, 10), data.Nonce, data.Seal)
		wasmCh <- wasmResult{res, err, time.Since(start)}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		token, err := SolveTurnstile(ctx, *capmonsterKey, "https://guns.lol/"+*username, data.Action, data.CData)
		turnstileCh <- turnstileResult{token, err, time.Since(start)}
	}()

	wg.Wait()

	wasmR := <-wasmCh
	if wasmR.err != nil {
		fatalf("Error solving PoW: %s", wasmR.err)
	}
	tR := <-turnstileCh
	if tR.err != nil {
		fatalf("Error solving Turnstile: %s", tR.err)
	}
	sp.succeed("Solved PoW and Turnstile")
	infoln("PoW        %s · proof %s", wasmR.duration.Round(time.Millisecond), truncateMiddle(wasmR.res.Oo, 24))
	infoln("Turnstile  %s", tR.duration.Round(time.Millisecond))

	sp = startSpinner("Submitting solution")
	err = SubmitSolution(ctx, SolutionPayload{
		Username:  *username,
		Version:   data.Version,
		ID:        data.ID,
		Timestamp: data.Timestamp,
		Nonce:     data.Nonce,
		Seal:      data.Seal,
		Challenge: data.Challenge,
		Proof:     wasmR.res.Oo,
		Token:     tR.token,
		Referrer:  "",
		DeviceNum: 0,
	}, *capmonsterKey)
	if err != nil {
		fatalf("Error submitting solution: %s", err)
	}
	sp.succeed("Submitted solution")

	doneln("View recorded for %s", cyan("guns.lol/"+*username))
}
