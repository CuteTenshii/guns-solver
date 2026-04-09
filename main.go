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

func main() {
	username := flag.String("username", "", "Profile username")
	capmonsterKey := flag.String("capmonster-key", "", "CapMonster API key for Turnstile solving")
	proxy := flag.String("proxy", "", "Proxy URL for guns.lol requests (e.g. http://user:pass@host:port)")
	linkID := flag.String("link-id", "", "Link UUID to record a click event instead of a profile view")
	flag.Parse()

	if *username == "" {
		println("Please provide a username using the -username flag.")
		os.Exit(1)
	}

	if *proxy != "" {
		if err := SetProxy(*proxy); err != nil {
			fmt.Println("Invalid proxy URL:", err.Error())
			os.Exit(1)
		}
	}

	if *linkID != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
		defer cancel()
		if _, err := FetchWorkerData(ctx, *username); err != nil {
			fmt.Println("Error acquiring clearance:", err.Error())
			os.Exit(1)
		}
		if err := SubmitLinkClick(*username, *linkID); err != nil {
			println("Error submitting link click:", err.Error())
			os.Exit(1)
		}
		println("Link click submitted successfully!")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	workerData, err := FetchWorkerData(ctx, *username)
	if err != nil {
		println("Error fetching worker data:", err.Error())
		os.Exit(1)
	} else if workerData == nil {
		println("Failed to fetch worker data.")
		os.Exit(1)
	}
	data := *workerData

	fmt.Println("Solving PoW and Turnstile simultaneously...")

	type wasmResult struct {
		res *WasmResult
		err error
	}
	type turnstileResult struct {
		token string
		err   error
	}

	wasmCh := make(chan wasmResult, 1)
	turnstileCh := make(chan turnstileResult, 1)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := SolveWithWasm(ctx, data.O09, 5, strconv.FormatInt(data.OriginalTimestamp, 10), data.Nonce, data.Underscore2xa)
		wasmCh <- wasmResult{res, err}
	}()

	if *capmonsterKey != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			token, err := SolveTurnstile(ctx, *capmonsterKey, "https://guns.lol/"+*username)
			turnstileCh <- turnstileResult{token, err}
		}()
	}

	wg.Wait()

	wasmR := <-wasmCh
	if wasmR.err != nil {
		fmt.Println("Error solving PoW:", wasmR.err.Error())
		os.Exit(1)
	}
	fmt.Printf("PoW solved! _oo=%s seal=%s\n", wasmR.res.Oo, wasmR.res.Seal)

	var turnstileToken string
	if *capmonsterKey != "" {
		tR := <-turnstileCh
		if tR.err != nil {
			fmt.Println("Error solving Turnstile:", tR.err.Error())
			os.Exit(1)
		}
		turnstileToken = tR.token
		fmt.Println("Turnstile solved!")
	}

	err = SubmitSolution(SolutionPayload{
		Username:          *username,
		Nonce:             data.Nonce,
		O09:               data.O09,
		Timestamp:         data.OriginalTimestamp,
		Underscore2xa:     data.Underscore2xa,
		Oo:                wasmR.res.Oo,
		Seal:              wasmR.res.Seal,
		TurnstileResponse: turnstileToken,
	})
	if err != nil {
		println("Error submitting solution:", err.Error())
		os.Exit(1)
	}

	println("Solution submitted successfully!\n")
}
