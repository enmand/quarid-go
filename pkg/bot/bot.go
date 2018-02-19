package bot

import (
	"github.com/enmand/quarid-go/pkg/adapter"
	"github.com/satori/go.uuid"
)

// Manager manages adapters and events
type Bot struct {
	Identity uuid.UUID
	adapter  adapter.Adapter
}

// NewManager returns a new instance of a Manager
func New(a adapter.Adapter) *Bot {
	return &Bot{
		Identity: uuid.NewV4(),

		adapter: a,
	}
}

func (b *Bot) Start() error {
	return b.adapter.Start()
}

func (b *Bot) Stop() error {
	return b.adapter.Stop()
}
