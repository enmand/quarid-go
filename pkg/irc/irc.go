// Package irc provides IRC client services in Golang
//
// About
//
// This package implements an simple IRC service, that can be used in Golang to
// build IRC clients, bots, or other tools.
//
// See also: https://tools.ietf.org/html/rfc2812
package irc

import (
	"fmt"
	"net"
	"time"

	"github.com/enmand/quarid-go/pkg/adapter"
	"github.com/enmand/quarid-go/pkg/logger"
	"github.com/renstrom/shortuuid"
)

// TIMEOUT is the connection timeout to the IRC server
const TIMEOUT = 1 * time.Minute

// IRC is the IRC client interface
type IRC interface {
	// Connect to an IRC server. Use the form address:port
	Connect(server string) error

	// Disconnect from an IRC server
	Disconnect() error

	// Read blocks while reading from the server
	Read(n int)

	adapter.EventsHandler
	adapter.Responder
}

// Client is the implementation of the IRC interface
type Client struct {
	// The client's nickname on the server
	Nick string

	// The client's Ident on the server
	Ident string

	// The client's hostname
	Host string

	// The client's masked hostname on the server (if masked)
	MaskedHost string

	// Server is the server name the Client is connecting to
	Server string

	// If this connection is a TLS connection
	TLS bool

	// Should this client verify the server's SSL certs
	TLSVerify bool

	// handlers for filtered events
	handlers []*adapter.Handler

	// Events broadcasted from the server
	events chan *adapter.Event

	// The network connection this client has to the server
	conn net.Conn
}

// NewClient returns a new IRC client
func NewClient(nick, ident string, tlsverify, tls bool) *Client {
	c := &Client{
		Nick:      nick,
		Ident:     ident,
		TLSVerify: tlsverify,
		TLS:       tls,
	}

	c.Handle(
		[]adapter.Filter{CommandFilter{Command: IRC_PING}},
		func(ev *adapter.Event, c adapter.Responder) {
			ev.Command = IRC_PONG
			c.Write(ev)
		},
	)

	c.Handle(
		[]adapter.Filter{CommandFilter{Command: CONNECTED}},
		func(ev *adapter.Event, r adapter.Responder) {
			c.authenticate(r)
		},
	)

	c.Handle(
		[]adapter.Filter{CommandFilter{Command: IRC_ERR_NICKNAMEINUSE}},
		func(ev *adapter.Event, r adapter.Responder) {
			c.Nick = fixNick(c.Nick, r)
			writeNick(c.Nick, r)
		},
	)

	return c
}

func (i *Client) authenticate(c adapter.Responder) {
	logger.Log.Infof("Authenticating for nick %s!%s", i.Nick, i.Ident)

	writeNick(i.Nick, c)

	// RFC 2812 USER command
	c.Write(&adapter.Event{
		Command: IRC_USER,
		Parameters: []string{
			i.Ident,
			"0",
			"*",
			i.Nick,
		},
	})
}

func fixNick(nick string, c adapter.Responder) string {
	uniq := shortuuid.UUID()

	newNick := fmt.Sprintf("%s_%s", nick, uniq)
	logger.Log.Debugf("Fixing nick to %s", nick)

	return newNick[:9] // minimum max length in 9
}

func writeNick(nick string, c adapter.Responder) {
	c.Write(&adapter.Event{
		Command: IRC_NICK,
		Parameters: []string{
			nick,
		},
	})

	timeout := make(chan bool)
	go func() {
		// Wait 400ms for an error
		time.Sleep(400 * time.Millisecond)
		timeout <- true
	}()
	<-timeout
}
