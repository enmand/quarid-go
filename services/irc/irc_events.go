package irc

// IRC Events
//
// The events system will coordinate events between reading from the server,
// and any actions that should be handled for those events, based on a Filter.

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

// EventsHandler preforms actions for incoming events, based on a filter, or
// set of filters
type EventsHandler interface {
	// Handle a portion of events, based on a filter
	Handle(f []Filter, h HandlerFunc)
}

// Event represents a single events read from the IRC server
type Event struct {
	// The event prefix (optional in spec)
	Prefix string

	// The command that the client (or server) is sending/sent
	Command string

	// The parameters to the command the client (or server) is sending/sent
	Parameters []string

	// The time the Event was recieved
	Timestamp time.Time
}

// HandlerFunc preforms some action, based on the Event given, and respsonds
// using the Responder c
type HandlerFunc func(ev *Event, c Responder)

// Handler filters events to be acted on by a HandlerFunc
type Handler struct {
	fs []Filter
	h  HandlerFunc
}

// Loop reads the data from the server, and handles events that happen
func (i *Client) Loop() {
	go i.read()

	for m := range i.events {
		go i.handleEvent(m)
	}

	<-i.dead
}

// Handle defines events that should be filtered to preform a handler function.
// Using "*" or "" for a filter, will cause all events to be passed to the
// HandlerFunc.
func (i *Client) Handle(fs []Filter, hf HandlerFunc) {
	h := &Handler{
		fs: fs,
		h:  hf,
	}

	i.handlers = append(i.handlers, h)
}

func (i *Client) handleEvent(ev *Event) {
	log.Infof("Handling event: %#v", ev)
	for _, h := range i.handlers {
		for _, f := range h.fs {
			log.Debugf("\tChecking filter: %#v", f)
			if f.Match(ev) {
				log.Debug("\t\tFilter matched")
				h.h(ev, i)
			}
		}
	}
}
