package main

import (
	"fmt"

	"github.com/enmand/quarid-go/services"

	"github.com/alecthomas/kingpin"
	"github.com/thoj/go-ircevent"
)

const VERSION = "0.1"

var (
	configFile = kingpin.Flag(
		"config",
		"Configuration file",
	).Required().String()
)

type Message struct {
	*irc.Event

	Room string
}

func main() {
	kingpin.Version(VERSION)
	kingpin.Parse()

	var err error
	config, err := services.NewConfig(*configFile)
	if err != nil {
		fmt.Printf("Could not load configuration: %s\n", err)
		return
	}

	bot := services.NewBot(config)

	bot.Connect()
	defer bot.Disconnect()

	bot.Connection.Loop()
}
