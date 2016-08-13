package bot

import (
	"github.com/enmand/quarid-go/pkg/adapter"
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

// IRC returns an IRC connections as an adapter.Adapter
func (b Manager) IRC(server, user string, tls, tlsVerify bool) adapter.Adapter {
	return NewIRC(server, b.Name, user, tls, tlsVerify)
}
