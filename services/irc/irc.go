package irc

import (
	"net"
	"time"

	"github.com/enmand/quarid-go/services/config"
)

const TIMEOUT = 1 * time.Minute

type IRC interface {
	// Connect to an IRC server,
	Connect(server string) error

	// Disconnect from an IRC server
	Disconnect() error

	// Write to the server
	Write(ev *IRCEvent) error
}

type IRCEvent struct {
	// The event prefix (optional in spec)
	Prefix string

	// The command that the client (or server) is sending/sent
	Command string

	// The parameters to the command the client (or server) is sending/sent
	Parameters []string
}

type IRCClient struct {
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

	// Addresses
	Addrs map[string]net.Addr

	// Alive while the connective is still active
	alive chan bool

	// Events broadcasted from the server
	events chan *IRCEvent

	// The network connection this client has to the server
	conn net.Conn
}

func NewClient(c *config.Config) IRC {
	return &IRCClient{
		Nick:      c.Nick,
		Ident:     c.Ident,
		TLSVerify: c.TLS.Verify,
		TLS:       c.TLS.Enable,
	}
}
