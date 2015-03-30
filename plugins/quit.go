package main

import "regexp"

const DEFAULT_QUIT_MSG = "Leaving..."

var quit = regexp.MustCompile("^@quit")

func init() {
	NewPlugin("quit",
		TriggerCommands{
			quit: Command{
				Handler: quitHandler,
			},
		},
	)
}

func quitHandler(q *App, msgCh chan *Message) {
	msg := <-msgCh

	for _, nick := range q.Cfg.Admins {
		if msg.Nick == nick {
			q.Con.SendRawf("QUIT %s", DEFAULT_QUIT_MSG)
			q.Con.Disconnect()
			q.Con.Quit()
		}
	}
}
