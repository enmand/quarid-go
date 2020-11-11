// go:generated ./vendor/UnnoTed/fileb0x
package js

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/enmand/quarid-go/generated/langsupport"
	"github.com/enmand/quarid-go/internal/engines"
	"github.com/internal/errors"
)

const (
	_compilerScript = "./js/compiler.js"
)

var _compiler = &goja.Program{}

type jsvm struct {
	vm       *goja.Runtime
	compiler goja.Callable
	modules  map[string]*goja.Program
}

func init() {
	_compilerSrc, err := langsupport.ReadFile(_compilerScript)
	if err != nil {
		panic(fmt.Sprintf("unable to load JavaScript compiler (%s): %s", _compilerScript, err))
	}
	_compiler, err = goja.Compile(_compilerScript, string(_compilerSrc), false)
	if err != nil {
		panic(fmt.Sprintf("unable to compile JavaScript copmiler (%s): %s", _compilerScript, err))
	}
}

// New returns a new Otto-based JavaScript virtual machine
func New() (engines.Engine, error) {
	v := &jsvm{
		vm: goja.New(),
	}

	if err := v.initialize(); err != nil {
		return nil, err
	}

	out, err := v.vm.RunProgram(_compiler)
	if err != nil {
		return nil, errors.Wrap(err, "unable to evaluate compiler")
	}

	call, ok := goja.AssertFunction(out)
	if !ok {
		return nil, errors.Wrap(err, "unable to find compiler")
	}

	exports := v.vm.NewObject()
	call(out, nil, exports)

	v.compiler, ok = goja.AssertFunction(exports.Get("_compile"))
	if !ok {
		return nil, errors.Wrap(err, "unable to find compiler function")
	}

	return v, nil
}

// Load the JavaScript engine
func (v *jsvm) Load() error {
	return v.initialize()
}

func (v *jsvm) Type() engines.Type {
	return engines.JS
}

func (v *jsvm) transpile(source string) (string, error) {
	val, err := v.compiler(nil, v.vm.ToValue(source))
	if err != nil {
		return "", errors.Wrap(err, "unable to compile src")
	}

	return val.String(), nil
}

func (v *jsvm) Runtime() interface{} {
	return v.vm
}

func (v *jsvm) Compile(path string, source string) (interface{}, error) {
	if _, ok := v.modules[path]; ok {
		return v.modules[path], nil
	}

	ts, err := v.transpile(source)
	if err != nil {
		return nil, err
	}

	s, err := goja.Compile(path, ts, true)
	if err != nil {
		return nil, fmt.Errorf("Could not compile %s: %s", path, err)
	}
	v.modules[path] = s

	return s, nil
}

func (v *jsvm) Run(path string, args ...interface{}) (string, error) {
	module := v.modules[path]

	val, err := v.vm.RunProgram(module)
	if err != nil {
		return "", fmt.Errorf("Could not run plugin %s: %s", val, err)
	}

	return val.String(), nil
}

func (v *jsvm) initialize() error {
	v.modules = make(map[string]*goja.Program)

	v.vm.Set("require", NewRequire(v.vm).Require)
	NewWeakMap(v.vm).Enable()
	NewMap(v.vm).Enable()
	//v.vm.Set("WeakSet", v.WeakSet)
	NewSet(v.vm).Enable("WeakSet")
	v.vm.Set("Symbol", NewSymbol)
	NewSet(v.vm).Enable()

	return nil
}
