package services

import (
	"fmt"

	"github.com/enmand/quarid-go/vm"

	log "github.com/Sirupsen/logrus"
	"github.com/thoj/go-ircevent"
)

type Bot interface {
	// Run initialization
	Init() error

	// Load all of the plugins in each `dir`
	LoadPlugins(dir []string) ([]Plugin, []error)

	// Connect to the configured server
	Connect() error

	// Disconnect from the server
	Disconnect()

	// Add a callback based on an IRC event
	AddCallback(event string, f func(e *irc.Event)) string

	// Clear all callbacks
	ClearCallback(event string) bool

	// Return a map of initialzed virtual machines for this bot
	VMs() map[int]vm.VM
}

func NewBot(cfg *Config) *quarid {
	bot := &quarid{
		Config: cfg,
	}
	err := bot.Init()
	if err != nil {
		panic(fmt.Sprintf("Could not initialize bot: %s", err))
	}

	bot.ClearCallback("CTCP_VERSION")
	bot.AddCallback("CTCP_VERSION", func(e *irc.Event) {
		bot.Connection.Noticef(e.Nick, "\x01VERSION Quarid %s", VERSION)
	})

	bot.AddCallback("001", func(e *irc.Event) {
		for _, room := range bot.Config.Channels {
			bot.Connection.Join(room)
		}
	})

	bot.AddCallback("*", func(e *irc.Event) {
		log.WithFields(log.Fields{
			"code":      e.Code,
			"message":   e.Message(),
			"arguments": e.Arguments,
			"nick":      e.Nick,
			"user":      e.User,
			"host":      e.Host,
		}).Debug("Got event")
	})

	return bot
}
