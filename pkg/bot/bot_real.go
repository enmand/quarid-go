package bot

import (
	"fmt"
	"io/ioutil"

	"github.com/enmand/quarid-go/pkg/config"
	"github.com/enmand/quarid-go/pkg/irc"
	"github.com/enmand/quarid-go/pkg/plugin"
	"github.com/enmand/quarid-go/vm"

	log "github.com/Sirupsen/logrus"
)

type quarid struct {
	// Connection to the IRC server
	IRC *irc.Client

	// Configuration from the user
	Config *config.Config

	// The Plugins we have loaded
	plugins []plugin.Plugin

	// The VM for our Plugins
	vms map[string]vm.VM
}

func (q *quarid) initialize() error {
	q.IRC = irc.NewClient(q.Config).(*irc.Client)

	// Initialize our VMs
	q.vms = map[string]vm.VM{
		vm.JS: vm.New(vm.JS),
	}

	var errs []error
	q.plugins, errs = q.LoadPlugins(q.Config.PluginsDirs)
	if errs != nil {
		log.Warningf(
			"Some plugins failed to load. The following are loaded: %q",
			q.plugins,
		)
		log.Warningf("But the follow errors occurred:")
		for _, e := range errs {
			log.Warning(e)
		}
	}

	return nil
}

func (q *quarid) LoadPlugins(dirs []string) ([]plugin.Plugin, []error) {
	var ps []plugin.Plugin
	var errs []error

	for _, d := range dirs {
		fis, err := ioutil.ReadDir(d)
		if err != nil {
			errs = append(errs, err)
		}

		for _, fi := range fis {
			if fi.IsDir() {
				p := plugin.NewPlugin(
					fi.Name(),
					fmt.Sprintf("%s/%s", d, fi.Name()),
				)
				if err := p.Load(q.VMs()); err != nil {
					errs = append(errs, err)
				} else {
					ps = append(ps, p)
				}

			}
		}
	}

	return ps, errs
}

func (q *quarid) Connect() error {
	err := q.IRC.Connect(q.Config.Server)

	q.IRC.Loop()

	return err
}

func (q *quarid) Disconnect() {
	q.IRC.Disconnect()
}

func (q *quarid) Plugins() []plugin.Plugin {
	return q.plugins
}

func (q *quarid) VMs() map[string]vm.VM {
	return q.vms
}
