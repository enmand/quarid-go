package plugin

import (
	"io/ioutil"

	"github.com/enmand/quarid-go/pkg/adapter"
	"github.com/enmand/quarid-go/pkg/engines"
	"github.com/pkg/errors"
)

// Type defines the runtime to use for the given Plugin
type Type string

var ()

// Plugin is a runnable set of instructions, that can react to an event
type Plugin struct {
	pluginType Type
	path       string
	filters    []adapter.Filter
}

// LoadPlugins a set of plugins from a path
func LoadPlugins(path []string) ([]Plugin, error) {
	for _, dir := range path {
		finfos, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, errors.Wrap(err, "unable to read plugin dir")
		}
		for _, plugin := range finfos {
			if plugin.IsDir() {
				// load plugin dir
			} else {
				// load plugin file
			}
		}
	}

	return nil, nil
}

// New returns a new Plugin that can be used to react to events
func New(t engines.Type, path string, fs []adapter.Filter) *Plugin {
	return &Plugin{
		pluginType: t,
		path:       path,
		filters:    fs,
	}
}

// Run the plugin for the event provided
func (p *Plugin) Run(ev *adapter.Event) error {
	return nil
}
