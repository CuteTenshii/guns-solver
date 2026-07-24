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
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(1)
}

func main() {
	username := flag.String("username", "", "Profile username")
	capmonsterKey := flag.String("capmonster-key", "", "CapMonster API key for solving Turnstile and minting cf_clearance")
	proxy := flag.String("proxy", "", "Proxy URL for guns.lol requests (e.g. http://user:pass@host:port)")
	linkID := flag.String("link-id", "", "Link UUID to record a click event instead of a profile view")
	flag.Parse()

	// Proxy is applied before any request (and passed to CapMonster when minting
	// cf_clearance) so the whole flow shares one egress IP.
	if *proxy != "" {
		if err := SetProxy(*proxy); err != nil {
			fatalf("Invalid proxy URL: %s", err)
		}
	}
	proxyURL = *proxy

	if *linkID != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
		defer cancel()
		if _, err := FetchWorkerData(ctx, *username); err != nil {
			fatalf("Error acquiring Cloudflare clearance: %s", err)
		}
		if err := SubmitLinkClick(*username, *linkID); err != nil {
			fatalf("Error submitting link click: %s", err)
		}
		fmt.Println("Link click submitted successfully!")
		return
	}

	if *username == "" {
		fatalf("Missing required flag: -username\nExample: guns-solver -username <username>")
	}
	if *capmonsterKey == "" {
		fatalf("Missing required flag: -capmonster-key (needed to solve Cloudflare Turnstile)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	fmt.Println("Fetching worker data...")
	workerData, err := FetchWorkerData(ctx, *username)
	if err != nil {
		fatalf("Error fetching worker data: %s", err)
	}
	if workerData == nil {
		fatalf("Error fetching worker data: no data returned")
	}
	data := *workerData

	fmt.Println("Fetching current PoW module...")
	powMod, err := FetchPowModule(ctx, data.ID, data.WorkerURL)
	if err != nil {
		fatalf("Error fetching PoW module: %s", err)
	}

	fmt.Println("Solving PoW and Turnstile simultaneously...")

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
	fmt.Printf("PoW solved in %s! proof=%s\n", wasmR.duration.Round(time.Millisecond), wasmR.res.Oo)

	tR := <-turnstileCh
	if tR.err != nil {
		fatalf("Error solving Turnstile: %s", tR.err)
	}
	fmt.Printf("Turnstile solved in %s!\n", tR.duration.Round(time.Millisecond))

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

	fmt.Println("Solution submitted successfully!")
}
