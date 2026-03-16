# guns.lol solver

A solver for the guns.lol WebAssembly script, which is used to record views on profile pages.

Why? I was bored and wanted to see how it worked. Seriously yes

This only send one request per program execution, if you want to bot your profile views then scroll down.

Got banned, but idc check out [Miwa.lol](https://miwa.lol)! It's better than guns.lol :)

Last update to the WASM script: December 20, 2025.

## Usage

Requires Go.

- Download or clone the source code
- Build it with `go build .`
- Run it:
```bash
Usage of ./guns-solver:
  -capmonster-key string
        CapMonster API key for Turnstile solving
  -link-id string
        Link UUID to record a click event instead of a profile view
  -proxy string
        Proxy URL for guns.lol requests (e.g. http://user:pass@host:port)
  -username string
        Profile username
```

## Bot your views

Personally I got banned after multiple days. Their mods are lazy asf

On Windows use this in PowerShell (replace 100 by the number of views you want to add):
```shell
1..100 | ForEach-Object -Parallel { Start-Process "guns-solver.exe" -ArgumentList "-username ....." } -ThrottleLimit 100
```

On Linux/macOS (100 is the number of repetitions, 40 is the number of parallel processes):
```shell
seq 100 | parallel -j40 ./guns-solver -username .....
```