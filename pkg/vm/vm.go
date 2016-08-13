package vm

const (
	// JavaScript runtime VM
	JavaScript = "js"
)

// A VM is a language-based virtual machine for running loading, and
// running code
type VM interface {
	// Load a and "compile" a script into the VM
	LoadScript(name string, source string) (interface{}, error)

	// Run a previously loaded script in the VM
	Run(name string) (string, error)

	Type() string
}
