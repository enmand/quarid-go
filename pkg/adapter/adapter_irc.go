package adapter

import (
	irc "github.com/enmand/go-ircclient"
)

// IRCResponder responds to IRC messages
type IRC struct {
	c *irc.Client

	server string
}

// NewIRCResponder creates a new responder to IRC requests
func NewIRC(server, nick, user string, tls, verifyCert bool) *IRC {
	c := irc.NewClient(nick, user, tls, verifyCert)

	return &IRC{c, server}
}

func (r *IRC) Start() error {
	go r.c.Loop()

	err := r.c.Connect(r.server)
	if err != nil {
		return err
	}

	return r.c.Read()
}

func (r *IRC) Stop() error {
	return r.c.Disconnect()
}

func (r *IRC) Handle(fs []Filter, h HandlerFunc) {
	r.handle(fs, h)
}

func (r *IRC) HandleFilter(f Filter, h HandlerFunc) {
	r.handle([]Filter{f}, h)
}

func (r *IRC) handle(fs []Filter, h HandlerFunc) {
	handleFs := []irc.Filter{}

	for _, f := range fs {
		filter := f.(IRCFilter).Filter
		handleFs = append(handleFs, filter)
	}

	r.c.Handle(handleFs, func(event *irc.Event, c irc.IRC) {
		h(FromIRCEvent(event), r)
	})

}

func (r *IRC) Write(ev *Event) error {
	ircEv := ToIRCEvent(ev)
	return r.c.Write(ircEv)
}

// IRCFilter matches an Filter for IRC events
type IRCFilter struct {
	Filter irc.Filter
}

// Match will match *Events for the given IRC filter
func (f IRCFilter) Match(ev *Event) bool {
	ircEv := ToIRCEvent(ev)

	return f.Filter.Match(ircEv)
}

// ToIRCEvent returns an *irc.Event for the *Event
func ToIRCEvent(ev *Event) *irc.Event {
	return &irc.Event{
		Prefix:     ev.Prefix,
		Command:    ev.Command,
		Parameters: ev.Parameters,
		Timestamp:  ev.Timestamp,
	}
}

// FromIRCEvent returns an *Event from an *irc.Event
func FromIRCEvent(ircEv *irc.Event) *Event {
	return &Event{
		Prefix:     ircEv.Prefix,
		Command:    ircEv.Command,
		Parameters: ircEv.Parameters,
		Timestamp:  ircEv.Timestamp,
	}
}
