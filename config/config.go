package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
)

var ConfigReadError = errors.New("config read fail, check config.yaml in root directory")
var NotReadFileError = errors.New("execute run ReadConfig function before querying config")
var NotDefined = errors.New("this item in config is allowed to omit")

// Implement Config interface in local linux machine setting
type Config struct {
	N int `yaml:"N"`
	F int `yaml:"F"`
	// ip list is for
	IpList   []string `yaml:"IpList"`
	PortList []int    `yaml:"PortList"`
	Isremote bool     `yaml:"Isremote"`
	// judge if execute read config function before
	// default is false in golang structure declare
	isRead bool
	MyID   int
}

func NewConfig(configname string) (Config, error) {
	ConfigName := configname
	c := Config{}
	err := c.ReadConfig(ConfigName)
	if err != nil {
		return Config{}, err
	}
	return c, err
}

// read config from ConfigName file location
func (c *Config) ReadConfig(ConfigName string) error {
	var id string
	var ok bool
	byt, err := ioutil.ReadFile(ConfigName)
	if err != nil {
		goto ret
	}

	err = yaml.Unmarshal(byt, c)
	if err != nil {
		goto ret
	}

	if c.N <= 0 || c.F < 0 {
		return errors.Wrap(errors.New("N or F is negative"),
			ConfigReadError.Error())
	}

	if c.N != len(c.IpList) || c.N != len(c.PortList) {
		return errors.Wrap(errors.New("ip list"+
			" length or port list length isn't match N"),
			ConfigReadError.Error())
	}
	c.isRead = true

	id, ok = os.LookupEnv("ID")
	if !ok {
		return errors.New("no ID environment in os")
	}

	c.MyID, err = strconv.Atoi(id)
	if err != nil {
		return errors.New("convert ID environment value to integer fail")
	}

	// id is begin from 0 to ... N-1
	if c.MyID >= c.N || c.MyID < 0 {
		return errors.New("ID is begin from 0 to N-1")
	}

	return nil
ret:
	return errors.Wrap(err, ConfigReadError.Error())
}

// Achieve numbers of total nodes
// the return value is a positive integer
func (c *Config) GetN() (int, error) {
	if !c.isRead {
		return 0, NotReadFileError
	}
	return c.N, nil
}

// Achieve number of corrupted nodes
// return value is a positive integer
func (c *Config) GetF() (int, error) {
	if !c.isRead {
		return 0, NotReadFileError
	}
	return c.F, nil
}

// Achieve ip list if defined
// return a ip list of defined ip in config file
func (c *Config) GetIPList() ([]string, error) {
	if !c.isRead {
		return nil, NotReadFileError
	}
	if len(c.IpList) == 0 {
		return nil, NotDefined
	}
	return c.IpList, nil
}

// Achieve port list if defined
// return a port list of defined port in config file
func (c *Config) GetPortList() ([]int, error) {
	if !c.isRead {
		return nil, NotReadFileError
	}
	if len(c.PortList) == 0 {
		return nil, NotDefined
	}
	return c.PortList, nil
}

func (c *Config) GetMyID() (int, error) {
	if !c.isRead {
		return 0, NotReadFileError
	}
	return c.MyID, nil
}
