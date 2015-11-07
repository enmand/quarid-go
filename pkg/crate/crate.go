package crate

import (
	"github.com/enmand/quarid-go/pkg/adapter"
	"github.com/enmand/quarid-go/pkg/database"
	"github.com/enmand/quarid-go/pkg/plugin"
	"github.com/satori/go.uuid"
)

const (
	// CrateStateRunning represents a "running", acive crate
	CrateStateRunning = "running"

	// CrateStateStopped represents a "stopped" crate
	CrateStateStopped = "stopped"
)

// A Crate defines a set of VMs with a specific name
type Crate struct {
	ID    uuid.UUID
	Name  string
	State string

	Adapter adapter.Adapter

	plugins map[string]plugin.Plugin
	db      database.VMDatabase
}

// NewCrate creates a Crate for VMs
func NewCrate(name string, a adapter.Adapter) *Crate {
	id := uuid.NewV4()

	return &Crate{ID: id, Name: name, State: CrateStateStopped}
}
