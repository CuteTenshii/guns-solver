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

func fatalf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(1)
}

func main() {
	username := flag.String("username", "", "Profile username")
	capmonsterKey := flag.String("capmonster-key", "", "CapMonster API key for Turnstile solving")
	proxy := flag.String("proxy", "", "Proxy URL for guns.lol requests (e.g. http://user:pass@host:port)")
	linkID := flag.String("link-id", "", "Link UUID to record a click event instead of a profile view")
	flag.Parse()

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

	if *proxy != "" {
		if err := SetProxy(*proxy); err != nil {
			fatalf("Invalid proxy URL: %s", err)
		}
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
		res, err := SolveWithWasm(ctx, data.O09, 5, strconv.FormatInt(data.OriginalTimestamp, 10), data.Nonce, data.Underscore2xa)
		wasmCh <- wasmResult{res, err, time.Since(start)}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		token, err := SolveTurnstile(ctx, *capmonsterKey, "https://guns.lol/"+*username)
		turnstileCh <- turnstileResult{token, err, time.Since(start)}
	}()

	wg.Wait()

	wasmR := <-wasmCh
	if wasmR.err != nil {
		fatalf("Error solving PoW: %s", wasmR.err)
	}
	fmt.Printf("PoW solved in %s! _oo=%s seal=%s\n", wasmR.duration.Round(time.Millisecond), wasmR.res.Oo, wasmR.res.Seal)

	tR := <-turnstileCh
	if tR.err != nil {
		fatalf("Error solving Turnstile: %s", tR.err)
	}
	fmt.Printf("Turnstile solved in %s!\n", tR.duration.Round(time.Millisecond))

	err = SubmitSolution(SolutionPayload{
		Username:          *username,
		Nonce:             data.Nonce,
		O09:               data.O09,
		Timestamp:         data.OriginalTimestamp,
		Underscore2xa:     data.Underscore2xa,
		Oo:                wasmR.res.Oo,
		Seal:              wasmR.res.Seal,
		TurnstileResponse: tR.token,
	})
	if err != nil {
		fatalf("Error submitting solution: %s", err)
	}

	fmt.Println("Solution submitted successfully!")
}
