// Events
//
// The events system will coordinate events between reading from the server,
// and any actions that should be handled for those events, based on a Filter.
package irc

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

type EventsHandler interface {
	// Handle a portion of events, based on a filter
	Handle(f []Filter, h HandlerFunc)
}

// Responder holds a connection, and can be used to write to the server
// Events are are read from the IRC server
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

// HandlerFuncs are given the event to handle, and a client to interact with the
// server with
type HandlerFunc func(ev *Event, c Responder)

// Handler filters events to be acted on by a HandlerFunc
type Handler struct {
	fs []Filter
	h  HandlerFunc
}

func (i *Client) Loop() {
	go i.read()

	for m := range i.events {
		go i.handleEvent(m)
	}

	<-i.dead
}

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
