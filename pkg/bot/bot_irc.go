package bot

import (
	"github.com/enmand/quarid-go/pkg/adapter"
	"github.com/enmand/quarid-go/pkg/irc"
	"github.com/enmand/quarid-go/pkg/logger"
)

type ircbot struct {
	// Connection to the IRC server
	IRC *irc.Client

	// IRC server to connect to
	server string
}

// NewIRC returns an adpater.Adpater, with an open connection to an IRC server
func NewIRC(server, nick, ident string, tls, tlsVerify bool) adapter.Adapter {
	c := irc.NewClient(nick, ident, tls, tlsVerify)

	b := &ircbot{
		IRC:    c,
		server: server,
	}

	return b
}

// Connect this bot to the IRC server
func (q *ircbot) Start() error {
	go q.IRC.Loop()

	err := q.IRC.Connect(q.server)
	if err != nil {
		return err
	}

	rCh := make(chan error)
	go func(ch chan error) {
		ch <- q.IRC.Read()
	}(rCh)

	if readErr := <-rCh; readErr != nil {
		logger.Log.Errorf(err.Error())
		return err
	}

	return err
}

func (q *ircbot) Stop() error {
	q.IRC.Disconnect()

	return nil
}

func (q *ircbot) Handle(fs []adapter.Filter, hf adapter.HandlerFunc) {

}

func (q *ircbot) Write(ev *adapter.Event) error {
	return q.IRC.Write(ev)
}
