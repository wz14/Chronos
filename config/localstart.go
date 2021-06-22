package config

import (
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

func NewLocalStart() *LocalStart {
	lo := logger.NewLogger("main")
	c, err := NewConfig("./config.yaml") // get a no-pointer config
	if err != nil {
		lo.Fatalf("read config fail: %s", err.Error())
	}
	return &LocalStart{
		c: c,
		l: logger.NewLogger("main"),
	}
}

type LocalStart struct {
	c   Config
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

func (s *LocalStart) GetConfig() *Config {
	return &s.c
}

func (s *LocalStart) CopySelf(id int) LocalStart {
	newc := s.c
	newc.MyID = id
	return LocalStart{c: newc}
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
	/*
		if s.c.MyID == 1 {
			err := core.Send(&pb.Message{
				Id:       "Send_1",
				Sender:   1,
				Receiver: 2,
				Data:     []byte("what are you doing"),
			}, s)
			if err != nil {
				s.l.Fatalf("send error: %s", err.Error())
			}
			s.l.Print("I'm id 1 node, I send: what are you doing")
		}

		if s.c.MyID == 2 {
			m, err := core.Receive(s.pig.GetRootPID("Send_1"))
			if err != nil {
				s.l.Fatal("receive error")
			}
			s.l.Print("I'm id 2 node , I receive", m)
		}
	*/
	wg.Done()
}
