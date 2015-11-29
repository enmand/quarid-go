package main

import (
	"os"

	"github.com/enmand/quarid-go/pkg/bot"
	"github.com/enmand/quarid-go/pkg/config"
	"github.com/enmand/quarid-go/pkg/logger"
)

func main() {
	c := config.Get()
	logger.Log.Info("Loading IRC bot...")

	q := bot.NewManager(c.GetString("irc.nick"))
	ircbot := q.IRC(
		c.GetString("irc.server"),
		c.GetString("irc.user"),
		c.GetBool("irc.tls.enable"),
		c.GetBool("irc.tls.verify"),
	)
	if err := ircbot.Start(); err != nil {
		logger.Log.Errorf("%s", err)
		os.Exit(-1)
	}
	defer func() {
		q.Stop()
	}()
}
