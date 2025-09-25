package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"runtime"
	"strings"
	"sync"
)

// DifficultyMode controls how difficulty is interpreted.
type DifficultyMode int

const (
	// LeadingHexNibbles difficulty = number of leading "0" nibbles in hex string.
	LeadingHexNibbles DifficultyMode = iota
	// LeadingZeroBits difficulty = number of leading zero bits in the raw hash.
	LeadingZeroBits
)

type GunsSolver struct {
	PublicSalt string
	Challenge  string
	Difficulty int
	NonceStart string
	Mode       DifficultyMode
	Charset    string
}

func NewGunsSolver(salt, challenge string, difficulty int, nonce string) GunsSolver {
	return GunsSolver{
		PublicSalt: salt,
		Challenge:  challenge,
		Difficulty: difficulty,
		NonceStart: nonce,
		Mode:       LeadingHexNibbles,
		Charset:    "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
	}
}

type SolveResult struct {
	Nonce string
	Hash  string // hex
}

// SolveConcurrent runs multiple workers to find a nonce. Returns the first found result or nil if ctx canceled.
func (s *GunsSolver) SolveConcurrent(parentCtx context.Context) *SolveResult {
	workers := runtime.NumCPU()
	if s.NonceStart == "" {
		s.NonceStart = "A"
	}

	// quick path for trivial difficulty
	if s.Difficulty <= 0 {
		h := sha256.Sum256([]byte(s.PublicSalt + s.Challenge + s.NonceStart))
		return &SolveResult{Nonce: s.NonceStart, Hash: hex.EncodeToString(h[:])}
	}

	// derive a cancelable context we can cancel internally
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	ch := make(chan string, workers*2)  // work channel (producer closes it)
	resCh := make(chan *SolveResult, 1) // buffered: first writer never blocks

	var wg sync.WaitGroup
	var once sync.Once

	// Producer goroutine: only this goroutine closes ch.
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(ch) // producer is sole closer of ch
		seed := s.NonceStart
		for {
			select {
			case ch <- seed:
				seed = incrementStringWithCharset(seed, s.Charset)
			case <-ctx.Done():
				return
			}
		}
	}()

	// Worker goroutines: range over ch until it's closed.
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for seed := range ch {
				nonce := seed
				for {
					// check for global cancel frequently
					select {
					case <-ctx.Done():
						return
					default:
					}

					hash := sha256.Sum256([]byte(s.PublicSalt + s.Challenge + nonce))
					if s.checkDifficulty(hash[:]) {
						once.Do(func() {
							resCh <- &SolveResult{Nonce: nonce, Hash: hex.EncodeToString(hash[:])}
							cancel() // stop producer & other workers
						})
						return
					}

					nonce = incrementStringWithCharset(nonce, s.Charset)
				}
			}
		}()
	}

	// Wait for either a result or cancellation from parentCtx.
	var result *SolveResult
	select {
	case result = <-resCh:
		// solution found
	case <-ctx.Done():
		// canceled or parent ctx timed out
	}

	// Wait for producer + workers to exit cleanly.
	wg.Wait()
	return result
}

func (s *GunsSolver) checkDifficulty(hash []byte) bool {
	switch s.Mode {
	case LeadingZeroBits:
		return checkLeadingZeroBits(hash, s.Difficulty)
	case LeadingHexNibbles:
		// hex nibble interpretation: e.s. difficulty=4 means hash hex starts with "0000"
		hexStr := hex.EncodeToString(hash)
		return strings.HasPrefix(hexStr, strings.Repeat("0", s.Difficulty))
	default:
		return false
	}
}

// checkLeadingZeroBits returns true if the hash has at least 'bits' leading zero bits.
func checkLeadingZeroBits(hash []byte, bits int) bool {
	if bits <= 0 {
		return true
	}
	fullBytes := bits / 8
	remBits := bits % 8
	// full bytes must be zero
	for i := 0; i < fullBytes; i++ {
		if i >= len(hash) || hash[i] != 0 {
			return false
		}
	}
	if remBits == 0 {
		return true
	}
	if fullBytes >= len(hash) {
		return false
	}
	// top remBits of next byte must be zero
	mask := byte(0xFF << (8 - remBits))
	return (hash[fullBytes] & mask) == 0
}

// incrementStringWithCharset increments s using the provided charset (radix). Charset[0] is "0".
func incrementStringWithCharset(s, charset string) string {
	if s == "" {
		return string(charset[0])
	}
	radix := []rune(charset)
	// build a map for rune->index for speed
	index := make(map[rune]int, len(radix))
	for i, r := range radix {
		index[r] = i
	}

	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		r := runes[i]
		idx, ok := index[r]
		if !ok {
			// if rune not in charset, reset it to first char
			runes[i] = radix[0]
			return string(runes)
		}
		if idx+1 < len(radix) {
			runes[i] = radix[idx+1]
			return string(runes)
		}
		// overflow this position
		runes[i] = radix[0]
	}
	// all positions overflowed -> prepend first charset rune
	return string(radix[0]) + string(runes)
}

func incrementString(s string) string {
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		switch {
		case runes[i] >= '0' && runes[i] < '9':
			runes[i]++
			return string(runes)
		case runes[i] == '9':
			runes[i] = 'A'
			return string(runes)
		case runes[i] >= 'A' && runes[i] < 'Z':
			runes[i]++
			return string(runes)
		case runes[i] == 'Z':
			runes[i] = 'a'
			return string(runes)
		case runes[i] >= 'a' && runes[i] < 'z':
			runes[i]++
			return string(runes)
		case runes[i] == 'z':
			runes[i] = '0'
		}
	}
	// overflow: add a new 'A' at the beginning
	return "A" + string(runes)
}
