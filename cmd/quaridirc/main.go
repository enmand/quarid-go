package main

import (
	"fmt"

	"regexp"

	"time"

	"github.com/Sirupsen/logrus"
	irc "github.com/enmand/go-ircclient"
	"github.com/enmand/quarid-go/pkg/adapter"
	"github.com/enmand/quarid-go/pkg/bot"
	"github.com/enmand/quarid-go/pkg/config"
	"github.com/enmand/quarid-go/pkg/plugin"
)

func main() {
	c := config.Get()
	logrus.Info("Loading IRC bot...")

	a := adapter.NewIRC(c.GetString("irc.server"),
		c.GetString("irc.nick"),
		c.GetString("irc.user"),
		c.GetBool("irc.tls.enable"),
		c.GetBool("irc.tls.verify"),
	)

	a.Handle(
		[]adapter.Filter{
			adapter.IRCFilter{
				Filter: irc.CommandFilter("*"),
			},
		},
		logIRC,
	)
	a.Handle(
		[]adapter.Filter{
			adapter.IRCFilter{
				Filter: irc.CommandFilter(irc.IRC_RPL_WELCOME),
			},
		},
		func(event *adapter.Event, r adapter.Responder) {
			_ = r.Write(&adapter.Event{
				Command:    irc.IRC_JOIN,
				Parameters: c.GetStringSlice("irc.channels"),
			})
		},
	)

	a.Handle(
		[]adapter.Filter{
			adapter.IRCFilter{
				Filter: &irc.RegExpFilter{
					Param:      irc.RegExpFilterParameters,
					Expression: *regexp.MustCompile("sup\\?"),
				},
			},
		},
		func(event *adapter.Event, r adapter.Responder) {
			event.Timestamp = time.Time{}

			fmt.Printf("%# v", event)

			r.Write(&adapter.Event{
				Command:    irc.IRC_PRIVMSG,
				Parameters: []string{event.Parameters[0], "not much, mate"},
			})
		})

	ircbot := bot.New([]adapter.Adapter{a})

	plugins, err := plugin.LoadPlugins(c.GetStringSlice("plugin_dirs"))
	if err != nil {
		logrus.Fatalf("error loading plugins: %s", err)
	}

	fmt.Printf("%# v", plugins)

	errCh := make(chan error)
	go start(ircbot, errCh)

	logrus.Fatalf("Error running bot: %s", <-errCh)
}

func start(b *bot.Bot, errCh chan error) {
	b.Start(errCh)
	if err := <-errCh; err != nil {
		panic(fmt.Sprintf("Unable to start: %s", err))
	}
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
