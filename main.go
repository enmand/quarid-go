package main

import (
	"os"

	"github.com/enmand/quarid-go/services"
	"github.com/enmand/quarid-go/services/bot"
	"github.com/enmand/quarid-go/services/config"

	log "github.com/Sirupsen/logrus"
	"github.com/alecthomas/kingpin"
)

var (
	configFile = kingpin.Flag(
		"config",
		"Configuration file",
	).Required().String()
)

func main() {
	kingpin.Version(services.VERSION)
	kingpin.Parse()

	var err error
	config, err := config.New(*configFile)
	if err != nil {
		log.Errorf("Could not load configuration: %s\n", err)
		return
	}

	q := bot.New(config)

	if err := q.Connect(); err != nil {
		log.Errorf("%s", err)
		os.Exit(-1)
	}
	defer func() {
		q.Disconnect()
	}()
}
