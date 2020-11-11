package js

import (
	"fmt"
	"unsafe"

	"github.com/dop251/goja"
	"github.com/mitchellh/hashstructure"
)

// WeakMap is an implementation of a WeakMap for the JavaScrpit VM
type WeakMap struct {
	runtime *goja.Runtime
}

// WeakMapMapper can hash objects going into the WeakMap
var WeakMapMapper = func(runtime *goja.Runtime) Hasher {
	return func(jsmap map[interface{}]goja.Value, v goja.Value) {
		kvp, ok := v.Export().([]interface{})
		if len(kvp) < 1 || !ok {
			panic("need key and value")
		}

		hash, err := hashstructure.Hash(kvp[0], nil)
		if err != nil {
			panic(fmt.Sprintf("unable to hash object: %s", err))
		}
		jsmap[hash] = runtime.ToValue(kvp[1])
	}
}

// NewWeakMap returns a new WeakMap
func NewWeakMap(r *goja.Runtime) *WeakMap {
	return &WeakMap{runtime: r}
}

// Enable WeakMap in the JavaScript runtime
func (wm *WeakMap) Enable() {
	wm.runtime.Set("WeakMap", wm.constructor)
}

func (wm *WeakMap) constructor(call goja.ConstructorCall) *goja.Object {
	jsmap := make(map[uint64]goja.Value)
	mapped := iterArgs(&call, wm.runtime, WeakMapMapper(wm.runtime))
	for k, v := range mapped {
		jsmap[k.(uint64)] = v
	}

	call.This.Set("clear", func(fcall goja.FunctionCall) goja.Value {
		jsmap = make(map[uint64]goja.Value)
		return nil
	})

	call.This.Set("delete", func(fcall goja.FunctionCall) goja.Value {
		assertArgument(fcall.Arguments, 1, wm.runtime)
		hash := hashArguments(fcall.Argument(0), wm.runtime)
		delete(jsmap, hash)
		return wm.runtime.ToValue(true)
	})

	call.This.Set("get", func(fcall goja.FunctionCall) goja.Value {
		assertArgument(fcall.Arguments, 1, wm.runtime)
		hash := hashArguments(fcall.Argument(0), wm.runtime)
		return jsmap[hash]

	})

	call.This.Set("has", func(fcall goja.FunctionCall) goja.Value {
		assertArgument(fcall.Arguments, 0, wm.runtime)
		hash := hashArguments(fcall.Argument(0), wm.runtime)
		_, ok := jsmap[hash]
		return wm.runtime.ToValue(ok)
	})

	call.This.Set("set", func(fcall goja.FunctionCall) goja.Value {
		assertArgument(fcall.Arguments, 2, wm.runtime)
		hash := hashArguments(fcall.Argument(0), wm.runtime)

		jsmap[hash] = fcall.Argument(1)
		return nil
	})

	return nil
}

func hashArguments(v goja.Value, runtime *goja.Runtime) uint64 {
	key := v.ToObject(runtime).Export()
	hash, err := hashstructure.Hash(key, nil)
	if err != nil {
		hash = *(*uint64)(unsafe.Pointer(&key))
	}

	return hash
}
