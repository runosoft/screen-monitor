GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get

all: deps build
build:
	$(GOBUILD)
deps:
	$(GOGET) github.com/mackerelio/go-osstat/memory
	$(GOGET) github.com/mackerelio/go-osstat/cpu
	$(GOGET) github.com/mackerelio/go-osstat/uptime
	$(GOGET) github.com/mackerelio/go-osstat/disk
	$(GOGET) github.com/mackerelio/go-osstat/network
	$(GOGET) github.com/mackerelio/go-osstat/loadavg
	$(GOGET) github.com/gorilla/mux
