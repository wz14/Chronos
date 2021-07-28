package config

import (
	"acc/crypto"
	"acc/idchannel"
	"acc/logger"
	"acc/pb"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"time"
)

type RemoteStart struct {
	c   Config
	nig *idchannel.NodeIDGroup
	pig *idchannel.PIDGroup
	l   *logger.Logger
	cc  *crypto.CCconfig
	e   *crypto.TPKE
	f   func(s Start)
}

func (s *RemoteStart) Getnig() *idchannel.NodeIDGroup {
	return s.nig
}

func (s *RemoteStart) Getpig() *idchannel.PIDGroup {
	return s.pig
}

func (s *RemoteStart) GetConfig() *Config {
	return &s.c
}

func (s *RemoteStart) GetCConfig() *crypto.CCconfig {
	return s.cc
}

func (s *RemoteStart) GetEConfig() *crypto.TPKE {
	return s.e
}

func NewRemoteStart(f func(s Start), c Config) *RemoteStart {
	logg := logger.NewLoggerWithID("config", c.MyID)
	ccConfig := crypto.CCconfig{}
	err := ccConfig.UnMarshal(c.CCconfig)
	if err != nil {
		logg.Fatalf("fail to unmarshal cc_config: %s", err.Error())
	}
	eConfig := crypto.TPKE{}
	err = eConfig.UnMarshal(c.Econfig)
	if err != nil {
		logg.Fatalf("fail to unmarshal e_config: %s", err.Error())
	}
	rs := &RemoteStart{
		c:   c,
		nig: nil,
		pig: nil,
		l:   logg,
		cc:  &ccConfig,
		e:   &eConfig,
		f:   f,
	}
	return rs
}

func (s *RemoteStart) Run() {
	s.HonestRun()
}

func (s *RemoteStart) HonestRun() {

	// TODO: wait others node start, e.g. 10s in config file
	// build id system for self
	// 1. start local grpc server for receive message
	// 2. start classifier in id channel
	// 3. wait some time for others' nodes start
	// 4. create id group for all nodes

	// start log
	localaddress := s.c.IpList[s.c.MyID] + ":" + strconv.Itoa(s.c.PortList[s.c.MyID])
	s.l.Infof("honest node id:%d in %s", s.c.MyID, localaddress)

	// init pid group
	var err error
	s.pig, err = idchannel.NewPIDGroup(&s.c)
	if err != nil {
		s.l.Fatal("primitive group create fail: %s", err.Error())
	}

	//create server && bind classifier
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(s.c.PortList[s.c.MyID]))
	if err != nil {
		s.l.Fatalf("tcp port open fail: %s in %d server", err, s.c.MyID)
	}
	defer lis.Close()
	server := grpc.NewServer()
	pb.RegisterNodeConServer(server, idchannel.NewClassifier(&s.c, s.c.MyID, s.pig))
	go server.Serve(lis)

	// wait time may config in file
	time.Sleep(time.Duration(s.c.WaitTime) * time.Second)

	// create id group
	s.nig, err = idchannel.NewIDGroup(&s.c)
	if err != nil {
		s.l.Fatalf("group create fail: %s", err.Error())
	}

	//run a consensus protocol, here is a send protocol for example
	// generate Data random
	s.f(s)
}
