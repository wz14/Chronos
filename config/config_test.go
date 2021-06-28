package config

import (
	"testing"
)

func TestConfig_ReadConfig(t *testing.T) {
	ConfigName := "./mock/config1.yaml"
	c := Config{}
	err := c.ReadConfig(ConfigName)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestConfig_ReadConfig2(t *testing.T) {
	ConfigName := "./mock/config1.yaml"
	c := Config{}
	err := c.ReadConfig(ConfigName)
	if err != nil {
		t.Error(err.Error())
	}

	if c.isRead != true {
		t.Error("isRead set error")
	}

	if c.Isremote != false {
		t.Error("isremote set error")
	}

	for _, ip := range c.IpList {
		if ip != "127.0.0.1" {
			t.Error("ip list read fail: ", c.IpList)
		}
	}

	if len(c.IpList) != 4 {
		t.Error("ip list read fail")
	}

	for i, port := range c.PortList {
		if port != 2001+i {
			t.Error("port list read fail")
		}
	}

	if len(c.PortList) != 4 {
		t.Error("port list read fail")
	}

	if c.N != 4 {
		t.Error(c.N)
	}

	if c.F != 1 {
		t.Error(c.F)
	}
}
