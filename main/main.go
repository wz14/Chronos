package main

import (
	"acc/benchmark"
	"acc/config"
	"acc/logger"
	"os"
)

const usage = "./start local|remoteGen|remote"

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
	} else {
		lo.Fatal(usage)
	}
}
