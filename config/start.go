package config

import (
	"acc/idchannel"
)

type Start interface {
	Run()
	Getnig() *idchannel.NodeIDGroup
	Getpig() *idchannel.PIDGroup
	GetConfig() *Config
}

func NewStart() Start {
	return NewLocalStart()
}
