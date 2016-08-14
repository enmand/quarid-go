package main

import (
	"github.com/Sirupsen/logrus"
	irc "github.com/enmand/go-ircclient"
	"github.com/enmand/quarid-go/pkg/adapter"
	"github.com/enmand/quarid-go/pkg/bot"
	"github.com/enmand/quarid-go/pkg/config"
)

func main() {
	c := config.Get()
	logrus.Info("Loading IRC bot...")

	q := bot.NewManager(c.GetString("irc.nick"))
	ircbot := q.IRC(
		c.GetString("irc.server"),
		c.GetString("irc.user"),
		c.GetStringSlice("irc.channels"),
		c.GetBool("irc.tls.enable"),
		c.GetBool("irc.tls.verify"),
	)

	ircbot.Handle(
		[]adapter.Filter{
			adapter.IRCFilter{
				Filter: irc.CommandFilter{
					Command: "*",
				},
			},
		},
		logIRC,
	)

	errCh := make(chan error)
	go start(ircbot, errCh)

	logrus.Fatalf("Error running bot: %s", <-errCh)
}

func start(a adapter.Adapter, errCh chan error) {
	err := a.Start()
	defer func() {
		err := a.Stop()
		if err != nil {
			errCh <- err
		}
	}()

	errCh <- err
}

func logIRC(ev *adapter.Event, r adapter.Responder) {
	logrus.Infof(
		"(%s) %s: %s (%q))",
		ev.Timestamp,
		ev.Prefix,
		ev.Command,
		ev.Parameters,
	)
}
