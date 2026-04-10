# guns.lol solver

A solver for the guns.lol WebAssembly script, which is used to record views on profile pages.
Why? Because I was bored and wanted to see how it worked.

This only sends one request per program execution. If you want to bot your profile views then scroll down.

Got banned, but idc check out [Miwa.lol](https://miwa.lol)! It's better than guns.lol :)

Last update to the WASM script: December 20, 2025.

## Usage

- Install [Go](https://go.dev/dl/)
- Download or clone the source code
- Open a Command Prompt/Terminal in the same folder, and run `go build .`
- Then run the program:
```bash
Usage of ./guns-solver.exe:
  -capmonster-key string
        CapMonster API key for Cloudflare Turnstile solving
  -link-id string
        Link UUID to record a click event instead of a profile view
  -proxy string
        Proxy URL for guns.lol requests (e.g. http://user:pass@host:port)
  -username string
        Profile username
```

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

## Bot your views

Personally I got banned after multiple days **but** I've seen people I botted not getting banned after multiple weeks. Their staff is lazy asf

**On Windows:** use this in PowerShell (replace 100 by the number of views you want to add):
```shell
1..100 | ForEach-Object -Parallel { Start-Process "guns-solver.exe" -ArgumentList "-username ....." } -ThrottleLimit 100
```

**On Linux/macOS:** you can use GNU Parallel (100 is the number of repetitions, 40 is the number of parallel processes):
```shell
seq 100 | parallel -j40 ./guns-solver -username .....
```