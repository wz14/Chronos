all: build, test

.PHONY: build
build:
	go build -o build/start ./main

.PHONY: test
test:
	go test ./main
	go test ./core

.PHONY: clean
clean:
	rm build/*
