// Package adapter adapts streams of information from third-party sources
// into into generic events that can be acted on by a Crate or a single Plugin
package adapter

import (
	"time"
)

// Event represents a single events read from the IRC server
type Event struct {
	// The event prefix (optional in spec)
	Prefix string

	// The command that the client (or server) is sending/sent
	Command string

	// The parameters to the command the client (or server) is sending/sent
	Parameters []string

	// The time the Event was received
	Timestamp time.Time
}

// An Adapter is an interface for adapting events
type Adapter interface {
	Start() error
	Stop() error

	EventsHandler
	Responder
}

// EventsHandler preforms actions for incoming events, based on a filter, or
// set of filters
type EventsHandler interface {
	// Handle a portion of events, based on a filter
	Handle(f []Filter, h HandlerFunc)
	HandleFilter(f Filter, h HandlerFunc)
}

// Responder writes an event to the server
type Responder interface {
	// Write to the server
	Write(ev *Event) error
}

// A Filter is ways to match Events from the server to some matching pattern,
// so that specific events can trigger specific actions (or not trigger actions)
type Filter interface {
	Match(ev *Event) bool
}

// Handler filters events to be acted on by a HandlerFunc
type Handler struct {
	Filters []Filter
	Handler HandlerFunc
}

// HandlerFunc preforms some action, based on the Event given, and respsonds
// using the Responder c
type HandlerFunc func(ev *Event, c Responder)
