package irc

import (
	"bufio"
	"crypto/tls"
	"net"
	"net/textproto"
	"strings"
	"time"
)

func (i *Client) Connect(server string) error {
	var err error

	i.events = make(chan *Event)
	i.dead = make(chan bool)

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

func (i *Client) Disconnect() error {
	err := i.Write(&Event{
		Command: IRC_QUIT,
	})

	i.dead <- false
	close(i.dead)
	close(i.events)

	return err
}

func (i *Client) authenticate() {
	var err error

	err = i.Write(&Event{
		Command: IRC_NICK,
		Parameters: []string{
			i.Nick,
		},
	})

	err = i.Write(&Event{
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

func (i *Client) read() {
	r := bufio.NewReader(i.conn)
	tp := textproto.NewReader(r)

	for {
		l, _ := tp.ReadLine()
		ws := strings.Split(l, " ")

		ev := &Event{}

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
		ev.Timestamp = time.Now()

		i.events <- ev

		if ev.Command == IRC_PING {
			ev.Command = IRC_PONG
			go i.Write(ev)
		}
	}
}
