package adapter

import (
	"fmt"

	"github.com/quarid-go/pkg/config"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type Matrix struct {
	c *mautrix.Client
}

func NewMatrix(homeserver, uid, token string) (Adapter, error) {
	id := id.UserID(uid)
	c, err := mautrix.NewClient(homeserver, id, token)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to homeserver: %w", err)
	}
	c.UserAgent = fmt.Sprintf("Go-adapter-matrix/%d", config.Version)

	return &Matrix{c}, nil
}

func (m *Matrix) Start() error {
	s := m.c.Syncer.(*mautrix.DefaultSyncer)
	s.OnEvent(func(source mautrix.EventSource, evt *event.Event) {
		fmt.Printf("\n\n")
		fmt.Printf("source: %s\n", source)
		fmt.Printf("evt: %# v\n", evt)
	})
	return m.c.Sync()
}

func (m *Matrix) Stop() error {
	m.c.StopSync()
	return nil
}

func (m *Matrix) Handle(f []Filter, h HandlerFunc) {
	s := m.c.Syncer.(*mautrix.DefaultSyncer)
	_ = s
}

func (m *Matrix) Write(ev *Event) error {
	return nil
}
