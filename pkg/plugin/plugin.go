package plugin

import "github.com/enmand/quarid-go/pkg/adapter"

// Type defines the runtime to use for the given Plugin
type Type string

var (
	// TypePython defines Python-based modules
	TypePython = Type("python")
)

// Plugin is a runnable set of instructions, that can react to an event
type Plugin interface {
	Load() error
	Run(ev *adapter.Event) error
}

// New returns a new Plugin that can be used to react to events
func New(t Type, path string, fs []adapter.Filter) Plugin {
	return &plugin{
		pluginType: t,
		path:       path,
		filters:    fs,
	}
}

type plugin struct {
	pluginType Type
	path       string
	filters    []adapter.Filter
}

func (p *plugin) Run(ev *adapter.Event) error {
	switch p.pluginType {
	case TypePython:
		// Do pythony stuff here
	}

	return nil
}
