VERSION := $(shell git describe --tags)
REVISION := $(shell git rev-parse --short HEAD)

# enterprise or open
VARIANT?=open

BINARY_NAME := dcos-signal

LDFLAGS := -X github.com/mesosphere/dcos-signal/signal.VERSION=$(VERSION) -X github.com/mesosphere/dcos-signal/signal.REVISION=$(REVISION) -X github.com/mesosphere/dcos-signal/config.VARIANT=$(VARIANT)

FILES := $(shell go list ./... | grep -v vendor)

# Testing Local Run 
ANON_PATH?=/tmp/anon-id.json
HOST?=localhost
CONFIG?=/tmp/signal-config.json
EXTRA?=

all: test install

test:
	@echo "+$@"
	go test $(FILES)  -cover

build: 
	@echo "+$@"
	go build -v -ldflags '$(LDFLAGS)' $(FILES)

install:
	@echo "+$@"
	go install -v -ldflags '$(LDFLAGS)' $(FILES)

run:
	@echo "+$@"
	go run dcos_signal.go -v -anonymous-id-path $(ANON_PATH) -report-host $(HOST) -report-port 1050 -c $(CONFIG) $(EXTRA)
