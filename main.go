package main

import (
	"crypto/tls"
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/thoj/go-ircevent"
)

var (
	debug      = kingpin.Flag("debug", "Debug mode").Bool()
	configFile = kingpin.Flag("config", "Configuration file").Required().String()
)

const VERSION = "0.1"

type Message struct {
	*irc.Event

	Room string
}

func main() {
	kingpin.Version(VERSION)
	kingpin.Parse()
	config, err := loadConfig(*configFile)
	if err != nil {
		fmt.Errorf("Could not load configuration: %s\n", err)
		return
	}
	fmt.Printf("Loading plugins...")
	InitializePlugin(config)

	con := irc.IRC(config.Nick, config.User)
	con.Debug = *debug

	con.UseTLS = config.TLS.Enable
	con.TLSConfig = &tls.Config{InsecureSkipVerify: !config.TLS.Verify}
	if err := con.Connect(config.Server); err != nil {
		fmt.Printf("%#v", err)
	}
	defer con.Disconnect()

	con.AddCallback("CTCP_VERSION", func(e *irc.Event) {
		con.Noticef(e.Nick, "\x01VERSION Quarid %s", VERSION)
	})

	con.AddCallback("001", func(e *irc.Event) {
		for _, room := range config.Channels {
			con.Join(room)
		}
	})

	con.AddCallback("PRIVMSG", func(event *irc.Event) {
		go func(e *irc.Event) {
			msg := &Message{e, e.Arguments[0]}

			msgCh := make(chan *Message)
			cmdCh := make(chan Command)
			go findCommand(msgCh, cmdCh)

			msgCh <- msg
			_cmd := <-cmdCh
			go _cmd.Handler(con, config, msgCh)

		}(event)
	})

	con.Loop()
}
