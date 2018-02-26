package js

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/robertkrimen/otto"
	Dbg "github.com/robertkrimen/otto/dbg"
)

// Implements require() in the JavaScript VM.
func RequireFunc(call otto.FunctionCall) otto.Value {

	v, _ := call.Otto.Get(_modpath)

	path, _ := call.Argument(0).ToString()
	fullPath := fmt.Sprintf("%s/%s", v.String(), path)

	if !strings.Contains(path, ".") {
		return _internalRequire(call, fullPath)
	} else {
		return _externalRequire(call, fullPath)
	}
}

func _internalRequire(call otto.FunctionCall, path string) otto.Value {
	requireError(path, "internal")
	return otto.UndefinedValue()
}

func _externalRequire(call otto.FunctionCall, path string) otto.Value {
	d, err := ioutil.ReadFile(path)

	if err != nil {
		// No external module found, let's search our internal path
		return _internalRequire(call, path)
	}

	_, v, err := otto.Run(d)
	if err != nil {
		requireError(path, "external")
		return otto.UndefinedValue()
	}

	return v
}

func requireError(path, requireType string) {
	_, dbf := Dbg.New()
	dbf(
		"%/warn//Module '%s' (%s module search) not found, module not loaded",
		path,
		requireType,
	)
}
