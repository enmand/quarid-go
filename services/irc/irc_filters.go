package irc

// Filters
//
// Filters define ways to match Events from the server to some matching pattern,
// so that specific events can trigger specific actions (or not trigger actions)

// Filter is an interface that matches Events
type Filter interface {
	Match(ev *Event) bool
}

// CommandFilter filters events based on the IRC command of the event
type CommandFilter struct {
	Command string
}

// Match this filter, against incoming events.
func (cf CommandFilter) Match(ev *Event) bool {
	if cf.Command == "*" || cf.Command == "" {
		// Match all events
		return true
	}
	return cf.Command == ev.Command
}
