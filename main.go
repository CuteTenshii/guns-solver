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
	}
	if workerData == nil {
		println("Failed to fetch worker data.")
		os.Exit(1)
	}
	data := *workerData
	println("Solving PoW...")
	// to test if the solver works with already known values. See assets/request.http for reference
	data.Nonce = "3e3xwVmOKMOYJxUZ"
	data.PublicSalt = "839198c382afb05ed1c5309486a1c9a4a634439866535f912db27d6a52d"
	data.Challenge = "ad0bc1b3adb82e6af84f81415318ecebbed06b6bb31747a55138e700a148e1e9"

	solver := NewGunsSolver(data.PublicSalt, data.Challenge, data.Difficulty, data.Nonce)

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
		PublicSalt:        data.PublicSalt,
		Challenge:         data.Challenge,
		Nonce:             res.Nonce,
		Timestamp:         data.OriginalTimestamp,
		Difficulty:        data.Difficulty,
		ResultHash:        res.Hash,
		TurnstileResponse: "", // TODO: Implement Turnstile challenge solving
	})
	if err != nil {
		println("Error submitting solution:", err.Error())
		os.Exit(1)
	}

	println("Solution submitted successfully!")
}
