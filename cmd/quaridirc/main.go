package main

import (
	"github.com/Sirupsen/logrus"
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
		c.GetBool("irc.tls.enable"),
		c.GetBool("irc.tls.verify"),
	)

	if err := ircbot.Start(); err != nil {
		logger.Log.Errorf("%s", err)
	}
	defer func() {
		ircbot.Stop()
	}()
}
