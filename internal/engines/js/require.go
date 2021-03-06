package js

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/dop251/goja"
	"github.com/enmand/quarid-go/generated/langsupport"
)

// Require implements the 'require(module)' pattern from NodeJS
type Require struct {
	modules map[string]*goja.Program
	runtime *goja.Runtime

	l sync.Mutex
}

// NewRequire returns a new Require object
func NewRequire(runtime *goja.Runtime) *Require {
	return &Require{
		modules: make(map[string]*goja.Program),
		runtime: runtime,
	}
}

// Require implements require() in the JavaScript VM.
func (r *Require) Require(call goja.FunctionCall) goja.Value {
	pathVal := call.Argument(0)
	if goja.IsNull(pathVal) && goja.IsUndefined(pathVal) {
		return r.runtime.NewGoError(errors.New("no path given"))
	}

	path := pathVal.String()
	if strings.HasPrefix(path, "./") {
		return r._externalRequire(call, path)
	}

	return r._internalRequire(call, path)

}

func (r *Require) _externalRequire(call goja.FunctionCall, path string) goja.Value {
	r.requireError(path, "external", errors.New("not implemented"))
	return goja.Undefined()
}

func (r *Require) _internalRequire(call goja.FunctionCall, path string) goja.Value {
	script, err := readBoxedFile(path)
	if err != nil || script == nil {
		// No external module found, let's search our internal path
		r.requireError(path, "internal", err)
		return r._internalRequire(call, path)

	}
	source := "(function(module, exports) {" + string(*script) + "\n})"
	p, err := goja.Compile(path, source, false)
	if err != nil {
		r.requireError(path, "internal", fmt.Errorf("unable to compile %s", path))
		return nil
	}

	out, err := r.runtime.RunProgram(p)
	if err != nil {
		r.requireError(path, "internal", errors.New("unable to run program"))
	}

	reqCall, ok := goja.AssertFunction(out)
	if ok != true {
		r.requireError(path, "internal", errors.New("unable to get transpiled function"))
		return nil
	}

	exports := r.runtime.NewObject()
	_, err = reqCall(out, nil, exports)

	return exports
}

func (r *Require) requireError(path, requireType string, err error) {
	fmt.Printf(" warn/Module '%s' (%s module search) module not loaded: %s\n", path, requireType, err)
}

func readBoxedFile(path string) (*string, error) {
	var fullPath string
	if strings.HasPrefix(path, "!") {
		fullPath = fmt.Sprintf("js/%s", path[1:])
	}

	if strings.HasPrefix(path, "@") {
		fullPath = fmt.Sprintf("js/node_modules/%s", path)
	}

	d, err := langsupport.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	script := string(d)
	return &script, nil
}
