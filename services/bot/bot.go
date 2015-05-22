package bot

import (
	"fmt"

	"github.com/enmand/quarid-go/services/config"
	"github.com/enmand/quarid-go/services/plugin"
	"github.com/enmand/quarid-go/vm"
)

type Event interface{}

type Bot interface {
	// Load all of the plugins in each `dir`
	LoadPlugins(dir []string) ([]plugin.Plugin, []error)

	// Connect to the configured server
	Connect() error

	// Disconnect from the server
	Disconnect()

	// Add a callback based on an IRC event
	AddCallback(event string, f func(e *Event)) string

	// Clear all callbacks
	ClearCallback(event string) bool

	// Return a map of initialzed Plugins for this bot
	Plugins() []plugin.Plugin

	// A list of VMs that the bot has available
	VMs() map[string]vm.VM

	Read(ch chan Event)
}

func New(cfg *config.Config) *quarid {
	bot := &quarid{
		Config: cfg,
	}

	err := bot.initialize()
	if err != nil {
		panic(fmt.Sprintf("Could not initialize bot: %s", err))
	}

	return bot
}
