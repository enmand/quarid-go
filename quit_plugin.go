package main

import (
	"regexp"

	"github.com/thoj/go-ircevent"
)

const DEFAULT_QUIT_MSG = "Leaving..."

func init() {
	const PLUGIN_NAME = "quit"

	trigger := regexp.MustCompile("^@quit")

	commands := TriggerCommands{
		trigger: Command{
			Handler: quitHandler,
		},
	}

	plugin := Plugin{
		Commands: commands,
	}

	RegisterPlugin(PLUGIN_NAME, plugin)
}

func quitHandler(con *irc.Connection, config *Config, msgCh chan *Message) {
	msg := <-msgCh

	for _, nick := range config.Admins {
		if msg.Nick == nick {
			con.SendRawf("QUIT %s", DEFAULT_QUIT_MSG)
			con.Disconnect()
		}
	}
}
