package vm

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
)

type js struct {
	vm      *otto.Otto
	modules map[string]Runnable
}

type jsHandler struct {
	Name     string
	Function func(call otto.FunctionCall) otto.Value
}

func newJsVm() *js {
	return &js{
		vm: otto.New(),
	}
}

func (v *js) Initialize() error {
	v.modules = make(map[string]Runnable)
	return nil
}

func (v *js) LoadScript(name string, source string) (Runnable, error) {
	if _, ok := v.modules[name]; ok {
		return nil, fmt.Errorf("Plugin named %s already exists", name)
	}

	s, err := v.vm.Compile("", source)
	if err != nil {
		return nil, fmt.Errorf("Could not compile %s: %s", name, err)
	}
	v.modules[name] = s

	log.Errorf("%#v", s)
	return s, nil
}

func (v *js) Run(name string) (string, error) {
	module := v.modules[name]

	val, err := v.vm.Run(module.(*otto.Script))
	if err != nil {
		return "", fmt.Errorf("Could not run plugin %s: %s", val, err)
	}

	ret, err := val.ToString()
	if err != nil {
		return "", fmt.Errorf("Could not convert return to response: %s", err)
	}
	fmt.Printf(ret)

	return ret, nil
}
