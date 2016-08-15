package bot

import (
	irc "github.com/enmand/go-ircclient"
	"github.com/enmand/quarid-go/pkg/adapter"
)

type ircbot struct {
	// Connection to the IRC server
	IRC *irc.Client

	// IRC server to connect to
	server string
}

// NewIRC returns an adapter.Adapter, with an open connection to an IRC server
func NewIRC(
	server, nick, ident string,
	autojoins []string,
	tls, tlsVerify bool,
) adapter.Adapter {
	c := irc.NewClient(nick, ident, tls, tlsVerify)

	b := &ircbot{
		IRC:    c,
		server: server,
	}

	b.Handle([]adapter.Filter{
		adapter.IRCFilter{
			Filter: irc.CommandFilter(irc.IRC_RPL_WELCOME),
		},
	}, autoJoinChans(autojoins))

	return b
}

// Connect this bot to the IRC server
func (q *ircbot) Start() error {
	go q.IRC.Loop() // loop handles IRC events

	err := q.IRC.Connect(q.server)
	if err != nil {
		return err
	}

	return q.IRC.Read()
}

func (q *ircbot) Stop() error {
	return q.IRC.Disconnect()
}

func (q *ircbot) Handle(fs []adapter.Filter, hf adapter.HandlerFunc) {
	ircFs := []irc.Filter{}

	for _, f := range fs {
		ircFilter := f.(adapter.IRCFilter).Filter
		ircFs = append(ircFs, ircFilter)
	}

	q.IRC.Handle(ircFs, func(ircEv *irc.Event, c irc.IRC) {
		resp := adapter.NewIRCResponder(c.(*irc.Client))
		hf(adapter.FromIRCEvent(ircEv), resp)
	})

}

func (q *ircbot) Write(ev *adapter.Event) error {
	return q.IRC.Write(adapter.ToIRCEvent(ev))
}

func autoJoinChans(
	chans []string,
) func(*adapter.Event, adapter.Responder) {
	return func(ev *adapter.Event, r adapter.Responder) {
		_ = r.Write(&adapter.Event{
			Command:    irc.IRC_JOIN,
			Parameters: chans,
		})
	}
}
