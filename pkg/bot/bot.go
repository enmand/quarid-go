package bot

import (
	"github.com/enmand/quarid-go/pkg/adapter"
	"github.com/satori/go.uuid"
)

// Bot manages adapters and events
type Bot struct {
	Identity uuid.UUID
	adapters []adapter.Adapter
}

// New returns a new instance of a Bot
func New(a []adapter.Adapter) *Bot {
	return &Bot{
		Identity: uuid.NewV4(),

		adapters: a,
	}
}

// Start the bot and it's adapters
func (b *Bot) Start() chan error {
	errCh := make(chan error)
	for _, a := range b.adapters {
		errCh <- a.Start()
	}

	return errCh
}

// Stop the bot and it's adapters
func (b *Bot) Stop() chan error {
	errCh := make(chan error)
	for _, a := range b.adapters {
		errCh <- a.Stop()
	}

	return errCh
}
