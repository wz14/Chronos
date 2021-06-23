package main

import (
	"acc/Benchmark"
	"acc/config"
	"acc/idchannel"
	"acc/logger"
	"acc/pb"
	"acc/rbc"
	"crypto/rand"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

func NewLocalStart() *LocalStart {
	lo := logger.NewLogger("main")
	c, err := config.NewConfig("./config.yaml") // get a no-pointer config
	if err != nil {
		lo.Fatalf("read config fail: %s", err.Error())
	}
	return &LocalStart{
		c: c,
		l: logger.NewLogger("main"),
	}
}

type LocalStart struct {
	c   config.Config
	nig *idchannel.NodeIDGroup
	pig *idchannel.PIDGroup
	l   *logger.Logger
}

func (s *LocalStart) Getnig() *idchannel.NodeIDGroup {
	return s.nig
}

func (s *LocalStart) Getpig() *idchannel.PIDGroup {
	return s.pig
}

func (s *LocalStart) GetConfig() *config.Config {
	return &s.c
}

func (s *LocalStart) CopySelf(id int) LocalStart {
	newc := s.c
	newc.MyID = id
	return LocalStart{c: newc, l: logger.NewLoggerWithID("main", id)}
}

func (s *LocalStart) Run() {
	if s.c.Isremote {
		s.l.Fatal("no implement remote deployment setting")
	}
	wg.Add(s.c.N)
	for i := 0; i < s.c.N; i++ {
		news := s.CopySelf(i)
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
	if s.c.MyID == 1 {
		TX := 250
		num := s.c.Txnum
		data := make([]byte, num*TX)
		n, err := rand.Read(data)
		if err != nil {
			s.l.Fatalf("random read data fail: %s", err.Error())
		}
		if n != num*TX {
			s.l.Fatal("fail to read enough bytes")
		}
		s.l.Infof("input message's length is %d", len(data))
		Benchmark.Begin("testRBC", s.c.MyID)
		go rbc.RBCBroadcast(&pb.Message{
			Id:       "testRBC",
			Sender:   1,
			Receiver: 0,
			Data:     data,
		}, s)
	}

	_, err = rbc.RBCDeliver("testRBC", s)
	if err != nil {
		s.l.Error("RBC deliver test fail: %s", err.Error())
	}
	Benchmark.End("testRBC", s.c.MyID)
	wg.Done()
}
