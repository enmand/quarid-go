package main

import (
	"regexp"
	"strings"
	"time"
)

const seenDb = "seen.plugin"

var (
	updateSeenAt = regexp.MustCompile(".*")
	seen         = regexp.MustCompile("^@seen (.*)")
)

func init() {
	NewPlugin("seen",
		TriggerCommands{
			updateSeenAt: Command{
				Handler: updateSeenHandler,
			},
			seen: Command{
				Handler: seenTriggerHandler,
			},
		},
	)
}

func updateSeenHandler(q *App, msgCh chan *Message) {
	msg := <-msgCh

	db, err := GetDatabase(seenDb)
	if err != nil {
		q.Con.Log.Fatal(err)
	}

	room := strings.ToLower(msg.Room)
	nick := strings.ToLower(msg.Nick)

	db.Insert(map[interface{}]interface{}{
		"seen":   time.Now().Unix(),
		"saying": msg.Message(),
	}, room, nick)
}

func seenTriggerHandler(q *App, msgCh chan *Message) {
	msg := <-msgCh

	args := strings.Split(msg.Message(), " ")
	name := strings.ToLower(args[1])

	db, err := GetDatabase(seenDb)
	if err != nil {
		q.Con.Log.Fatal(err)
	}

	data, err := db.Find(msg.Room, name)
	if err != nil {
		q.Con.Log.Print(err)
		q.Con.Privmsgf(msg.Room, "I've never seen '%s'", name)

		return
	}

	seen := time.Unix(data["seen"].(int64), 0).
		In(q.Cfg.Timezone).
		Format(time.RFC822)

	q.Con.Privmsgf(msg.Room, "Last saw '%s' at '%s' saying '%s'",
		name,
		seen,
		strings.TrimSpace(data["saying"].(string)),
	)
}
