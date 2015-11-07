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

	q := bot.New(&c)

	if err := q.Connect(); err != nil {
		logger.Log.Errorf("%s", err)
		os.Exit(-1)
	}
	defer func() {
		q.Disconnect()
	}()
}
