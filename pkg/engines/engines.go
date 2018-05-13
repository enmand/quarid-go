package engines

type Type int

const (
	JS Type = iota
)

// An Engine is a plugin runtime
type Engine interface {
	// Load the Engine
	Load() error

	// Load a and "compile" a script into the VM
	Compile(name string, source string) (interface{}, error)

	// Run a previously loaded script in the VM
	Run(name string, args ...interface{}) (string, error)

	Type() Type

	// Return the runtime for the engine, if one is available
	Runtime() interface{}
}
