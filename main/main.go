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
		log.Fatal("not be implemented yet")
	} else if os.Args[1] == "remote" {
		log.Fatal("not be implemented yet")
	} else if os.Args[1] == "help" {
		fmt.Println(usage)
	} else {
		lo.Fatal(usage)
	}
}
