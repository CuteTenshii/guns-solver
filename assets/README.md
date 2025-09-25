# Assets

This directory contains various stuff which has been used to reverse engineer how guns.lol uses their WebAssembly module.

## Files

- `guns.wasm` - The WebAssembly module used by guns.lol, in raw.
- `guns.wat` - The WebAssembly module in assembly text format, generated using `wasm2wat`.
- `guns.c` - The C code generated using `wasm2c`.
- `guns.h` - The header file generated using `wasm2c`.
- `w2c_guns_gunssolver_solve_pow_0.c` - The C code for the `w2c_guns_gunssolver_solve_pow_0` function, extracted from `guns.c`. Variables have been renamed for clarity.
- `next.js` - The JS script found in profile pages.
- `4999-fa0fe6ccdca6b8ab.js` - The Next.js chunk containing the view recording logic.
- `_gunslolpow.js` - The JS script loader by the `next.js` script.

### C Files

Those files have been extracted from `guns.c` - generated using `wasm2c` - and had their variables renamed for clarity.

- `w2c_guns_gunssolver_new_0.c` - The C code for the `w2c_guns_gunssolver_new_0` function, which creates a new `gunssolver` object.
- `w2c_guns_gunssolver_solve_pow_0.c` - The C code for the `w2c_guns_gunssolver_solve_pow_0` function, which solves the proof-of-work challenge.
