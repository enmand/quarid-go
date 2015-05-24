//
// IRC client services in Golang
//
// About
//
// This package implements an simple IRC service, that can be used in Golang to
// build IRC clients, bots, or other tools.
//
// See also: https://tools.ietf.org/html/rfc2812
package irc

import (
	"net"
	"time"

	"github.com/enmand/quarid-go/services/config"
)

const TIMEOUT = 1 * time.Minute

type IRC interface {
	// Connect to an IRC server. Use the form address:port
	Connect(server string) error

	// Disconnect from an IRC server
	Disconnect() error

	// Loop blocks while reading from the server
	Loop()

	EventsHandler
	Responder
}

type Client struct {
	// The client's nickname on the server
	Nick string

	// The client's Ident on the server
	Ident string

	// The client's hostname
	Host string

	// The client's masked hostname on the server (if masked)
	MaskedHost string

	// If this connection is a TLS connection
	TLS bool

	// Should this client verify the server's SSL certs
	TLSVerify bool

	// Dead is blocks until the conn
	dead chan bool

	// Events broadcasted from the server
	events chan *Event

	// The network connection this client has to the server
	conn net.Conn
}

func NewClient(c *config.Config) IRC {
	return &Client{
		Nick:      c.Nick,
		Ident:     c.Ident,
		TLSVerify: c.TLS.Verify,
		TLS:       c.TLS.Enable,
	}
}
