package js

import (
	"fmt"

	"github.com/dop251/goja"
)

// Set is an implemention of a Set for the JavaScript VM
type Set struct {
	runtime *goja.Runtime
}

// NewSetMapper maps the objects to a map (with int keys)
var NewSetMapper = func(runtime *goja.Runtime) Hasher {
	return func(jsmap map[interface{}]goja.Value, v goja.Value) {
		val, ok := v.Export().([]interface{})
		if !ok {
			panic(runtime.ToValue(fmt.Sprintf("unable to export")))
		}

		for i, v := range val {
			jsmap[i] = runtime.ToValue(v)
		}
	}
}

// NewSet ruterns a new Set
func NewSet(r *goja.Runtime) *Set {
	return &Set{runtime: r}
}

// Enable Set in the JavaScript VM
func (s *Set) Enable(alias ...string) {
	if alias != nil {
		for _, a := range alias {
			s.runtime.Set(a, s.constructor)
		}

		return
	}
	s.runtime.Set("Set", s.constructor)
}

func (s *Set) constructor(call goja.ConstructorCall) *goja.Object {
	jsset := []goja.Value{}
	mapped := iterArgs(&call, s.runtime, NewSetMapper(s.runtime))
	for _, v := range mapped {
		jsset = append(jsset, v)
	}

	call.This.Set("has", func(fcall goja.FunctionCall) goja.Value {
		for _, v := range jsset {
			if v.Equals(fcall.Argument(0)) {
				return s.runtime.ToValue(true)
			}
		}

		return s.runtime.ToValue(false)
	})

	call.This.Set("add", func(fcall goja.FunctionCall) goja.Value {
		assertArgument(fcall.Arguments, 1, s.runtime)
		jsset = append(jsset, fcall.Argument(0))

		return call.This
	})

	return nil
}
