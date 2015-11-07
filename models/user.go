package models

import (
	"fmt"

	"github.com/enmand/quarid-go/pkg"
	"github.com/enmand/quarid-go/pkg/crate"
	"github.com/enmand/quarid-go/pkg/database"
	"github.com/enmand/quarid-go/vm"
)

// User is a user that has access to the system
type User struct {
	Name  string
	Email string

	Crates []*crate.Crate
}

// NewUser creates a User
func NewUser(name, email string) *User {
	c := services.GetConfig()
	d := database.Open(c.GetString("database"))

	return &Crates{
		Name: name,
	}
}

// Register a Crate to our Crates registry, with a name. If this Crate already
// exists, return the existing Cratem and
func (cs *User) Register(name string, crate vm.Crate) (*vm.Crate, error) {
	c, err := cs.Find(name)
	if err != nil {
		return c, fmt.Errorf("A crate named %s is already registered", name)
	}

	cs.Crates = append(cs.Crates, &crate)

	return &crate, nil
}

// Find the Crate with the given name, if it exists
func (cs *User) Find(name string) (*vm.Crate, error) {
	for _, c := range cs.Crates {
		if c.Name == name {
			return c, nil
		}
	}

	return nil, fmt.Errorf("No Crate named %s found", name)
}
