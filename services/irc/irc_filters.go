// Filters
//
// Filters define ways to match Events from the server to some matching pattern,
// so that specific events can trigger specific actions (or not trigger actions)
package irc

// Filter is an interface that matches Events
type Filter interface {
	Match(ev *Event) bool
}

type CommandFilter struct {
	Command string
}

func (cf CommandFilter) Match(ev *Event) bool {
	return cf.Command == ev.Command
}
