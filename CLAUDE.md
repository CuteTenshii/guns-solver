# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go build          # build the binary
go run .          # run without building
./guns-solver -username <username> -capmonster-key <key>   # run the solver
go vet ./...      # lint
go test ./...     # run tests (includes a ~3s known-answer WASM solve)
```

Depends on `github.com/tetratelabs/wazero` (runs the PoW WASM) plus the standard library.

## Architecture

This is a Go reverse-engineering tool that automates the proof-of-work challenge guns.lol uses to record profile views (`POST /api/analytics/view`).

### Files

- **`main.go`** — Entry point. Parses flags (`-username`, `-capmonster-key`, `-proxy`, `-link-id`), fetches the challenge, solves the PoW and Turnstile concurrently, then submits. `userAgent` is shared across every request (cf_clearance is bound to it).
- **`worker_data.go`** — `FetchWorkerData(username)` fetches the profile page and `parseGppChallenge` extracts the `_gpp_ch` object out of the Next.js flight payload (it lives inside a JS string literal, so the captured JSON is un-escaped before `json.Unmarshal`). Also holds the clearance-cookie plumbing and the `guns_clearance` 401 interstitial flow (`solveChallenge`), a **separate** subsystem that uses the `_gs_sets` format and the embedded clearance binary.
- **`pow_module.go`** — `FetchPowModule(id, workerPath)` resolves the current view-PoW module for a challenge. The worker, glue, and binary all live under `/_challenge/pow/<id>/`; `fetchPowModule` fetches the worker, XOR-decodes the glue import name, parses the glue for the wasm URL and the rotated constructor/solve export symbols, and fetches the binary, returning a `powModule{wasm, ctorExport, solveExport}`. Results are cached to disk keyed by challenge id (`powCacheRoot`, under the temp dir) — the module is immutable under a given id and the id only changes on rotation, so a cache hit is always current and skips all three round-trips.
- **`wasm.go`** — `SolveWithWasm(m *powModule, ...)` runs a given PoW module via wazero, reimplementing the wasm-bindgen JS glue (a JS heap with free-list, the host imports, string passing). The constructor ABI is `(challenge, difficulty, orgTs, nonce, seal) -> ptr`; the solve method writes 3×i32 to a return slot and yields a result object whose `_oo` field is the proof. `clearancePowModule()` wraps the **embedded** stable clearance binary (`assets/gpp_gunslol_bg.wasm`, exports `gunssolver_new`/`gunssolver_solve_pow`) as a `powModule`.
- **`turnstile.go`** — `SolveTurnstile` solves the Cloudflare Turnstile (sitekey `0x4AAAAAAAgU7T2niLQD-TLm`) via CapMonster, using the challenge's `action` and `cd` values.
- **`cf_clearance.go`** — `SolveCfClearance` mints a `cf_clearance` cookie via CapMonster when the analytics endpoint answers with a 403 Managed Challenge.
- **`submit.go`** — `SubmitSolution` builds the positional-array body and POSTs it to `/api/analytics/view`, minting cf_clearance and retrying once on 403. `SubmitLinkClick` still targets the older `/api/analytics/record` labeled-object endpoint.
- **`ui.go`** — Dependency-free terminal styling for the CLI: ANSI helpers (`bold`/`dim`/`cyan`/…) gated by `colorEnabled` (auto-disabled on non-TTY stdout, `NO_COLOR`, or `TERM=dumb`), a braille `spinner` for per-step progress that resolves to `✓`/`✗`, and `banner`/`infoln`/`warnln`/`doneln` line helpers. `fatalf` tears down the `activeSpinner` before printing an error; `infoln`/`warnln` clear a running spinner's frame so subsystem logs print above it without corrupting the animation. In non-TTY mode the spinner degrades to a static `→` line.
- **`assets/`** — `gpp_gunslol_bg.wasm`, the embedded clearance binary (exports `gunssolver_new`; used only by `clearancePowModule()`). `gpp_gunslol_bg_old.wasm` is a retained older variant, unreferenced by code. The rotating view-PoW binaries are not stored here — they are fetched and cached under the temp dir.

### View proof-of-work flow (rotating)

guns.lol replaced the old cross-origin WASM solver with a **same-origin, server-issued worker** served from `/_challenge/pow/<id>/<worker>`. The worker (XOR-obfuscated) dynamically imports a wasm-bindgen glue module and runs `new Solver(c, Number(d), String(t), n, s).solve()`, posting back `_oo`.

**This module rotates periodically (≈hourly, server-side): the binary content, its hashed export symbols, the glue/worker filenames, and the worker's XOR key all change together.** The binary embeds a rotation-specific key the constructor uses to validate the challenge `seal`; a stale binary rejects a current seal with **error code 1003**. So the tool must **fetch the current module per run** (`FetchPowModule`) rather than embed one — do not re-embed the view binary.

What is stable across rotations: the host imports (`__wbg*`), the `malloc`/`realloc`/`add_to_stack_pointer` export names, the constructor ABI `(challenge, difficulty, orgTs, nonce, seal) -> ptr`, and the func indices (ctor=8, solve=7). Only the ctor/solve export *names* rotate — hence `FetchPowModule` parses them out of the glue (the wasm export called with 9 args is the ctor; the one whose 2nd arg is `this.__wbg_ptr` is solve).

The `guns_clearance` 401 interstitial is a **separate** subsystem and did not change; it keeps using the stable embedded binary via `clearancePowModule()`.

### `_gpp_ch` schema (short keys → fields)

The challenge object on the profile page uses single-letter keys (`WorkerData` expands them):

- `v` → `Version` — payload schema version (3)
- `e` → `ID` — challenge id (first path segment of the worker URL)
- `u` → `WorkerURL` — same-origin path of the PoW worker module
- `t` → `Timestamp` — challenge issue time (unix seconds), passed to the solver as a string
- `n` → `Nonce`
- `s` → `Seal` — opaque server seal, fed to the solver and replayed in the solution
- `c` → `Challenge` — 64-char hex challenge input
- `d` → `Difficulty`
- `cd` → `CData` — Turnstile `data-cdata`
- `a` → `Action` — Turnstile `data-action` (`guns_view`)

### Submission payload

`/api/analytics/view` now takes a positional JSON array (`Content-Type: application/json`):

```
[token, [v, e, t, n, s, c, proof], username, deviceNum, referrer]
```

`deviceNum` is a numeric enum: desktop=0, mobile=1, tablet=2. `proof` is the solver's `_oo`.

### Testing

- `go test ./...` runs offline unit tests: `TestSolveWithWasmKnownSample` (known-answer against the embedded clearance binary), `TestParseGlue`/`TestDecodeGlueImport` (the runtime glue/worker parsing), `TestPowModuleCache` (the disk cache round-trip and path-safety), `TestParseGppChallenge`, and `TestBuildViewPayloadShape`.
- `GUNS_LIVE=1 [GUNS_USER=<name>] go test -run TestLivePowSolve` runs the network-gated integration test: fetch the current module and solve the current challenge. This is the regression test for the rotation 1003 bug.

### Known unknowns

- The live fetch-and-solve path is verified (no 1003), but a full end-to-end **submit** to `/api/analytics/view` (real CapMonster key + cf_clearance) accepting the proof has not been confirmed against the origin.
