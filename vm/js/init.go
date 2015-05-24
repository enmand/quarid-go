package js

import (
	"fmt"

	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
)

const _modpath = "__modpath_%s___jsvm"

type vm struct {
	vm      *otto.Otto
	modules map[string]interface{}
}

func NewVM() *vm {
	v := &vm{
		vm: otto.New(),
	}

	v.initialize()

	return v
}

func (v *vm) initialize() error {
	v.modules = make(map[string]interface{})

	v.vm.Set("require", RequireFunc)

	return nil
}

func (v *vm) LoadScript(path string, source string) (interface{}, error) {
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

func (v *vm) Run(path string) (string, error) {
	module := v.modules[path]

	v.vm.Set(_modpath, path)

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
