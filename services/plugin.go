package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/enmand/quarid-go/vm"
)

type Plugin interface {
	Load(b Bot) error
}

func NewPlugin(name, path string) *plugin {
	p := &plugin{
		Name: name,
		path: path,
	}

	return p
}

type plugin struct {
	path string
	vm   vm.VM

	Name          string      `json:"name"`
	VM            string      `json:"vm"`
	Main          string      `json:"main"`
	Configuration interface{} `json:"configuration"`
}

func (p *plugin) Load(b Bot) error {
	pp, err := p.pluginConfig()
	if err != nil {
		return err
	}
	p = pp.(*plugin) // Set our plugin configuration our loaded config
	log.Infof("Loading plugin: %s (in VM: %s)", p.Name, p.VM)

	return nil
}

// Load JSON-based plugin configuration from plugin.json
func (p *plugin) pluginConfig() (Plugin, error) {
	cf := fmt.Sprintf("%s/plugin.json", p.path)

	cb, err := ioutil.ReadFile(cf)
	if err != nil {
		return nil, fmt.Errorf("Cannot load plugin '%s': %s", p.Name, err)
	}

	var pp Plugin
	err = json.Unmarshal(cb, pp)
	if err != nil {
		return nil, fmt.Errorf(
			"Cannot load plugin '%s' configuration: %s",
			p.Name,
			err,
		)
	}

	if p.Name != pp.(*plugin).Name {
		log.Warningf(`
			Sanity check! %s has a different name than it was configured
			for ('%s')
		`, p.Name, pp.(*plugin).Name)
	}

	return pp, err
}

// Compile our plugin, using the VM given
func (p *plugin) compile() error {
	return nil
}
