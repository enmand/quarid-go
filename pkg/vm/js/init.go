package js

import (
	"fmt"

	"github.com/enmand/quarid-go/pkg/vm"
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore" // underscorejs for otto
)

const _modpath = "__modpath_%s___jsvm"

type jsvm struct {
	vm      *otto.Otto
	modules map[string]interface{}
}

// NewVM returns a new Otto-based JavaScript virtual machine
func NewVM() (vm.VM, error) {
	v := &jsvm{
		vm: otto.New(),
	}

	err := v.initialize()
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (v *jsvm) Type() string {
	return vm.JavaScript
}

func (v *jsvm) LoadScript(path string, source string) (interface{}, error) {
	if _, ok := v.modules[path]; ok {
		return nil, fmt.Errorf("Plugin named %s already exists", path)
	}

	s, err := v.vm.Compile("", source)
	if err != nil {
		return nil, fmt.Errorf("Could not compile %s: %s", path, err)
	}
	v.modules[path] = s

	return s, nil
}

func (v *jsvm) Run(path string) (string, error) {
	module := v.modules[path]

	if err := v.vm.Set(_modpath, path); err != nil {
		return "", err
	}

	val, err := v.vm.Run(module.(*otto.Script))
	if err != nil {
		return "", fmt.Errorf("Could not run plugin %s: %s", val, err)
	}

	ret, err := val.ToString()
	if err != nil {
		return "", fmt.Errorf("Could not convert return to response: %s", err)
	}

	return ret, nil
}

func (v *jsvm) initialize() error {
	v.modules = make(map[string]interface{})

	return v.vm.Set("require", RequireFunc)
}
