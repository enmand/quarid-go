package main

import (
	"fmt"
	"net/http"
	"regexp"
)

var urlRegex = regexp.MustCompile("@(https?)://(-\\.)?([^\\s/?\\.#-]+\\.?)+(/[^\\s]*)?$@iS")

func init() {
	NewPlugin("urltitle",
		TriggerCommands{
			urlRegex: Command{
				Handler: urlTitleHandler,
			},
		},
	)
}

func urlTitleHandler(q *App, msgCh chan *Message) {
	msg := <-msgCh

	if matches := urlRegex.FindAll([]byte(msg.Message()), 0); len(matches) > 0 {
		for _, _match := range matches {
			fmt.Printf("%#v", _match)
			go func(url []byte) {
				resp, err := http.Head(string(url))
				if err != nil {
					q.Con.Log.Printf(
						"Could not fetch URL '%s': %s'",
						string(url),
						err,
					)
				}

				q.Con.Privmsgf(
					msg.Room,
					"Looks like it's %d",
					resp.ContentLength,
				)

			}(_match)
		}
	}
}
