# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go build          # build the binary
go run .          # run without building
./guns-solver -username <username>   # run the solver
go vet ./...      # lint
go test ./...     # run tests (none exist yet)
```

No external dependencies — `go.mod` uses only the standard library.

## Architecture

This is a Go reverse-engineering tool that automates the proof-of-work challenge used by guns.lol to record profile views. It has no tests and no build scripts beyond `go build`.

### Files

- **`main.go`** — Entry point. Parses `-username` flag, fetches worker data, runs the solver, and submits. Currently hardcodes known test values for `Nonce`, `Underscore2xa`, and `O09` to validate the solver against a known challenge. The actual HTTP submission in `SubmitSolution` is commented out.
- **`worker_data.go`** — `FetchWorkerData(username)` scrapes the guns.lol profile page HTML via regex to extract `_n` (nonce), `o09`, `_2xa`, and `org_ts` values needed to construct the PoW challenge. Handles `guns_clearance` cookie redirect flow automatically.
- **`solver.go`** — `GunsSolver` struct. `SolveConcurrent` runs CPU-count goroutines that brute-force a nonce string such that `SHA-256(Nonce + Underscore2xa + candidate)` meets the difficulty target. Supports two modes: `LeadingHexNibbles` (default, difficulty=5 means 5 leading hex zeros) and `LeadingZeroBits`. Nonce candidates are enumerated via `incrementStringWithCharset` over an alphanumeric charset.
- **`submit.go`** — `SubmitSolution` builds the `POST /api/analytics/record` JSON payload and prints it. The actual HTTP call is commented out pending full reverse engineering.
- **`assets/`** — Research artifacts: `gpp_gunslol.js` (WASM JS loader), `next.js` (profile page script), `request.http` (reference request), `gpp_gunslol_bg.wasm` (the WASM binary).

### Current Status / Known Unknowns

The payload field `_gpp_ch` contains a `seal` object (`{_oo, seal}`) whose values are not yet known. The Turnstile token (`_t`) is not solved. The WASM module (`gpp_gunslol_bg.wasm`, written in Rust) has been converted to C via `wasm2c` for analysis — its `solve_pow` output is a 5-character hex string used as `ResultHash`.

The `_` and `_2` dynamic hashes inside the old `__hcm` payload structure have been partially reversed; see README.md and commit history for current progress.

### Key Field Naming

The fields use obfuscated names from the JS source:
- `Underscore2xa` → `_2xa` in the page HTML / `ps` (public salt) conceptually
- `O09` → `o09` in the page HTML (a 64-char hex string, likely the challenge)
- `Nonce` → `_n` in the worker script
