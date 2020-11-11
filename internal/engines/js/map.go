package js

import (
	"unsafe"

	"github.com/dop251/goja"
)

// TODO: Just use WeakMap here, or expand this implentation

// Map is a partial Map imlementation for Babel
type Map struct {
	runtime *goja.Runtime
}

// NewMap returns a new Map object for a goja.Runtime
func NewMap(r *goja.Runtime) *Map {
	return &Map{runtime: r}
}

// MapMapper can hash objects going into the WeakMap
var MapMapper = func(runtime *goja.Runtime) func(map[interface{}]goja.Value, goja.Value) {
	return func(jsmap map[interface{}]goja.Value, v goja.Value) {
		kvp, ok := v.Export().([]interface{})
		if len(kvp) < 1 || !ok {
			panic("need key and value")
		}

		jsmap[kvp[0]] = runtime.ToValue(kvp[1])
	}
}

// Enable the Map in the JavaScript runtime
func (m *Map) Enable() {
	m.runtime.Set("Map", m.constructor)
}

func (m *Map) constructor(call goja.ConstructorCall) *goja.Object {
	jsmap := iterArgs(&call, m.runtime, MapMapper(m.runtime))

	call.This.Set("get", func(fcall goja.FunctionCall) goja.Value {
		assertArgument(fcall.Arguments, 1, m.runtime)
		val := fcall.Argument(0).Export()
		return jsmap[*(*uint64)(unsafe.Pointer(&val))]
	})

	call.This.Set("set", func(fcall goja.FunctionCall) goja.Value {
		assertArgument(fcall.Arguments, 2, m.runtime)

		val := fcall.Argument(0).Export()
		jsmap[*(*uint64)(unsafe.Pointer(&val))] = fcall.Argument(0)

		return nil
	})

	call.This.Set("has", func(fcall goja.FunctionCall) goja.Value {
		assertArgument(fcall.Arguments, 2, m.runtime)
		val := fcall.Argument(0).Export()
		_, ok := jsmap[*(*uint64)(unsafe.Pointer(&val))]
		return m.runtime.ToValue(ok)
	})

	return nil
}
