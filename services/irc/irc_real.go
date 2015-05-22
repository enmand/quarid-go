// IRC client services in Golang
//
// This package implements an simple IRC service, that can be used in Golang to
// build IRC clients, bots, or other tools.
//
// Connecting
// To connect
package irc

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"net/textproto"
	"strings"
)

// Dial a connect to the IRC server. Use the form address:port
func (i *IRCClient) Connect(server string) error {
	var err error

	i.events = make(chan *IRCEvent)
	i.alive = make(chan bool)

	if !i.TLS {
		i.conn, err = net.DialTimeout("tcp", server, TIMEOUT)
	} else {
		i.conn, err = tls.DialWithDialer(&net.Dialer{
			Timeout: TIMEOUT,
		}, "tcp", server, &tls.Config{
			InsecureSkipVerify: i.TLSVerify,
		})
	}

	go i.authenticate()

	return err
}

func (i *IRCClient) Disconnect() error {
	err := i.Write(&IRCEvent{
		Command: IRC_QUIT,
	})

	i.alive <- false
	close(i.alive)
	close(i.events)

	return err
}

func (i *IRCClient) Write(ev *IRCEvent) error {
	var payload [][]byte

	payload = append(payload, []byte(ev.Command))
	for i, p := range ev.Parameters {
		if i == len(ev.Parameters)-1 && len(ev.Parameters) > 1 {
			p = fmt.Sprintf(":%s\r\n", p)
		}
		payload = append(payload, []byte(p))
	}

	payload = append(payload, []byte("\r\n"))
	full := bytes.Join(payload, []byte(" "))

	_, err := i.conn.Write(full)
	return err
}

func (i *IRCClient) authenticate() {
	var err error

	err = i.Write(&IRCEvent{
		Command: IRC_NICK,
		Parameters: []string{
			i.Nick,
		},
	})

	err = i.Write(&IRCEvent{
		Command: IRC_USER,
		Parameters: []string{
			i.Ident,
			"0.0.0.0",
			"0.0.0.0",
			i.Ident,
			i.Nick,
		},
	})

	if err != nil {
		i.Disconnect()
	}
}

func (i *IRCClient) read() {
	r := bufio.NewReader(i.conn)
	tp := textproto.NewReader(r)

	for {
		l, _ := tp.ReadLine()
		ws := strings.Split(l, " ")

		ev := &IRCEvent{}

		if prefix := ws[0]; prefix[0] == ':' {
			ev.Prefix = prefix[1:]
		} else {
			ev.Prefix = ""
			ev.Command = prefix
		}

		trailingIndex := 1
		if ev.Prefix != "" {
			trailingIndex = 2
			ev.Command = ws[1]
		}

		var trailing []string
		for _, param := range ws[trailingIndex:len(ws)] {
			if len(param) > 0 && (param[0] == ':' || len(trailing) > 0) {
				if param[0] == ':' {
					param = param[1:]
				}
				trailing = append(trailing, param)
			} else if len(trailing) == 0 {
				ev.Parameters = append(ev.Parameters, param)
			}
		}

		ev.Parameters = append(ev.Parameters, strings.Join(trailing, " "))

		if ev.Command == IRC_PING {
			ev.Command = IRC_PONG
			i.Write(ev)
		}

		i.events <- ev
	}
}

func scanEvents(data []byte, eof bool) (int, []byte, error) {
	if eof {
		return len(data), data[0:len(data)], nil
	}

	if i := bytes.Index(data, []byte("\n")); i >= 0 {
		return i + 2, data[0:i], nil
	}

	return 0, nil, nil
}
