GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY_NAME=screen-monitor

all: deps build run
build:
	$(GOBUILD)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
deps:
	$(GOGET) github.com/mackerelio/go-osstat/memory
	$(GOGET) github.com/mackerelio/go-osstat/cpu
	$(GOGET) github.com/mackerelio/go-osstat/uptime
	$(GOGET) github.com/mackerelio/go-osstat/disk
	$(GOGET) github.com/mackerelio/go-osstat/network
	$(GOGET) github.com/mackerelio/go-osstat/loadavg
	$(GOGET) github.com/gorilla/mux
