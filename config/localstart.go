package config

import (
	"acc/crypto"
	"acc/idchannel"
	"acc/logger"
	"acc/pb"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

func NewLocalStartWithReadLocalConfig(f func(s Start), configpath string) *LocalStart {
	lo := logger.NewLogger("config")
	c, err := NewConfig(configpath, true) // get a no-pointer config
	if err != nil {
		lo.Fatalf("read config fail: %s", err.Error())
	}
	return &LocalStart{
		c: c,
		l: logger.NewLogger("config"),
		f: f,
	}
}

func NewLocalStart(f func(s Start), c Config) *LocalStart {
	return &LocalStart{
		c: c,
		l: logger.NewLogger("config"),
		f: f,
	}
}

type LocalStart struct {
	c   Config
	nig *idchannel.NodeIDGroup
	pig *idchannel.PIDGroup
	l   *logger.Logger
	cc  *crypto.CCconfig
	e   *crypto.TPKE
	f   func(s Start)
}

func (s *LocalStart) Getnig() *idchannel.NodeIDGroup {
	return s.nig
}

func (s *LocalStart) Getpig() *idchannel.PIDGroup {
	return s.pig
}

func (s *LocalStart) GetConfig() *Config {
	return &s.c
}

func (s *LocalStart) GetCConfig() *crypto.CCconfig {
	return s.cc
}

func (s *LocalStart) GetEConfig() *crypto.TPKE {
	return s.e
}

func (s *LocalStart) CopySelf(id int, cconfig *crypto.CCconfig, tpke *crypto.TPKE) LocalStart {
	newc := s.c
	newc.MyID = id
	return LocalStart{
		c:  newc,
		l:  logger.NewLoggerWithID("main", id),
		f:  s.f,
		cc: cconfig,
		e:  tpke,
	}
}

func (s *LocalStart) Run() {

	// init common coin
	cconfigs, err := crypto.NewCCconfigs(s.c.F, s.c.N)
	if err != nil {
		s.l.Fatalf("common coin init fail: %s", err.Error())
	}

	// init tpke
	tpkes := crypto.NewTPKE(s.c.N, s.c.F)

	if s.c.Isremote {
		s.l.Fatal("no implement remote deployment setting")
	}
	wg.Add(s.c.N)
	for i := 0; i < s.c.N; i++ {
		news := s.CopySelf(i, cconfigs[i], tpkes[i])
		go news.HonestRun()
	}
	wg.Wait()
}

func (s *LocalStart) HonestRun() {

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
	time.Sleep(2 * time.Second)

	// create id group
	s.nig, err = idchannel.NewIDGroup(&s.c)
	if err != nil {
		s.l.Fatalf("group create fail: %s", err.Error())
	}

	//run a consensus protocol, here is a send protocol for example
	// generate Data random
	s.f(s)
	wg.Done()
}
