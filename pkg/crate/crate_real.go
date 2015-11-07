package crate

import "github.com/enmand/quarid-go/vm"

// RegisterVM registers a VM to a Crate. If the a VM with the same Type is
// already registered, it will return the registered VM
func (c *Crate) RegisterVM(virtualMachine vm.VM) vm.VM {
	if v, ok := c.vms[virtualMachine.Type()]; !ok {
		return v
	}

	c.vms[virtualMachine.Type()] = virtualMachine

	return virtualMachine
}
