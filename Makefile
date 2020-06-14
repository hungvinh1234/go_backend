.PHONY: build run

build:
	go build -ldflags "-s -w" -o ./bin/server ./cmd/main.go

run:
	@GO111MODULE=off go get -v github.com/phanletrunghieu/fresh
	DEBUG=true PORT=3000 fresh -c fresh.conf
