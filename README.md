# guns.lol solver

<sub>Last update to the WASM script: July 14, 2026.</sub>

A solver for the guns.lol WebAssembly script, which is used to record views on profile pages.
Made it because I was bored and wanted to see how it worked.

This only sends one request per program execution. If you want to bot your profile views then scroll down.

For an alternative to guns.lol, check out [Miwa.lol](https://miwa.lol)! It's better than them :)

## Usage

### Prerequisites

- [Go](https://go.dev/dl/)
- [Docker](https://docs.docker.com/desktop/#next-steps) (for FlareSolverr)

### Steps 

- Download or clone the source code
- Open a Command Prompt/Terminal in the same folder, and run `go build .`
- Then run the program:
```bash
Usage of ./guns-solver.exe:
  -capmonster-key string
        CapMonster API key for Cloudflare Turnstile solving
  -flaresolverr string
        FlareSolverr endpoint (e.g. http://localhost:8191/v1) to POST the analytics record through a real browser, bypassing Cloudflare's bot check
  -link-id string
        Link UUID to record a click event instead of a profile view
  -proxy string
        Proxy URL for guns.lol requests (e.g. http://user:pass@host:port)
  -username string
        Profile username
```

## Cloudflare

guns.lol puts a Cloudflare bot check on the analytics endpoint that gates by the presence of `cf_clearance`, a client without this cookie gets a Cloudflare challenge.

To get around it, the tool `POST`s the record through a real browser via FlareSolverr:
```shell
docker compose up -d # starts FlareSolverr on 127.0.0.1:8191
```

Then pass `-flaresolverr http://localhost:8191/v1`. Your `-proxy` is forwarded to FlareSolverr (credentials split out, since Chrome rejects inline-auth proxy URLs) so the browser and the tool share one egress IP.

### Examples

On Linux and macOS, the `.exe` extension is not present so remove it.

To add a view to a user:
```shell
./guns-solver.exe -username <username> -capmonster-key ... -proxy ...
```

To add a link click:
```shell
./guns-solver.exe -link-id <link-id>
```

## Botting your views

**Note:** You NEED proxies with a randomized IP for this to work. This is NOT optional since I estimate they rate limit user views to **1 view/per ip/per day**.

Personally I got banned after multiple days **BUT** I've seen people I botted not getting banned after multiple weeks. Their staff is lazy asf

**On Windows:** use this in PowerShell (replace 100 by the number of views you want to add):
```shell
1..100 | ForEach-Object -Parallel { Start-Process "guns-solver.exe" -ArgumentList "-username ....." } -ThrottleLimit 100
```

**On Linux/macOS:** you can use GNU Parallel (100 is the number of repetitions, 40 is the number of parallel processes):
```shell
seq 100 | parallel -j40 ./guns-solver -username .....
```
