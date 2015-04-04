package main

import (
	"github.com/enmand/quarid-go/services"

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
	config, err := services.NewConfig(*configFile)
	if err != nil {
		log.Errorf("Could not load configuration: %s\n", err)
		return
	}

	bot := services.NewBot(config)

	bot.Connect()
	defer bot.Disconnect()

	ch, errCh := bot.Connection.Loop()
	go func(e chan error) {
		log.Errorf("Got an error: %s", <-e)
	}(errCh)

	<-ch
}
