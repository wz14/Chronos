package main

import (
	"acc/benchmark"
	"acc/config"
	"acc/consensus"
	"acc/logger"
	"acc/pb"
	"crypto/rand"
	"fmt"
	"os"
	"time"
)

const usage = "usage: ./start local|remoteGen|remote|help"

var log = logger.NewLogger("main")

func main() {
	lo := logger.NewLogger("main")
	if len(os.Args) != 2 {
		lo.Fatal(usage)
	}

	c, err := config.NewConfig("./config.yaml", true) // get a no-pointer config
	if err != nil {
		lo.Fatalf("read config fail: %s", err.Error())
	}

	benchmark.InitBenchmark(c)

	if os.Args[1] == "local" {
		f := makeFunc(c)
		start := config.NewLocalStart(f, c)
		start.Run()
		_ = benchmark.BenchmarkOuput()
	} else if os.Args[1] == "remoteGen" {
		err := c.RemoteGen(".")
		if err != nil {
			log.Fatal("remoteGen fail: %s", err.Error())
		}
		log.Info("remoteGen ok !")
	} else if os.Args[1] == "remote" {
		if c.Isremote != true {
			log.Fatal("the config.yaml is not remote config")
		}
		f := makeFunc(c)
		start := config.NewRemoteStart(f, c)
		start.Run()
		_ = benchmark.BenchmarkOuput()
		// TODO:config 1.server prepare time && 2.server wait time && 3. which consensus is used
		// TODO:rename main/localstart or delete and add to main
		// TODO:add RBC benchmark in CHOS
		fmt.Println("wait some time(120) for others aba is ok")
		time.Sleep(time.Second * time.Duration(c.PrepareTime))
	} else if os.Args[1] == "help" {
		fmt.Println(usage)
	} else {
		lo.Fatal(usage)
	}
}

func makeFunc(c config.Config) func(start config.Start) {
	var typ consensus.CONSTYPE
	var idName string
	if c.ConsenType == "CHOS" {
		idName = "chronos"
		typ = consensus.CHOSCon
	} else if c.ConsenType == "HB" {
		idName = "hb"
		typ = consensus.HBCon
	} else {
		log.Fatal("no such consensus type")
	}

	return func(s config.Start) {
		//run a consensus protocol, here is a send protocol for example
		// generate Data random
		conf := s.GetConfig()
		benchmark.Create(idName)
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
		benchmark.Begin(idName, conf.MyID)

		results, _, err := consensus.Consensus(typ, &pb.Message{
			Id:       idName,
			Sender:   uint32(conf.MyID),
			Receiver: 0,
			Data:     data,
		}, s)

		consenedTXnum := len(results) * len(results[0].Data) / TX
		log.Infof("results length: %d, tx nums: %d", len(results), consenedTXnum)
		benchmark.End(idName, conf.MyID)
		benchmark.Nums(idName, conf.MyID, consenedTXnum)
	}
}
