package main

import (
	"acc/benchmark"
	"acc/config"
	"acc/consensus"
	"acc/logger"
	"acc/pb"
	"crypto/rand"
)

var log = logger.NewLogger("main")

func f(s config.Start) {
	//run a consensus protocol, here is a send protocol for example
	// generate Data random
	conf := s.GetConfig()
	benchmark.Create("CHOS")
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
	benchmark.Begin("CHOS", conf.MyID)

	results, _, err := consensus.Consensus(consensus.CHOSCon, &pb.Message{
		Id:       "chronos",
		Sender:   uint32(conf.MyID),
		Receiver: 0,
		Data:     data,
	}, s)

	consenedTXnum := len(results) * len(results[0].Data) / TX
	log.Infof("results length: %d, tx nums: %d", len(results), consenedTXnum)
	benchmark.End("CHOS", conf.MyID)
	benchmark.Nums("CHOS", conf.MyID, consenedTXnum)
}
