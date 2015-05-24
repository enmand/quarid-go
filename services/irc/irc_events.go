// Events
//
// The events system will coordinate events between reading from the server,
// and any actions that should be handled for those events, based on a Filter.
package irc

import "time"

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

// Filter is an interface that matches Events
type Filter interface {
	Match(ev *Event) bool
}

type handleFunc func(ev *Event)
