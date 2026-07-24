# guns.lol solver

<sub>Last update to the WASM script: July 24, 2026 — guns.lol now rotates the PoW module (binary + export names) server-side, so the solver fetches and caches it per challenge.</sub>

A solver for the guns.lol WebAssembly script, which is used to record views on profile pages.
Made it because I was bored and wanted to see how it worked.

By default it records one view per run. To record many at once — each on its own IP — use the built-in `-count`/`-concurrency` flags; scroll down to [Botting your views](#botting-your-views).

For an alternative to guns.lol, check out [Miwa.lol](https://miwa.lol)! It's better than them :)

## Usage

### Prerequisites

- [Go](https://go.dev/dl/)
- A [CapMonster](https://capmonster.cloud/) API key (solves the Turnstile and mints the Cloudflare `cf_clearance` cookie)

### Steps 

- Download or clone the source code
- Open a Command Prompt/Terminal in the same folder, and run `go build .`
- Then run the program:
```bash
Usage of ./guns-solver.exe:
  -capmonster-key string
        CapMonster API key for solving Turnstile and minting cf_clearance
  -concurrency int
        Maximum views to record simultaneously (default: min(count, 5))
  -count int
        Number of profile views to record (default 1)
  -link-id string
        Link UUID to record a click event instead of a profile view
  -proxy string
        Proxy URL for guns.lol requests (e.g. http://user:pass@host:port). A {session} placeholder is replaced with a fresh token per worker for rotating sticky-session proxies
  -username string
        Profile username
```

### Examples

On Linux and macOS, the `.exe` extension is not present so remove it.

To add a view to a user:
```shell
./guns-solver.exe -username <username> -capmonster-key ... -proxy ...
```

To add many views at once (50 total, 8 in flight, each on its own IP):
```shell
./guns-solver.exe -username <username> -capmonster-key ... -count 50 -concurrency 8 -proxy "http://user:pass_session-{session}@host:port"
```

To add a link click:
```shell
./guns-solver.exe -link-id <link-id>
```

## Cloudflare

guns.lol puts a Cloudflare bot check on the analytics endpoint that gates by the presence of `cf_clearance`, a client without this cookie gets a Cloudflare challenge (403).

When the tool hits that challenge, it hands the interstitial to CapMonster's Cloudflare Challenge task (`cloudflareTaskType: cf_clearance`), which returns a `cf_clearance` cookie. The tool caches it under the temp dir (`<tmp>/guns-solver-pow/cf_clearance.txt`, reused across runs) and retries the request directly.

`cf_clearance` is bound to a single IP and User-Agent, so `-proxy` is **required** to mint one, and CapMonster egresses through that same proxy so the cookie matches the tool's IP.

## Botting your views

**Note:** You NEED proxies with a randomized IP for this to work. This is NOT optional since I estimate they rate limit user views to **1 view/per ip/per day**.

### Rotating proxies (sticky sessions)

If your proxy provider supports **sticky sessions** (a session id in the credentials that pins one exit IP), put a `{session}` placeholder where that id goes. Every worker replaces it with a fresh random token, so each view gets a new IP — exactly what you want when firing off many at once.

For example, with [IPRoyal](https://iproyal.com):
```
-proxy "http://username:password_session-{session}_lifetime-30s@geo.iproyal.com:12345"
```
Without the placeholder the proxy URL is used as-is (all views share one IP), so make sure to add it for botting — the tool warns you if it's missing.

Personally I got banned after multiple days **BUT** I've seen people I botted not getting banned after multiple weeks. Their staff is lazy asf

### Firing off many views

This is built in — no shell loops or GNU Parallel needed. Use `-count` for how many views to record and `-concurrency` for how many run at the same time:
```shell
./guns-solver -username <username> -capmonster-key ... -count 100 -concurrency 40 -proxy "http://user:pass_session-{session}@host:port"
```
Each view runs on its own proxy session (its own IP and `cf_clearance`), so the counts are what you'd expect. A live progress board shows each worker; when the output is piped/redirected it falls back to one line per completed view. Press `Ctrl-C` to stop early — it reports what already completed.
