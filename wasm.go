package main

import (
	"context"
	_ "embed"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

//go:embed assets/gpp_gunslol_bg.wasm
var wasmBytes []byte

type WasmResult struct {
	Oo   string
	Seal string
}

// jsHeap mirrors the wasm-bindgen JS heap: a growable slice with an embedded free-list.
type jsHeap struct {
	objects  []any
	freeNext int
}

func newJsHeap() *jsHeap {
	// Indices 0-131 are reserved (0-127 = undefined, 128=undefined, 129=null, 130=true, 131=false).
	return &jsHeap{objects: make([]any, 132), freeNext: 132}
}

func (h *jsHeap) add(obj any) int32 {
	if h.freeNext == len(h.objects) {
		h.objects = append(h.objects, len(h.objects)+1)
	}
	idx := h.freeNext
	if next, ok := h.objects[idx].(int); ok {
		h.freeNext = next
	} else {
		h.freeNext = idx + 1
	}
	for len(h.objects) <= idx {
		h.objects = append(h.objects, nil)
	}
	h.objects[idx] = obj
	return int32(idx)
}

func (h *jsHeap) get(idx int32) any { return h.objects[idx] }

func (h *jsHeap) drop(idx int32) {
	if idx < 132 {
		return
	}
	h.objects[idx] = h.freeNext
	h.freeNext = int(idx)
}

func (h *jsHeap) take(idx int32) any {
	obj := h.get(idx)
	h.drop(idx)
	return obj
}

// SolveWithWasm runs the guns.lol WASM proof-of-work solver via wazero.
// Parameters match the GunsSolver constructor: (o09, difficulty, orgTs, nonce, twoXa).
func SolveWithWasm(ctx context.Context, o09 string, difficulty int, orgTs, nonce, twoXa string) (*WasmResult, error) {
	hp := newJsHeap()

	cacheDir := filepath.Join(os.TempDir(), "guns-solver-wazero")
	cache, err := wazero.NewCompilationCacheWithDir(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("wazero cache: %w", err)
	}
	rt := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().WithCompilationCache(cache))
	defer rt.Close(ctx)

	// readStr reads a UTF-8 string from WASM linear memory.
	readStr := func(mod api.Module, ptr, length int32) string {
		b, ok := mod.Memory().Read(uint32(ptr), uint32(length))
		if !ok {
			return ""
		}
		return string(b)
	}

	hb := rt.NewHostModuleBuilder("wbg")

	// __wbg_new_405e22f390576ce2: () -> i32  — new Object()
	hb.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, _ api.Module, stack []uint64) {
			stack[0] = api.EncodeI32(hp.add(make(map[string]any)))
		}), nil, []api.ValueType{api.ValueTypeI32}).
		Export("__wbg_new_405e22f390576ce2")

	// __wbg_set_3807d5f0bfc24aa7: (i32, i32, i32) -> void  — obj[key] = val
	hb.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, _ api.Module, stack []uint64) {
			obj := hp.get(api.DecodeI32(stack[0])).(map[string]any)
			key := hp.take(api.DecodeI32(stack[1])).(string)
			val := hp.take(api.DecodeI32(stack[2]))
			obj[key] = val
		}), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, nil).
		Export("__wbg_set_3807d5f0bfc24aa7")

	// __wbindgen_number_new: (f64) -> i32
	hb.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, _ api.Module, stack []uint64) {
			stack[0] = api.EncodeI32(hp.add(api.DecodeF64(stack[0])))
		}), []api.ValueType{api.ValueTypeF64}, []api.ValueType{api.ValueTypeI32}).
		Export("__wbindgen_number_new")

	// __wbindgen_object_clone_ref: (i32) -> i32
	hb.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, _ api.Module, stack []uint64) {
			stack[0] = api.EncodeI32(hp.add(hp.get(api.DecodeI32(stack[0]))))
		}), []api.ValueType{api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).
		Export("__wbindgen_object_clone_ref")

	// __wbindgen_object_drop_ref: (i32) -> void
	hb.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, _ api.Module, stack []uint64) {
			hp.drop(api.DecodeI32(stack[0]))
		}), []api.ValueType{api.ValueTypeI32}, nil).
		Export("__wbindgen_object_drop_ref")

	// __wbindgen_rethrow: (i32) -> void
	hb.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, _ api.Module, stack []uint64) {
			panic(fmt.Sprintf("wasm rethrow: %v", hp.take(api.DecodeI32(stack[0]))))
		}), []api.ValueType{api.ValueTypeI32}, nil).
		Export("__wbindgen_rethrow")

	// __wbindgen_string_new: (i32, i32) -> i32
	hb.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, mod api.Module, stack []uint64) {
			s := readStr(mod, api.DecodeI32(stack[0]), api.DecodeI32(stack[1]))
			stack[0] = api.EncodeI32(hp.add(s))
		}), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).
		Export("__wbindgen_string_new")

	// __wbindgen_throw: (i32, i32) -> void
	hb.NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(_ context.Context, mod api.Module, stack []uint64) {
			msg := readStr(mod, api.DecodeI32(stack[0]), api.DecodeI32(stack[1]))
			panic(fmt.Sprintf("wasm throw: %s", msg))
		}), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, nil).
		Export("__wbindgen_throw")

	if _, err := hb.Instantiate(ctx); err != nil {
		return nil, fmt.Errorf("host module: %w", err)
	}

	mod, err := rt.Instantiate(ctx, wasmBytes)
	if err != nil {
		return nil, fmt.Errorf("wasm instantiate: %w", err)
	}

	malloc := mod.ExportedFunction("__wbindgen_export_0")
	addToSP := mod.ExportedFunction("__wbindgen_add_to_stack_pointer")
	newSolver := mod.ExportedFunction("gunssolver_new")
	solvePow := mod.ExportedFunction("gunssolver_solve_pow")

	// allocStr copies s into WASM linear memory and returns (ptr, len).
	allocStr := func(s string) (uint32, uint32) {
		b := []byte(s)
		res, err := malloc.Call(ctx, uint64(len(b)), 1)
		if err != nil {
			panic(fmt.Sprintf("malloc failed: %v", err))
		}
		ptr := uint32(res[0])
		mod.Memory().Write(ptr, b)
		return ptr, uint32(len(b))
	}

	p0, l0 := allocStr(o09)
	p1, l1 := allocStr(orgTs)
	p2, l2 := allocStr(nonce)
	p3, l3 := allocStr(twoXa)

	// gunssolver_new(ptr_o09, len_o09, difficulty, ptr_org_ts, len_org_ts, ptr_nonce, len_nonce, ptr_2xa, len_2xa) -> i32
	solverRes, err := newSolver.Call(ctx,
		uint64(p0), uint64(l0),
		uint64(difficulty),
		uint64(p1), uint64(l1),
		uint64(p2), uint64(l2),
		uint64(p3), uint64(l3),
	)
	if err != nil {
		return nil, fmt.Errorf("gunssolver_new: %w", err)
	}
	solverPtr := solverRes[0]

	// Allocate return slot on WASM stack (16 bytes for 3×i32 + padding).
	spRes, err := addToSP.Call(ctx, api.EncodeI32(-16))
	if err != nil {
		return nil, fmt.Errorf("stack alloc: %w", err)
	}
	retptr := uint32(spRes[0])

	// gunssolver_solve_pow(retptr, solver_ptr) writes 3×i32 at retptr.
	_, err = solvePow.Call(ctx, uint64(retptr), solverPtr)
	addToSP.Call(ctx, api.EncodeI32(16)) // always restore stack
	if err != nil {
		return nil, fmt.Errorf("solve_pow: %w", err)
	}

	mem := mod.Memory()
	r0 := int32(binary.LittleEndian.Uint32(readMem(mem, retptr+0, 4)))
	r1 := int32(binary.LittleEndian.Uint32(readMem(mem, retptr+4, 4)))
	r2 := int32(binary.LittleEndian.Uint32(readMem(mem, retptr+8, 4)))

	if r2 != 0 {
		return nil, fmt.Errorf("solve_pow error: %v", hp.take(r1))
	}

	result, ok := hp.take(r0).(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from solve_pow")
	}
	return &WasmResult{
		Oo:   result["_oo"].(string),
		Seal: result["seal"].(string),
	}, nil
}

func readMem(mem api.Memory, offset, size uint32) []byte {
	b, ok := mem.Read(offset, size)
	if !ok {
		panic(fmt.Sprintf("wasm memory read out of bounds at 0x%x size %d", offset, size))
	}
	return b
}
