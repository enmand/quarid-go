package irc

// IRC Events
//
// The events system will coordinate events between reading from the server,
// and any actions that should be handled for those events, based on a Filter.

import (
	"bufio"
	"fmt"
	"io"
	"net/textproto"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/enmand/quarid-go/pkg/adapter"
	"github.com/enmand/quarid-go/pkg/logger"
)

const (
	InvalidLineSize = "Could not parse line: %d parameter given, %d expected"
)

// Read reads the data from the server, and handles events that happen
func (i *Client) Read(n int) {
	go i.read(n)

	for m := range i.events {
		go i.handleEvent(m)
	}

	if n == 0 {
		<-i.dead
	}
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

// handleEvent will forward events to the proper handlers
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

// read n lines from the server. if n is 0, continue reading until we can't
func (i *Client) read(n int) {
	r := bufio.NewReader(i.conn)
	tp := textproto.NewReader(r)

	for current := 0; n == 0 || current < n; current++ {
		l, err := tp.ReadLine()
		switch err {
		case io.EOF:
			logger.Log.Debugf("Read EOF after %d lines", current-1)
			break
		case nil:
			ev, err := parseLine(l)
			if err != nil {
				logger.Log.Errorf(err.Error())
				continue
			}
			i.events <- ev
		default:
			logger.Log.Errorf("Error reading from server... passing: %s", err)
			break
		}
	}
}

// parseLine read from the IRC server
func parseLine(l string) (*adapter.Event, error) {
	ev := &adapter.Event{}
	ws := strings.Split(l, " ") // split args on " "
	var paramIndex int          // the argument index where the parameters are

	// Make sure we have at least two params (PREFIX and COMMAND)
	if len(ws) < 1 {
		return nil, fmt.Errorf(InvalidLineSize, len(ws), 2)
	}

	// Check if our "prefix" has ":"
	if ws[0][0] == ':' {
		// Server sent a prefix
		ev.Prefix = ws[0][1:]
		ev.Command = ws[1]

		paramIndex = 2
	} else {
		// Server did not send a prefix
		ev.Prefix = ""
		ev.Command = ws[0]

		paramIndex = 2
	}

	ev.Parameters = readParams(ws, paramIndex)
	ev.Timestamp = time.Now()

	return ev, nil
}

// readParams p with parameteres starting at index
func readParams(ws []string, index int) []string {
	var params []string
	wsSize := len(ws)

	for i, p := range ws[index:wsSize] {
		if len(p) > 0 && p[0] == ':' {
			params = append(params, strings.Join(ws[index+i:wsSize], " "))
			break
		} else {
			params = append(params, p)
		}
	}

	return params
}
