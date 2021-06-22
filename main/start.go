package main

import (
	"acc/core"
	"acc/idchannel"
	"acc/logger"
	"acc/pb"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"sync"
	"time"
)

var log = logger.NewLogger("main")
var wg = sync.WaitGroup{}

func NewStart() *Start {
	c, err := NewConfig() // get a no-pointer config
	if err != nil {
		log.Fatalf("read config fail: %s", err.Error())
	}
	return &Start{c: c}
}

type Start struct {
	c   Config
	nig *idchannel.NodeIDGroup
	pig *idchannel.PIDGroup
}

func (s *Start) CopySelf(id int) Start {
	newc := s.c
	newc.MyID = id
	return Start{c: newc}
}

func (s *Start) Run() {
	if s.c.Isremote {
		log.Fatal("no implement remote deployment setting")
	}
	wg.Add(s.c.N)
	for i := 0; i < s.c.N; i++ {
		news := s.CopySelf(i)
		go news.HonestRun()
	}
	wg.Wait()
}

func (s *Start) HonestRun() {
	// TODO: wait others node start, e.g. 10s in config file
	// build id system for self
	// 1. start local grpc server for receive message
	// 2. start classifier in id channel
	// 3. wait some time for others' nodes start
	// 4. create id group for all nodes

	localaddress := s.c.IpList[s.c.MyID] + ":" + strconv.Itoa(s.c.PortList[s.c.MyID])
	log.Infof("honest node id:%d in %s", s.c.MyID, localaddress)
	//create server & classifer
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(s.c.PortList[s.c.MyID]))
	if err != nil {
		log.Fatalf("tcp port open fail: %s in %d server", err, s.c.MyID)
	}
	server := grpc.NewServer()

	s.pig, err = idchannel.NewPIDGroup(&s.c)
	if err != nil {
		log.Fatal("primitive group create fail: %s", err.Error())
	}

	pb.RegisterNodeConServer(server, idchannel.NewClassifier(&s.c, s.c.MyID, s.pig))
	go server.Serve(lis)

	// wait time
	time.Sleep(5 * time.Second)

	// create id group
	s.nig, err = idchannel.NewIDGroup(&s.c)
	if err != nil {
		log.Fatalf("group create fail: %s", err.Error())
	}

	//TODO: run a consensus protocol

	if s.c.MyID == 1 {
		err := core.Send(pb.Message{
			Id:       "Send_1",
			Sender:   1,
			Receiver: 2,
			Data:     []byte("what are you doing"),
		}, s.pig.GetRootPID("Send_1"), s.nig)
		if err != nil {
			log.Fatalf("send error: %s", err.Error())
		}
		log.Print("I'm id 1 node, I send: what are you doing")
	}

	if s.c.MyID == 2 {
		m, err := core.Receive(s.pig.GetRootPID("Send_1"))
		if err != nil {
			log.Fatal("receive error")
		}
		log.Print("I'm id 2 node , I receive", m)
	}

	wg.Done()
}

func main() {
	log.Info("asynchronous agreement components start")
	start := NewStart()
	start.Run()
	log.Info("asynchronous agreement components stop here")
}
