package vm

const (
	// JS is the JavaScript VM
	JS = "js"

	// TODO
	// PY is the Python VM
	PY = "py"
	// LUA is the Lua VM
	LUA = "lua"

	// PROLOG is the Prolog VM
	PROLOG = "prolog"

	// PHP is the PHP VM
	PHP = "php"

	// Haskell is the Haskell VM
	HASKELL = "haskell"
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
