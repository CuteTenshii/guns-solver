/*
Usage in _gunslolpow.js:

const ptr0 = passStringToWasm0(public_salt, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
const len0 = WASM_VECTOR_LEN;
const ptr1 = passStringToWasm0(challenge, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
const len1 = WASM_VECTOR_LEN;
const ptr2 = passStringToWasm0(nonce, wasm.__wbindgen_malloc, wasm.__wbindgen_realloc);
const len2 = WASM_VECTOR_LEN;

const ret = wasm.gunssolver_new(ptr0, len0, ptr1, len1, difficulty, numeric, ptr2, len2); <-- Usage here
 */

u32 w2c_guns_gunssolver_new_0(w2c_guns* instance, u32 public_salt_ptr, u32 public_salt_len, u32 challenge_ptr, u32 challenge_len, u32 difficulty, u32 numeric, u32 nonce_ptr, u32 nonce_len) {
  u32 solver_obj = 0;
  FUNC_PROLOGUE;
  u32 var_i0, var_i1, var_i2;
  u64 var_j1;
  
  var_i0 = 1063025u;
  var_i0 = i32_load8_u_default32(&instance->w2c_memory, (u64)(var_i0));
  var_i0 = 56u;  // Allocate 56 bytes for solver object
  var_i1 = 4u;   // Alignment: 4 bytes
  var_i0 = w2c_guns_f203(instance, var_i0, var_i1);
  solver_obj = var_i0;
  var_i0 = !(var_i0);
  if (var_i0) {
    // Allocation failed - call error handler
    var_i0 = 4u;
    var_i1 = 56u;
    w2c_guns_f237(instance, var_i0, var_i1);
    UNREACHABLE;
  }
  
  // Initialize solver object fields
  var_i0 = solver_obj;
  var_i1 = difficulty;
  i32_store_default32(&instance->w2c_memory, (u64)(var_i0) + 48, var_i1);  // +48: difficulty
  var_i0 = solver_obj;
  var_i1 = nonce_len;
  i32_store_default32(&instance->w2c_memory, (u64)(var_i0) + 44, var_i1);  // +44: nonce_len
  var_i0 = solver_obj;
  var_i1 = nonce_ptr;
  i32_store_default32(&instance->w2c_memory, (u64)(var_i0) + 40, var_i1);  // +40: nonce_ptr
  var_i0 = solver_obj;
  var_i1 = nonce_len;
  i32_store_default32(&instance->w2c_memory, (u64)(var_i0) + 36, var_i1);  // +36: nonce_len (backup)
  var_i0 = solver_obj;
  var_i1 = challenge_len;
  i32_store_default32(&instance->w2c_memory, (u64)(var_i0) + 32, var_i1);  // +32: challenge_len
  var_i0 = solver_obj;
  var_i1 = challenge_ptr;
  i32_store_default32(&instance->w2c_memory, (u64)(var_i0) + 28, var_i1);  // +28: challenge_ptr
  var_i0 = solver_obj;
  var_i1 = challenge_len;
  i32_store_default32(&instance->w2c_memory, (u64)(var_i0) + 24, var_i1);  // +24: challenge_len (backup)
  var_i0 = solver_obj;
  var_i1 = public_salt_len;
  i32_store_default32(&instance->w2c_memory, (u64)(var_i0) + 20, var_i1);  // +20: public_salt_len
  var_i0 = solver_obj;
  var_i1 = public_salt_ptr;
  i32_store_default32(&instance->w2c_memory, (u64)(var_i0) + 16, var_i1);  // +16: public_salt_ptr
  var_i0 = solver_obj;
  var_i1 = public_salt_len;
  i32_store_default32(&instance->w2c_memory, (u64)(var_i0) + 12, var_i1);  // +12: public_salt_len (backup)
  var_i0 = solver_obj;
  var_i1 = 0u;
  i32_store_default32(&instance->w2c_memory, (u64)(var_i0) + 8, var_i1);   // +8: current_iteration (initialized to 0)
  var_i0 = solver_obj;
  var_j1 = 4294967297ull;  // 0x100000001 = reference_count(1) + flags(1)
  i64_store_default32(&instance->w2c_memory, (u64)(var_i0), var_j1);       // +0: ref_count + flags
  var_i0 = solver_obj;
  var_i1 = numeric;
  var_i2 = 0u;
  var_i1 = var_i1 != var_i2;
  i32_store8_default32(&instance->w2c_memory, (u64)(var_i0) + 52, var_i1); // +52: numeric_mode (boolean)
  
  // Return pointer to offset +8 (current_iteration field)
  var_i0 = solver_obj;
  var_i1 = 8u;
  var_i0 += var_i1;
  FUNC_EPILOGUE;
  return var_i0;
}
