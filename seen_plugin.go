package main

import (
	"regexp"
	"strings"
	"time"

	"github.com/thoj/go-ircevent"
)

const SEEN_DB = "seen.plugin"

func init() {
	const PLUGIN_NAME = "seen"

	updateTrigger := regexp.MustCompile(".*")
	seenTrigger := regexp.MustCompile("^@seen (.*)")

	commands := TriggerCommands{
		seenTrigger: Command{
			Handler: seenTriggerHandler,
		},
		updateTrigger: Command{
			Handler: updateSeenHandler,
		},
	}

	RegisterPlugin(PLUGIN_NAME, Plugin{
		Commands: commands,
	})
}

func updateSeenHandler(
	con *irc.Connection,
	config *Config,
	msgCh chan *Message,
) {
	msg := <-msgCh

	db, err := GetDatabase(SEEN_DB)
	if err != nil {
		con.Log.Fatal(err)
	}

	room := strings.ToLower(msg.Room)
	nick := strings.ToLower(msg.Nick)

	db.Data = map[string]interface{}{
		room: map[string]interface{}{
			nick: map[string]interface{}{
				"seen":   time.Now().Unix(),
				"saying": msg.Message(),
			},
		},
	}

	db.Save()
}

func seenTriggerHandler(
	con *irc.Connection,
	config *Config,
	msgCh chan *Message,
) {
	msg := <-msgCh

	args := strings.Split(msg.Message(), " ")
	name := strings.ToLower(args[1])

	db, err := GetDatabase(SEEN_DB)
	if err != nil {
		con.Log.Fatal(err)
	}

	if roomIdx, ok := db.Data[msg.Room]; ok {
		roomIdx := roomIdx.(map[interface{}]interface{})
		if nickIdx, ok := roomIdx[name]; ok {
			nickIdx := nickIdx.(map[interface{}]interface{})

			con.Privmsgf(
				msg.Room,
				"Last saw '%s' at '%s' saying '%s'",
				name,
				time.Unix(nickIdx["seen"].(int64), 0),
				nickIdx["saying"].(string),
			)
		} else {
			con.Privmsgf(msg.Room, "I've never seen '%s'", name)
		}
	}
}
