package main

import (
	"context"
	"flag"
	"os"
	"time"
)

func main() {
	username := flag.String("username", "", "Profile username")
	flag.Parse()

	if *username == "" {
		println("Please provide a username using the -username flag.")
		os.Exit(1)
	}

	workerData, err := FetchWorkerData(*username)
	if err != nil {
		println("Error fetching worker data:", err.Error())
		os.Exit(1)
	} else if workerData == nil {
		println("Failed to fetch worker data.")
		os.Exit(1)
	}
	data := *workerData
	println("Solving PoW...")

	// to test if the solver works with already known values. See assets/request.http for reference
	data.Nonce = "Zn6F5oHCaRnX1vT9zQ81a4Ty7faAAGCb"
	data.Underscore2xa = "oUAFFCgIIAUAAgEDBFRWE1Bbr_TxYzk5MjZmYWI4NTZkZjlhNjQ5YmY1MWM1NWUyOTk5ZDJmZjE1Y2ZiZGNjYjhjY2Q1YzgzNzMzNTc2NGIXqZPHa_eHIw"
	data.O09 = "1867764c175fc766ce086157515853444b999d7f64d03cbe1d96b78cf6224e02"

	solver := NewGunsSolver(data.Nonce, data.Underscore2xa, data.O09)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	res := solver.SolveConcurrent(ctx)
	if res == nil {
		println("Failed to solve PoW within timeout")
		os.Exit(1)
	}
	println("Solved!")

	err = SubmitSolution(SolutionPayload{
		Username:          *username,
		Nonce:             data.Nonce,
		O09:               data.O09,
		Timestamp:         data.OriginalTimestamp,
		Underscore2xa:     data.Underscore2xa,
		ResultHash:        "db54e", //res.Hash,
		TurnstileResponse: "",      // TODO: Implement Turnstile challenge solving
	})
	if err != nil {
		println("Error submitting solution:", err.Error())
		os.Exit(1)
	}

	println("Solution submitted successfully!")
}
