package main

import (
	"math/rand"
	"regexp"
	"time"
)

var random = regexp.MustCompile("^@random")

func init() {
	NewPlugin("random",
		TriggerCommands{
			random: Command{
				Handler: randomHandler,
			},
		},
	)
}

func randomHandler(q *App, msgCh chan *Message) {
	msg := <-msgCh
	rand.Seed(time.Now().Unix())
	rand := rand.ExpFloat64()
	Quarid.Con.Privmsgf(msg.Room, "%f", rand)
}
