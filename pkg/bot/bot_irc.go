package bot

import (
	"fmt"
	"io/ioutil"

	"github.com/enmand/quarid-go/pkg/adapter"
	"github.com/enmand/quarid-go/pkg/config"
	"github.com/enmand/quarid-go/pkg/irc"
	"github.com/enmand/quarid-go/pkg/plugin"
	"github.com/enmand/quarid-go/vm"
	"github.com/enmand/quarid-go/vm/js"
	"github.com/renstrom/shortuuid"

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
	q.IRC = irc.NewClient(
		q.Config.GetString("irc.nick"),
		q.Config.GetString("irc.user"),
		q.Config.GetBool("irc.tls.verify"),
		q.Config.GetBool("irc.tls.enable"),
	)

	// Initialize our VMs
	q.vms = map[string]vm.VM{
		vm.JS: js.NewVM(),
	}

	var errs []error
	q.plugins, errs = q.LoadPlugins(q.Config.GetStringSlice("plugins_dirs"))
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
	err := q.IRC.Connect(q.Config.GetString("irc.server"))

	q.IRC.Handle(
		[]adapter.Filter{irc.CommandFilter{Command: irc.IRC_ERR_NICKNAMEINUSE}},
		q.fixNick,
	)

	q.IRC.Handle(
		[]adapter.Filter{irc.CommandFilter{Command: irc.IRC_RPL_MOTD}},
		q.joinChan,
	)

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

func (q *quarid) fixNick(
	ev *adapter.Event,
	c adapter.Responder,
) {
	nick := q.IRC.Nick
	uniq := shortuuid.UUID()

	newNick := fmt.Sprintf("%s_%s", nick, uniq)

	fixnickCmd := &adapter.Event{
		Command:    irc.IRC_NICK,
		Parameters: []string{newNick},
	}

	if err := c.Write(fixnickCmd); err == nil {
		q.IRC.Nick = newNick
	}
}

func (q *quarid) joinChan(
	ev *adapter.Event,
	c adapter.Responder,
) {
	chans := q.Config.GetStringSlice("irc.channels")

	joinCmd := &adapter.Event{
		Command:    irc.IRC_JOIN,
		Parameters: chans,
	}
	c.Write(joinCmd)
}