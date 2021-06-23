all: build, test

.PHONY: build
build:
	go build -o build/start ./main

.PHONY: test
test:
	go test ./config
	go test ./core
	go test ./rbc

.PHONY: clean
clean:
	rm build/*
