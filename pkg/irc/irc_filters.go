package irc

import (
	"github.com/enmand/quarid-go/pkg/adapter"
)

// CommandFilter filters events based on the IRC command of the event
type CommandFilter struct {
	Command string
}

// Match this filter, against incoming events.
func (cf CommandFilter) Match(ev *adapter.Event) bool {
	if cf.Command == "*" || cf.Command == "" {
		// Match all events
		return true
	}
	return cf.Command == ev.Command
}
