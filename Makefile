GO=go
GOFLAGS=-race
DEV_BIN=bin

all: build

.PHONY: build
build:
	go build $(GOFLAGS) -o $(DEV_BIN)/grg2p main.go

