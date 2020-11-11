package main

import (
	"fmt"
	"log"

	"github.com/enmand/quarid-go/internal/adapter"
	"github.com/enmand/quarid-go/internal/bot"
	"github.com/enmand/quarid-go/internal/config"
)

func main() {
	c := config.Get()
	a, err := adapter.NewMatrix(c.GetString("matrix.homeserver"), c.GetString("matrix.user"), c.GetString("matrix.token"))
	if err != nil {
		log.Fatalln(err)
	}
	bot := bot.New([]adapter.Adapter{a})

	errCh := make(chan error)
	go start(bot, errCh)
	log.Fatalln(<-errCh)
}

func start(b *bot.Bot, errCh chan error) {
	b.Start(errCh)
	if err := <-errCh; err != nil {
		panic(fmt.Sprintf("Unable to start: %s", err))
	}
}
