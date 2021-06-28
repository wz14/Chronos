package config

import (
	"acc/crypto"
	"acc/idchannel"
)

// start is a runtime config
type Start interface {
	Run()
	Getnig() *idchannel.NodeIDGroup
	Getpig() *idchannel.PIDGroup
	GetConfig() *Config
	GetCConfig() *crypto.CCconfig
}
