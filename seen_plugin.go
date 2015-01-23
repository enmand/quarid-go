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

	db.Insert(map[interface{}]interface{}{
		"seen":   time.Now().Unix(),
		"saying": msg.Message(),
	}, room, nick)
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

	data, err := db.Find(msg.Room, name)
	if err != nil {
		con.Log.Print(err)
		con.Privmsgf(msg.Room, "I've never seen '%s'", name)

		return
	}

	seen := time.Unix(data["seen"].(int64), 0).
		Local().
		Format(time.RFC822)

	con.Privmsgf(msg.Room, "Last saw '%s' at '%s' saying '%s'",
		name,
		seen,
		data["saying"].(string),
	)
}
