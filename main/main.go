package main

import (
	"acc/benchmark"
	"acc/config"
	"acc/logger"
	"fmt"
	"os"
)

const usage = "usage: ./start local|remoteGen|remote|help"

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
		start := config.NewRemoteStart(f, c)
		start.Run()
		_ = benchmark.BenchmarkOuput()
	} else if os.Args[1] == "help" {
		fmt.Println(usage)
	} else {
		lo.Fatal(usage)
	}
}
