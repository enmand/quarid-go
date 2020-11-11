package js

import (
	"github.com/dop251/goja"
)

var symbols map[goja.Value]string

// NewSymbol returns the JavaScript Symbol type
var NewSymbol = func(runtime *goja.Runtime) map[string]interface{} {
	return map[string]interface{}{
		"for": func(v goja.FunctionCall) goja.Value {
			assertArgument(v.Arguments, 1, runtime)
			symb := v.Argument(0).String()
			symbols[v.Argument(0)] = symb
			return runtime.ToValue(symb)
		},
		"keyFor": func(v goja.FunctionCall) goja.Value {
			assertArgument(v.Arguments, 1, runtime)
			return runtime.ToValue(symbols[v.Argument(0)])
		},
		"iterator": "@@iterator",
	}
}
