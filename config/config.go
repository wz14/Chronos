package config

import (
	"acc/crypto"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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
	Txnum    int      `yaml:"Txnum"`
	// judge if execute read config function before
	// default is false in golang structure declare
	isRead    bool
	MyID      int      `yaml:"MyID"`
	Statistic string   `yaml:"Statistic"`
	CCconfig  []string `yaml:"CCconfig"`
	Econfig   []string `yaml:"Econfig"`
	// server start time
	PrepareTime int `yaml:"PrepareTime"`
	WaitTime    int `yaml:"WaitTime"`
}

func NewConfig(configName string, isLocal bool) (Config, error) {
	c := Config{}
	err := c.ReadConfig(configName, isLocal)
	if err != nil {
		return Config{}, err
	}
	return c, err
}

// read config from ConfigName file location
func (c *Config) ReadConfig(ConfigName string, isLocal bool) error {
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

	if !isLocal {
		// id is begin from 0 to ... N-1
		if c.MyID >= c.N || c.MyID < 0 {
			return errors.New("ID is begin from 0 to N-1")
		}
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

func (c *Config) Marshal(location string) error {
	byts, err := yaml.Marshal(c)
	if err != nil {
		return errors.Wrap(err, "marshal config fail")
	}
	err = ioutil.WriteFile(location, byts, 0777)
	if err != nil {
		return errors.Wrap(err, "marshal config fail")
	}
	return nil
}

func (c *Config) RemoteGen(dir string) error {
	c.Isremote = true
	// myid; CCconfig; TPKconfig
	ccconfigs, err := crypto.NewCCconfigs(c.F, c.N)
	if err != nil {
		return errors.Wrap(err, "generate cc_config fail")
	}
	econfigs := crypto.NewTPKE(c.N, c.F)
	for i := 0; i < c.N; i++ {
		c.MyID = i
		ccM, err := ccconfigs[i].Marshal()
		if err != nil {
			return errors.Wrap(err, "generate cc_config fail")
		}
		c.CCconfig = ccM
		eM, err := econfigs[i].Marshal()
		if err != nil {
			return errors.Wrap(err, "generate e_config fail")
		}
		c.Econfig = eM
		err = c.Marshal(dir + "/config_" + strconv.Itoa(i) + ".yaml")
		if err != nil {
			return errors.Wrap(err, "marshal config fail")
		}
	}
	return nil
}
