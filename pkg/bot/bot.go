package bot

import (
	"github.com/enmand/quarid-go/pkg/adapter"
	"github.com/enmand/quarid-go/pkg/irc"
	"github.com/satori/go.uuid"
)

// Manager manages adapters and events
type Manager struct {
	Identity uuid.UUID
	Name     string
}

// NewManager returns a new instance of a Manager
func NewManager(nickname string) Manager {
	return Manager{
		Identity: uuid.NewV4(),
		Name:     nickname,
	}
}

func (b adapter.Manager) IRC(server, user string, tls, tlsVerify bool) adapter.Adapter {
	bot := &ircbot{}

	bot.server = server
	bot.IRC = irc.NewClient(b.Name, user, tls, tlsVerify)

	return bot
}
