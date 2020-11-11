package bot

import (
	"github.com/enmand/quarid-go/internal/adapter"
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
func (b *Bot) Start(errCh chan error) {
	for _, a := range b.adapters {
		go func(a adapter.Adapter, ch chan error) {
			ch <- a.Start()
		}(a, errCh)
	}
}

// Stop the bot and it's adapters
func (b *Bot) Stop(errCh chan error) {
	for _, a := range b.adapters {
		errCh <- a.Stop()
	}
}
