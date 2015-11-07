package irc

// IRC Events
//
// The events system will coordinate events between reading from the server,
// and any actions that should be handled for those events, based on a Filter.

import (
	log "github.com/Sirupsen/logrus"
	"github.com/enmand/quarid-go/pkg/adapter"
)

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
func (i *Client) Handle(fs []adapter.Filter, hf adapter.HandlerFunc) {
	h := &adapter.Handler{
		Filters: fs,
		Handler: hf,
	}

	i.handlers = append(i.handlers, h)
}

func (i *Client) handleEvent(ev *adapter.Event) {
	log.Infof("Handling event: %#v", ev)
	for _, h := range i.handlers {
		for _, f := range h.Filters {
			log.Debugf("\tChecking filter: %#v", f)
			if f.Match(ev) {
				log.Debug("\t\tFilter matched")
				h.Handler(ev, i)
			}
		}
	}
}
