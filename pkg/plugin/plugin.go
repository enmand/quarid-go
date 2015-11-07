package plugin

import "github.com/enmand/quarid-go/vm"

type Plugin interface {
	Load(i map[string]vm.VM) error
	Run() error
}

func NewPlugin(name, path string) *plugin {
	p := &plugin{
		Name: name,
		path: path,
	}

	return p
}
