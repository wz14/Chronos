package main

import (
	"acc/Benchmark"
	"acc/config"
	"acc/logger"
	"acc/pb"
	"acc/rbc"
	"crypto/rand"
)

var log = logger.NewLogger("main")

func f(s config.Start) {
	//run a consensus protocol, here is a send protocol for example
	// generate Data random
	conf := s.GetConfig()
	if conf.MyID == 1 {
		TX := 250
		num := conf.Txnum
		data := make([]byte, num*TX)
		n, err := rand.Read(data)
		if err != nil {
			log.Fatalf("random read data fail: %s", err.Error())
		}
		if n != num*TX {
			log.Fatal("fail to read enough bytes")
		}
		log.Infof("input message's length is %d", len(data))
		Benchmark.Begin("testRBC", conf.MyID)
		go rbc.RBCBroadcast(&pb.Message{
			Id:       "testRBC",
			Sender:   1,
			Receiver: 0,
			Data:     data,
		}, s)
	}

	_, err := rbc.RBCDeliver("testRBC", s)
	if err != nil {
		log.Error("RBC deliver test fail: %s", err.Error())
	}
	Benchmark.End("testRBC", conf.MyID)
}
