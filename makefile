all: build, test

.PHONY: build
build:
	go build -o build/start ./main

# need to be serial
.PHONY: buildtest
buildtest:
	go test -c -o build/config.test ./config
	go test -c -o build/core.test ./core
	go test -c -o build/rbc.test ./rbc
	go test -c -o build/aba.test ./aba

.PHONY: runtest
runtest:
	AACDEUBG=1 && cd ./build && ./config.test
	AACDEUBG=1 && cd ./build && ./core.test
	AACDEUBG=1 && cd ./build && ./rbc.test
	AACDEUBG=1 && cd ./build && ./aba.test

.PHONY: clean
clean:
	rm build/*
