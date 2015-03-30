package vm

type Type int
type Runnable interface{}

type FunctionCall interface{}
type Value interface{}

type Handler struct {
	Name     string
	Function func(call FunctionCall) Value
}

const (
	JS   = iota
	JSv8 // developmental... more so than everything else
	PY
	LUA
	PROLOG
	PHP
)

type VM interface {
	// Initialize the virtual machine
	Initialize() error

	// Load a and "compile" a script into the VM
	LoadScript(name string, source string) (Runnable, error)

	// Run a previously loaded script in the VM
	Run(name string) (string, error)
}

func New(typ Type) VM {
	switch typ {
	case JS:
		return newJsVm()
	}

	return nil
}
