package main

import (
	"acc/Benchmark"
)

func main() {
	// command process or something
	start := NewLocalStart()
	start.Run()
	Benchmark.BenchmarkOuput()
}
