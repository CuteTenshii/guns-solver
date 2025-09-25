# Guns.lol Solver

A solver for the guns.lol WebAssembly script, which is used to record views on profile pages. Was cool to reverse engineer it :)

**If you know anything about WASM or C, please help the project by contributing!**

## How it works - Reverse Engineering

This section explains how guns.lol's WebAssembly script works and how the solver interacts with it.

### Server-Side script

When you load a profile page, a server-side script which loads another one is called (see [Worker script](#worker-script) below).

Here the first script is also interesting because it contains some data which is used to construct the view recording payload.

### Worker script

A web worker is created with the following script:
```js
import init, { GunsSolver } from 'https://assets.guns.lol/pkg/_gunslolpow.js';
self.onmessage = async function (event) {
    const { ps, d, c, _n, org_ts } = event.data;
    await init();
    const _get_s = new GunsSolver(ps, c, parseInt(d), false, _n);
    const ts = Math.floor(Date.now() / 1000);
    const _tsoff = ts - Number(org_ts);
    const tsp = ts - _tsoff;
    const _res = await _get_s.solve_pow();
    if (_res !== null) {
        self.postMessage({ _res, ps, tsp, _n, c });
    }
};
```

`init()` initializes the WebAssembly module.
Then, a `worker.postMessage` is called with an object that has the following fields:
- `ps`: public salt
- `d`: difficulty
- `c`: challenge
- `_n`: nonce
- `org_ts`: original timestamp?

The worker then creates a new `GunsSolver` instance with the parameters and calls its `solve_pow` method, which returns a result `_res`.
Finally, it sends back an object containing `_res`, `ps`, `tsp`, `_n`, and `c`.

```js
worker.onmessage = async function (event) {
    const { _res, ps, tsp, _n, c } = event.data;
    if (_res !== undefined) {
        await getResult(_res);
    }
};
```

The main script listens for messages from the worker and calls `getResult` with `_res` when it receives a message (the rest is not used). `getResult` is defined in a Next.js chunk (see below).

### `_gunslolpow.js`

This script fetches a WebAssembly module (written in Rust) from `https://assets.guns.lol/pkg/_gunslolpow_bg.wasm` and instantiates it.
Also, it exports the `GunsSolver` class, used in the worker script (see above).

#### `GunsSolver` class

It's interesting to look at the `GunsSolver` class and its constructor. It takes the following parameters:
- `public_salt` (string)
- `challenge` (string)
- `difficulty` (number)
- `numeric` (boolean)
- `nonce` (string)

```js
class GunsSolver {
    constructor(public_salt, challenge, difficulty, numeric, nonce) {
        const ptr0 = passStringToWasm0(public_salt, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len0 = WASM_VECTOR_LEN;
        const ptr1 = passStringToWasm0(challenge, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len1 = WASM_VECTOR_LEN;
        const ptr2 = passStringToWasm0(nonce, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
        const len2 = WASM_VECTOR_LEN;
        const ret = wasm.gunssolver_new(ptr0, len0, ptr1, len1, difficulty, numeric, ptr2, len2);
        this.__wbg_ptr = ret >>> 0;
        GunsSolverFinalization.register(this, this.__wbg_ptr, this);
        return this;
    }
}
```

The constructor sends the parameters to the WebAssembly module with `passStringToWasm0()` for the salt, challenge, and nonce to be converted to pointers.

### WASM reverse engineering

The WebAssembly module has been converted to C using `wasm2c`, which made it easier to understand how the PoW challenge is solved. See the `assets` directory for the generated C code and their header files.

It has also been converted to WAT using `wasm2wat`, which is easier to read than raw WebAssembly, but still not as easy as C.

### Payload construction

Request payload is constructed in the following way:
```json
{
  "_t": "turnstile response token",
  "_c": "challenge",
  "_ps": "base64 of public salt",
  "_ts": "calculated timestamp",
  "_s": {
    "_s": "wasm solver result (5 chars [a-f0-9])",
    "n": "nonce",
    "ps": "public salt (same as _ps but not base64)"
  },
  "__hcm": {
    "_s": "wasm solver result",
    "_": "sha-256 hash, idk what it's for",
    "_2": "another sha-256 hash, idk what it's for either",
    "__meta": {
      // the values here never change, they're always 5x "64" and 1x "5"
      "maybe hex?": 64,
      "maybe hex?": 64,
      "maybe hex?": 64,
      "maybe hex?": 64,
      "maybe hex?": 64,
      "maybe hex?": 5,
    },
    // probably not in this order, this is just to better illustrate the structure
    "random string 1": "sha-256 hash of: public salt + challenge",
    "random string 2": "sha-256 hash of: idk",
    "random string 3": "sha-256 hash of: idk",
    "random string 4": "sha-256 hash of: idk",
    "random string 5": "sha-256 hash of: idk",
    "random string 6": "sha-256 hash of: idk",
  },
  "username": "profile username",
  "deviceType": "device type: desktop, tablet, mobile",
  "event": "view",
  "linkId": null,
  "referrer": "document.referrer"
}
```

**Keys which have an unknown purpose needs to be reverse engineered.**
- `_t`: Turnstile response token, obtained from the Turnstile widget on the profile page.
- `_c`: Challenge, obtained from the server-side script.
- `_ps`: Base64-encoded public salt, obtained from the server-side script.
- `_ts`: Calculated timestamp. In fact, it's just the `org_ts` value sent to the worker.
- `_s`: An object containing:
  - `_s`: The result from the WebAssembly solver (5 characters long, hexadecimal).
  - `n`: Nonce, obtained from the server-side script.
  - `ps`: Public salt (same as `_ps` but not base64-encoded).
- `__hcm`: An object containing:
  - `_s`: The result from the WebAssembly solver (same as in `_s._s`).
  - `_`: A SHA-256 hash (purpose unknown).
  - `_2`: Another SHA-256 hash (purpose unknown).
  - `__meta`: An object with constant values (always five 64s and one 5).
  - Six random strings as keys, each containing a SHA-256 hash:
    - One is the hash of `public salt + challenge`.
    - The other five hashes have unknown purposes.
- `username`: The profile username.
- `deviceType`: The type of device (e.g., desktop, tablet, mobile).
- `event`: The event type, always "view".
- `linkId`: Always null.
- `referrer`: The document referrer.

This payload is sent to `POST https://guns.lol/api/analytics/record` to record the view.