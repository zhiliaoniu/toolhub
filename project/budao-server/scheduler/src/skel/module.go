package skel

import (
	"service"
)

type Module interface {
	Init(*service.ServiceContext) error
	Start() error
}

type ModuleBase struct {
}

func (m *ModuleBase) Start() error {
	return nil
}

func (sk *Skel) RegisterModule(name string, m Module) {
	sk.modules[name] = m
}
