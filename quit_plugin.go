package main

import (
	"fmt"
	"regexp"

	"github.com/thoj/go-ircevent"
)

const DEFAULT_QUIT_MSG = "Leaving..."

func init() {
	const PLUGIN_NAME = "quit"

	trigger, err := regexp.Compile("^@quit")
	if err != nil {
		fmt.Errorf("Could not load plugin %s", PLUGIN_NAME)
	}

	commands := map[*regexp.Regexp]Command{
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
