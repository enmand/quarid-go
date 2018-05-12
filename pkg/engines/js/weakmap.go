package js

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/mitchellh/hashstructure"
)

func (v *jsvm) WeakMap(call goja.ConstructorCall) *goja.Object {
	return nil
}

// WeakMap is an implementation of a WeakMap for the JavaScrpit VM
type WeakMap struct {
	runtime *goja.Runtime
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

	if len(call.Arguments) > 0 {
		jsIter := call.Argument(0).ToObject(wm.runtime)
		iter, ok := goja.AssertFunction(call.Argument(0).ToObject(wm.runtime).Get("next"))
		if ok {
			for {
				v, err := iter(jsIter)
				if err != nil {
					panic(fmt.Sprintf("%s", err.Error()))
				}
				obj := v.ToObject(wm.runtime)
				if goja.IsUndefined(v) && obj.Get("done").ToBoolean() == true {
					break
				}

				val := obj.Get("value")
				if goja.IsNull(val) || goja.IsUndefined(val) {
					panic("no key values given")
				}

				if val == nil {
					break
				}

				kvp, ok := val.Export().([]interface{})
				if len(kvp) < 1 || !ok {
					panic("need key and value")
				}

				hash, err := hashstructure.Hash(kvp[0], nil)
				if err != nil {
					panic(fmt.Sprintf("unable to hash object: %s", err))
				}
				jsmap[hash] = wm.runtime.ToValue(kvp[1])
			}
		}
	}

	call.This.Set("clear", func(fcall goja.FunctionCall) goja.Value {
		jsmap = make(map[uint64]goja.Value)
		return nil
	})

	call.This.Set("delete", func(fcall goja.FunctionCall) goja.Value {
		hash := assertArgument(fcall.Arguments, 1, wm.runtime)
		delete(jsmap, hash)
		return wm.runtime.ToValue(true)
	})

	call.This.Set("get", func(fcall goja.FunctionCall) goja.Value {
		hash := assertArgument(fcall.Arguments, 1, wm.runtime)
		return jsmap[hash]

	})

	call.This.Set("has", func(fcall goja.FunctionCall) goja.Value {
		hash := assertArgument(fcall.Arguments, 1, wm.runtime)
		_, ok := jsmap[hash]
		return wm.runtime.ToValue(ok)
	})

	call.This.Set("set", func(fcall goja.FunctionCall) goja.Value {
		hash := assertArgument(fcall.Arguments, 2, wm.runtime)

		jsmap[hash] = fcall.Argument(1)
		return nil
	})

	return nil
}

func assertArgument(v []goja.Value, length int, runtime *goja.Runtime) uint64 {
	if len(v) < length {
		panic(fmt.Sprintf("no element given. required %d args", length))
	}

	key := v[0].ToObject(runtime).Export()
	hash, err := hashstructure.Hash(key, nil)
	if err != nil {
		panic(fmt.Sprintf("unable to hash object: %s", err))
	}

	return hash
}
