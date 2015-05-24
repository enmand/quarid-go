package vm

import "github.com/enmand/quarid-go/vm/js"

const (
	JS = "js"

	// TODO
	JSv8    = "jsv8"
	PY      = "py"
	LUA     = "lua"
	PROLOG  = "prolog"
	PHP     = "php"
	HASKELL = "haskell"
)

type VM interface {
	// Load a and "compile" a script into the VM
	LoadScript(name string, source string) (interface{}, error)

	// Run a previously loaded script in the VM
	Run(name string) (string, error)
}

func New(typ string) VM {
	switch typ {
	case JS:
		return js.NewVM()
	}

	return nil
}
