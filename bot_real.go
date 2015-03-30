package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"

	"github.com/enmand/quarid-go/services"
	"github.com/enmand/quarid-go/vm"

	log "github.com/Sirupsen/logrus"
	"github.com/thoj/go-ircevent"
)

type quarid struct {
	// Connection to the IRC server
	Connection *irc.Connection
	// Configuration from the user
	Config *services.Config
	// The VM for our plugins
	vms map[int]vm.VM
	// The bot's logging service
	Log *log.Logger
}

func (q *quarid) Init() error {
	q.Connection = irc.IRC(q.Config.Nick, q.Config.User)

	q.Connection.UseTLS = q.Config.TLS.Enable
	q.Connection.TLSConfig = &tls.Config{InsecureSkipVerify: !q.Config.TLS.Verify}

	q.vms = map[int]vm.VM{
		vm.JS: vm.New(vm.JS),
	}
	q.vms[vm.JS].Initialize()
	script, err := ioutil.ReadFile("./plugins/main.js")
	if err != nil {
		return err
	}
	_, err = q.vms[vm.JS].LoadScript("seen", string(script))
	if err != nil {
		return fmt.Errorf("Could not load script: %s", err)
	}
	_, err = q.vms[vm.JS].Run("seen")
	if err != nil {
		return err
	}

	q.Log = q.Config.Logger

	return nil
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
