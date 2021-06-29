package main

import (
	"acc/benchmark"
	"acc/config"
)

func main() {
	// command process or something
	start := config.NewLocalStart(f, "./config.yaml")
	start.Run()
	_ = benchmark.BenchmarkOuput()
}
