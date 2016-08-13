package plugin

import "github.com/enmand/quarid-go/pkg/vm"

// Plugin defines the interface for running or loading plugins from different
// runtime environments
type Plugin interface {
	Load(i map[string]vm.VM) error
	Compile() error
	Run() error
}

// NewPlugin returns a new plugin at the given path
func NewPlugin(name, path string) Plugin {
	p := &plugin{
		Name: name,
		path: path,
	}

	return p
}
