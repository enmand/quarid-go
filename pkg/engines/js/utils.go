package js

import (
	"fmt"

	"github.com/dop251/goja"
)

type Hasher func(map[interface{}]goja.Value, goja.Value)

func iterArgs(call *goja.ConstructorCall, runtime *goja.Runtime, mapper func(map[interface{}]goja.Value, goja.Value)) map[interface{}]goja.Value {
	jsmap := make(map[interface{}]goja.Value)
	if len(call.Arguments) > 0 {
		jsIter := call.Argument(0).ToObject(runtime)
		iter, ok := goja.AssertFunction(call.Argument(0).ToObject(runtime).Get("next"))
		if ok {
			for {
				v, err := iter(jsIter)
				if err != nil {
					panic(fmt.Sprintf("%s", err.Error()))
				}
				obj := v.ToObject(runtime)
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

				mapper(jsmap, val)
			}
		}
	}

	return jsmap
}

func assertArgument(v []goja.Value, length int, runtime *goja.Runtime) {
	if len(v) < length {
		panic(runtime.ToValue(fmt.Sprintf("no element given. required %d args", length)))
	}
}
