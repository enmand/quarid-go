package adapter

import (
	irc "github.com/enmand/go-ircclient"
)

// IRCResponder responds to IRC messages
type IRCResponder struct {
	c *irc.Client
}

// NewIRCResponder creates a new responder to IRC requests
func NewIRCResponder(c *irc.Client) IRCResponder {
	return IRCResponder{c}
}

func (r IRCResponder) Write(ev *Event) error {
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
