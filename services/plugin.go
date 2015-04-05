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
	path string `json:"-"`
	vm   vm.VM  `json:"-"`

	Name          string      `json:"name"`
	VM            string      `json:"vm"`
	Main          string      `json:"main"`
	Configuration interface{} `json:"configuration"`
}

func (p *plugin) Load(b Bot) error {
	pp, err := p.pluginConfig(b)
	if err != nil {
		return err
	}
	p = pp.(*plugin) // Set our plugin configuration our loaded config
	log.Infof("Loading plugin: %s (in VM: %s)", p.Name, p.VM)

	return nil
}

// Load JSON-based plugin configuration from plugin.json
func (p *plugin) pluginConfig(b Bot) (Plugin, error) {
	cf := fmt.Sprintf("%s/plugin.json", p.path)

	cb, err := ioutil.ReadFile(cf)
	if err != nil {
		return nil, fmt.Errorf("Cannot load plugin '%s': %s", p.Name, err)
	}

	pp := &plugin{}
	err = json.Unmarshal(cb, pp)
	if err != nil {
		return nil, fmt.Errorf(
			"Cannot load plugin '%s' configuration: %s",
			p.Name,
			err,
		)
	}

	if p.Name != pp.Name {
		log.Warningf(`
			Sanity check! %s has a different name than it was configured
			for ('%s')
		`, p.Name, pp.Name)
	}

	vs := b.VMs()
	if vm, ok := vs[pp.VM]; !ok {
		return nil, fmt.Errorf("The VM '%s' is not available", p.VM)
	} else {
		p.vm = vm
	}

	return pp, nil
}

// Compile our plugin, using the VM given
func (p *plugin) compile() error {
	m, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", p.path, p.Main))
	if err != nil {
		return err
	}
	_, err = p.vm.LoadScript(p.Name, string(m))
	if err != nil {
		return err
	}

	return nil
}
