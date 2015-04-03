package services

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"

	"github.com/enmand/quarid-go/vm"

	log "github.com/Sirupsen/logrus"
	"github.com/thoj/go-ircevent"
)

type quarid struct {
	// Connection to the IRC server
	Connection *irc.Connection
	// Configuration from the user
	Config *Config
	// The Plugins we have loaded
	plugins []Plugin
	// The VM for our Plugins
	vms map[int]vm.VM
	// The bot's logging service
	Log *log.Logger
}

func (q *quarid) Init() error {
	q.Connection = irc.IRC(q.Config.Nick, q.Config.User)

	q.Connection.UseTLS = q.Config.TLS.Enable
	q.Connection.TLSConfig = &tls.Config{InsecureSkipVerify: !q.Config.TLS.Verify}

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

	q.vms = map[int]vm.VM{
		vm.JS: vm.New(vm.JS),
	}
	q.vms[vm.JS].Initialize()

	q.Log = q.Config.Logger

	return nil
}

func (q *quarid) LoadPlugins(dirs []string) ([]Plugin, []error) {
	var ps []Plugin
	var errs []error

	for _, d := range dirs {
		fis, err := ioutil.ReadDir(d)
		if err != nil {
			errs = append(errs, err)
		}

		for _, fi := range fis {
			if fi.IsDir() {
				p := NewPlugin(
					fi.Name(),
					fmt.Sprintf("%s/%s", d, fi.Name()),
				)
				if err := p.Load(); err != nil {
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
	return q.Connection.Connect(q.Config.Server)
}

func (q *quarid) Disconnect() {
	q.Connection.Disconnect()
}

func (q *quarid) AddCallback(event string, f func(e *irc.Event)) string {
	return q.Connection.AddCallback(event, f)
}

func (q *quarid) ClearCallback(event string) bool {
	return q.Connection.ClearCallback(event)
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
}
