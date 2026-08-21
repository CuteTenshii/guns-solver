# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go build          # build the binary
go run .          # run without building
./guns-solver -username <username> -captcha-key <key>   # record one view (CapMonster by default)
./guns-solver -username <username> -captcha-provider solverify -captcha-key <key>   # same, via Solverify
./guns-solver -username <username> -captcha-provider uncaptcha -captcha-key <key>   # same, via uncaptcha.io
./guns-solver -username <username> -captcha-key <key> -count 50 -concurrency 8 -proxy 'http://user:pass_session-{session}@host:port'   # record many views in parallel, each on its own IP
go vet ./...      # lint
go test ./...     # run tests (includes a ~3s known-answer WASM solve)
```

Depends on `github.com/tetratelabs/wazero` (runs the PoW WASM) plus the standard library.

## Architecture

This is a Go reverse-engineering tool that automates the proof-of-work challenge guns.lol uses to record profile views (`POST /api/analytics/view`).

### Files

- **`main.go`** — Entry point. Parses flags (`-username`, `-captcha-provider`, `-captcha-key`, `-proxy`, `-link-id`, `-count`, `-concurrency`), resolves them into a `captchaProvider`, then dispatches: `runLinkClick` (link event), `runSingle` (one view, animated per-step spinners), or `runParallel` (`-count` > 1, live board). `userAgent` is shared across every request (cf_clearance is bound to it).
- **`session.go`** — `session` is the per-worker request state (own `http.Client`, resolved proxy, `guns_clearance`/`cf_clearance` cookies). `newSession(rawProxy, persist)` resolves a `{session}` proxy token to a fresh value so each session egresses through a distinct IP; `persist` mirrors the clearance cookies to disk (single-run only — IP-bound cookies must not leak across parallel workers). Holds the cookie plumbing (`addClearanceCookies`, `captureCfClearance`, `save*`) and `infof`/`warnf`, which route through `onLog` (the board) in parallel mode. A session is single-goroutine, so it needs no locking.
- **`worker_data.go`** — `(*session).FetchWorkerData(username)` fetches the profile page and `parseGppChallenge` extracts the `_gpp_ch` object out of the Next.js flight payload (it lives inside a JS string literal, so the captured JSON is un-escaped before `json.Unmarshal`). Also holds the `guns_clearance` 401 interstitial flow (`solveChallenge`), a **separate** subsystem that uses the `_gs_sets` format and the embedded clearance binary, plus `resolveProxySession`/`randomSessionID`.
- **`run.go`** — `(*session).runView(cfg, reporter)` is the single record-a-view pipeline used by both single and parallel modes: fetch challenge → fetch PoW module → `solvePowAndTurnstile` (concurrent) → submit. The `reporter` interface has two impls: `spinnerReporter` (single-run animated steps) and `rowReporter` (parallel board row).
- **`board.go`** — `board` is the live multi-slot progress display for parallel mode: one animated row per worker slot + a running summary, with a persistent `✓`/`✗` line emitted *above* the live region per finished view. Falls back to one plain line per completion when color is disabled (piped output), which is what fixed the "external `parallel` strips styling" problem. All mutating methods take `b.mu`; `render`/`logLine` assume it is held.
- **`pow_module.go`** — `(*session).FetchPowModule(id, workerPath)` resolves the current view-PoW module for a challenge. The worker, glue, and binary all live under `/_challenge/pow/<id>/`; `fetchPowModule` fetches the worker, XOR-decodes the glue import name, parses the glue for the wasm URL and the rotated constructor/solve export symbols, and fetches the binary, returning a `powModule{wasm, ctorExport, solveExport}`. Results are cached to disk keyed by challenge id (`powCacheRoot`, under the temp dir) — the module is immutable under a given id and the id only changes on rotation, so a cache hit is always current and skips all three round-trips. Temp files use a random suffix so concurrent workers caching the same rotation don't clobber each other.
- **`wasm.go`** — `SolveWithWasm(m *powModule, ...)` runs a given PoW module via wazero, reimplementing the wasm-bindgen JS glue (a JS heap with free-list, the host imports, string passing). The constructor ABI is `(challenge, difficulty, orgTs, nonce, seal) -> ptr`; the solve method writes 3×i32 to a return slot and yields a result object whose `_oo` field is the proof. `clearancePowModule()` wraps the **embedded** stable clearance binary (`assets/gpp_gunslol_bg.wasm`, exports `gunssolver_new`/`gunssolver_solve_pow`) as a `powModule`.
- **`captcha.go`** — The provider abstraction. `captchaProvider` has `Name`/`SolveTurnstile`/`SolveCfClearance`; `turnstileRequest`, `clearanceRequest` and `clearanceResult` are the provider-neutral payloads. Also holds `turnstileSiteKey` (`0x4AAAAAAAgU7T2niLQD-TLm`), `newCaptchaProvider`, the `captchaFlags.resolve()` flag rules, `parseProxyParts` (services want the proxy split into fields, not a URL), and the shared `postJSON`/`postJSONWithHeaders`/`pollTask` helpers. `postJSONWithHeaders` adds request headers and returns the HTTP status; `postJSON` is the thin wrapper the two task-id services use, since they authenticate in the body and always answer 200. Implementations are stateless, so one instance serves every parallel worker.
- **`capmonster.go`** — CapMonster provider. Turnstile via `TurnstileTaskProxyless`; `cf_clearance` via the Cloudflare Challenge task (`TurnstileTask` + `cloudflareTaskType: cf_clearance`), which replays the captured 403 interstitial and needs the session's proxy. Integer task ids, `ready` status, token in `solution.token`.
- **`solverify.go`** — Solverify provider (`https://solver.solverify.net`). Turnstile via the `turnstile` task (note the different field names: `action`/`cdata`, not `pageAction`/`data`); `cf_clearance` via the `interstitial` task, which navigates the challenged request itself (no captured HTML) and requires an **http** proxy with the port as a *string*. It must be pointed at `clearanceRequest.ChallengeURL` with `ChallengeMethod` (`POST /api/analytics/view`), not the profile page: Cloudflare only issues `cf_clearance` when it actually challenges, and a browser GETting the profile page — or GETting the POST-only analytics endpoint — sails through and comes back with an empty jar. String task ids, `completed`/`failed` statuses, token in `solution.value`, cookie in `solution.cookies.cf_clearance`, and the browser identity in `solution.useragent` — mismatches against `userAgent` are warned about in `SubmitSolution`.
- **`uncaptcha.go`** — uncaptcha.io provider (`https://api.uncaptcha.io/v1`). The odd one out: **synchronous** — a single `POST /task/execute` blocks until solved or the service's 30s execution limit, so there is no task id, no `getTaskResult`, and `pollTask` is unused. Authenticates with an `X-Api-Key` **header** (not a body field), takes the proxy as one **URL string** (not split fields), nests input under `task_data`, and signals failure through the real HTTP status plus `{success:false, message}` rather than an `errorId`. Turnstile via the `turnstile` task (`sitekey`/`cdata`, snake_case) → `data.solution.token`; `cf_clearance` via the `waf` task, which replays our captured 403 interstitial (`html`) like CapMonster rather than navigating itself like Solverify → `data.solution.clearance`, with the bound identity read case-insensitively out of `solution.headers["user-agent"]` (`solvedUserAgent`). Because the call is synchronous, `execute` retries `429`/5xx itself with a doubling backoff (`uncaptchaRetryable`) — 429 means the key's concurrency limit is saturated, which is the normal case when `-concurrency` exceeds it, and failed solves are refunded. The `wafauto` task is deliberately unused: it makes the post-challenge request itself, but ours must carry our own signed view payload. guns.lol sets no `__cf_bm`, so the WAF solution's `cf_bm`/`cf_rt`/`attributes` extras are ignored.
- **`submit.go`** — `(*session).SubmitSolution` builds the positional-array body and POSTs it to `/api/analytics/view`, minting cf_clearance through the `captchaProvider` (bound to the session's `proxyURL`) and retrying once on 403. `(*session).SubmitLinkClick` still targets the older `/api/analytics/record` labeled-object endpoint.
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

### Parallel view recording

`-count N` records N views; `-concurrency K` (default `min(count, 5)`) caps how many run at once. `runParallel` fans out K worker goroutines that pull view indices off an atomic counter until `count` is exhausted. **Each view runs on its own `newSession`**, so a `{session}` proxy resolves to a fresh token — a distinct exit IP and a distinct `cf_clearance` — per view; without a `{session}` proxy the workers share one IP and the tool warns that the views will likely de-dupe. The whole point is that this is built in: it replaces piping the binary through GNU `parallel`, which stripped the TTY and killed the styling. Ctrl-C (`signal.NotifyContext`) stops claiming new views and reports what completed. `-link-id` and `-count`/`-concurrency` are independent — a link click is always a single request.

### Testing

- `go test ./...` runs offline unit tests: `TestSolveWithWasmKnownSample` (known-answer against the embedded clearance binary), `TestParseGlue`/`TestDecodeGlueImport` (the runtime glue/worker parsing), `TestPowModuleCache` (the disk cache round-trip and path-safety), `TestParseGppChallenge`, `TestBuildViewPayloadShape`, and the captcha-provider tests in `captcha_test.go` (flag resolution, proxy splitting, and pinned request/response JSON shapes for all three services — the field names differ per provider, so a rename breaks solving silently — plus `TestUncaptchaRetryable` for which HTTP statuses are re-sent and `TestSolvedUserAgent` for the case-insensitive header lookup).
- `GUNS_LIVE=1 [GUNS_USER=<name>] go test -run TestLivePowSolve` runs the network-gated integration test: fetch the current module and solve the current challenge. This is the regression test for the rotation 1003 bug.

### Status

- The full end-to-end flow is confirmed working against the origin: `/api/analytics/view` accepts the proof with a real CapMonster key + cf_clearance, in both single and parallel (`-count`/`-concurrency`) modes. The Solverify and uncaptcha.io providers are built to their published APIs and unit-tested against the documented payloads, but neither has been exercised against the origin with a live key.
