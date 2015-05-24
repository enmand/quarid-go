package bot

import (
	"github.com/enmand/quarid-go/services/config"
	"github.com/enmand/quarid-go/services/irc"
	"github.com/enmand/quarid-go/services/plugin"
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

	// The bot's logging service
	Log *log.Logger

	// The event channel for IRC to use for this bot
	IRCEventCh chan *irc.Event

	// Exit flag
	exit chan bool
}

func (q *quarid) initialize() error {
	q.IRC = irc.NewClient(q.Config).(*irc.Client)
	q.IRCEventCh = make(chan *irc.Event)
	q.exit = make(chan bool)

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
	/*
		q.ClearCallback("CTCP_VERSION")
		q.AddCallback("CTCP_VERSION", func(e *Event) {
			q.IRC.Noticef(e.Nick, "\x01VERSION Quarid %s", VERSION)
		})

		ready := make(chan bool)
		q.IRC.AddCallback("001", func(e *Event) {
			for _, room := range q.Config.Channels {
				q.IRC.Join(room)
			}

			ready <- true
		})

		go func(ch chan bool) {
			select {
			case <-ch:
				q.IRC.AddCallback("*", func(e *Event) {
					log.WithFields(log.Fields{
						"code":      e.Code,
						"message":   e.Message(),
						"arguments": e.Arguments,
						"nick":      e.Nick,
						"user":      e.User,
						"host":      e.Host,
					}).Debug("Got event")
				})
			}
		}(ready)*/

	return nil
}

func (q *quarid) LoadPlugins(dirs []string) ([]plugin.Plugin, []error) {
	var ps []plugin.Plugin
	var errs []error
	/*
		for _, d := range dirs {
			fis, err := ioutil.ReadDir(d)
			if err != nil {
				errs = append(errs, err)
			}

			for _, fi := range fis {
				if fi.IsDir() {
					p := plugins.New(
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
	*/

	return ps, errs
}

func (q *quarid) Connect() error {
	err := q.IRC.Connect(q.Config.Server)

	return err
}

func (q *quarid) Disconnect() {
	q.IRC.Disconnect()
}

/*func (q *quarid) AddCallback(event string, f func(e *Event)) string {
	return q.IRC.AddCallback(event, f)
}

func (q *quarid) ClearCallback(event string) bool {
	return q.IRC.ClearCallback(event)
}

func (q *quarid) Plugins() []Plugin {
	return q.plugins
}

func (q *quarid) VMs() map[string]vm.VM {
	return q.vms
}

func (q *quarid) Debugf(s string, f ...interface{}) {
	log.Debugf(s, f)
}

func (q *quarid) Infof(s string, f ...interface{}) {
	log.Infof(s, f)
}

func (q *quarid) Warningf(s string, f ...interface{}) {
	log.Warningf(s, f)
}

func (q *quarid) Errorf(s string, f ...interface{}) {
	log.Errorf(s, f)
}*/
