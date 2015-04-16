package js

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/robertkrimen/otto"
	Dbg "github.com/robertkrimen/otto/dbg"
)

// Implements require() in the JavaScript VM.
func Require(call otto.FunctionCall) otto.Value {
	path, _ := call.Argument(0).ToString()
	if !strings.Contains(path, ".") {
		return _internalRequire(call, path)
	} else {
		return _externalRequire(call, path)
	}
}

func _internalRequire(call otto.FunctionCall, path string) otto.Value {
	requireError(path)
	return otto.Value{}
}

func _externalRequire(call otto.FunctionCall, path string) otto.Value {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		// No external module found, let's search our internal path
		return _internalRequire(call, path)
	}

	_, v, err := otto.Run(d)
	if err != nil {
		ef, _ := otto.ToValue(fmt.Errorf("Could not compile module %s", path))
		return ef
	}

	return v
}

func requireError(path string) {
	_, dbf := Dbg.New()
	dbf("%/panic//Module '%s' not found", path)
}
