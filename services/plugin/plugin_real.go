package plugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	qvm "github.com/enmand/quarid-go/vm"

	log "github.com/Sirupsen/logrus"
)

type plugin struct {
	path string `json:"-"`
	vm   qvm.VM `json:"-"`

	Name          string      `json:"name"`
	VM            string      `json:"vm"`
	Main          string      `json:"main"`
	Configuration interface{} `json:"configuration"`
}

func (p *plugin) Run() error {
	_, err := p.vm.Run(p.path)
	return err
}

// Compile our plugin, using the VM given
func (p *plugin) Compile() error {
	m, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", p.path, p.Main))
	if err != nil {
		return err
	}

	_, err = p.vm.LoadScript(p.path, string(m))
	if err != nil {
		return err
	}

	return nil
}

func (p *plugin) Load(vms map[string]qvm.VM) error {
	pp, err := p.pluginConfig(vms)
	if err != nil {
		return err
	}

	p = pp.(*plugin) // Set our plugin configuration our loaded config
	log.Infof("Loading plugin: %s (in VM: %s)", p.Name, p.VM)

	if err := p.Compile(); err != nil {
		return err
	}

	return nil
}

// Load JSON-based plugin configuration from plugin.json
func (p *plugin) pluginConfig(vms map[string]qvm.VM) (Plugin, error) {
	cf := fmt.Sprintf("%s/plugin.json", p.path)

	cb, err := ioutil.ReadFile(cf)
	if err != nil {
		return nil, fmt.Errorf("Cannot load plugin '%s': %s", p.Name, err)
	}

	pp := &plugin{}
	pp.path = p.path

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

	if v, ok := vms[pp.VM]; !ok {
		return nil, fmt.Errorf("The VM '%s' is not available", p.VM)
	} else {
		pp.vm = v
	}

	return pp, nil
}
