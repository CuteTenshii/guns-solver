# guns.lol solver

A solver for the guns.lol WebAssembly script, which is used to record views on profile pages.

Why? I was bored and wanted to see how it worked. Seriously yes

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