package plugin

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/enmand/quarid-go/pkg/adapter"
	"github.com/enmand/quarid-go/pkg/engines"
	"github.com/pkg/errors"
	"strings"
)

// Plugin is a runnable set of instructions, that can react to an event
type Plugin struct {
	pluginType engines.Type
	plugin     plugin
	path       string
	filters    []adapter.Filter
}

type plugin interface {
	Filters() ([]adapter.Filter, error)
}

// LoadPlugins a set of plugins from a path
func LoadPlugins(path []string) ([]*Plugin, error) {
	ps := []*Plugin{}

	for _, dir := range path {
		finfos, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, errors.Wrap(err, "unable to read plugin dir")
		}
		for _, plugin := range finfos {
			path := fmt.Sprintf("%s/%s", dir, plugin.Name())
			if plugin.IsDir() {
				// load plugin dir
			} else {
				p, err := Load(path)
				if err != nil {
					return nil, err
				}
				ps = append(ps, p)
			}
		}
	}

	return ps, nil
}

// Load loads a plugin from a file source
func Load(path string) (*Plugin, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open file")
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't read file")
	}

	if strings.HasSuffix(path, ".js") {
		jsplugin, err := newJavaScript(path, string(b))
		if err != nil {
			return nil, err
		}

		f, err := jsplugin.Filters()
		if err != nil {
			return nil, errors.Wrap(err, "unable to load javascript plugin filters")
		}

		return &Plugin{
			plugin:     jsplugin,
			pluginType: engines.JS,
			path:       path,
			filters:    f,
		}, nil
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
