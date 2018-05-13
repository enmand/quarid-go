package js

import (
	"fmt"
	"strings"

	"github.com/kr/pretty"

	"github.com/dop251/goja"
)

// Console imeplements some console functions for interaction with the
// user's console
type Console struct {
	runtime *goja.Runtime
}

// NewConsole returns a new console implementation
func NewConsole(runtime *goja.Runtime) *Console {
	return &Console{runtime: runtime}
}

// Enable the Console in the runtime
func (c *Console) Enable() {
	c.runtime.Set("console", console)
}

var console = map[string]func(goja.FunctionCall) goja.Value{
	"log": func(fcall goja.FunctionCall) goja.Value {
		log := ""
		for _, a := range fcall.Arguments {
			log += fmt.Sprintf("%s ", a)
		}

		fmt.Printf("%s", strings.Trim(log, " "))

		return nil
	},

	"dir": func(fcall goja.FunctionCall) goja.Value {
		log := ""
		for _, a := range fcall.Arguments {
			exp := a.Export()
			log += fmt.Sprintf("%# v", pretty.Formatter(exp))
		}

		fmt.Printf("%s", strings.Trim(log, " "))

		return nil
	},
}
