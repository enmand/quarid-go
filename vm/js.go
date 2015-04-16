package vm

import (
	"fmt"

	"github.com/enmand/quarid-go/vm/js"
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
)

type vm struct {
	vm      *otto.Otto
	modules map[string]Runnable
}

func newJsVm() *vm {
	return &vm{
		vm: otto.New(),
	}
}

func (v *vm) Initialize() error {
	v.modules = make(map[string]Runnable)

	v.vm.Set("require", js.Require)

	return nil
}

func (v *vm) LoadScript(name string, source string) (Runnable, error) {
	if _, ok := v.modules[name]; ok {
		return nil, fmt.Errorf("Plugin named %s already exists", name)
	}

	s, err := v.vm.Compile("", source)
	if err != nil {
		return nil, fmt.Errorf("Could not compile %s: %s", name, err)
	}
	v.modules[name] = s

	return s, nil
}

func (v *vm) Run(name string) (string, error) {
	module := v.modules[name]

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
