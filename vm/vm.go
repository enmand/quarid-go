package vm

type Runnable interface{}

type FunctionCall interface{}
type Value interface{}

type Handler struct {
	Name     string
	Function func(call FunctionCall) Value
}

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
	// Initialize the virtual machine
	Initialize() error

	// Load a and "compile" a script into the VM
	LoadScript(name string, source string) (Runnable, error)

	// Run a previously loaded script in the VM
	Run(name string) (string, error)
}

func New(typ string) VM {
	switch typ {
	case JS:
		return newJsVm()
	}

	return nil
}
