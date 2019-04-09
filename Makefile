VERSION := $(shell git describe --tags)
REVISION := $(shell git rev-parse --short HEAD)

# enterprise or open
VARIANT?=open

BINARY_NAME := dcos-signal

LDFLAGS := -X github.com/dcos/dcos-signal/signal.VERSION=$(VERSION) -X github.com/dcos/dcos-signal/signal.REVISION=$(REVISION) -X github.com/dcos/dcos-signal/config.VARIANT=$(VARIANT)

FILES := $(shell go list ./... | grep -v vendor)

# Testing Local Run 
ANON_PATH?=/tmp/anon-id.json
HOST?=localhost
CONFIG?=/tmp/signal-config.json
EXTRA?=

all: test install

test: unit integration

build:
	@echo "+$@"
	go build -v -o signal_'$(VERSION)' -ldflags '$(LDFLAGS)' dcos_signal.go

linux:
	@echo "+$@"
	GOOS=linux go build -v -o signal_'$(VERSION)'_linux -ldflags '$(LDFLAGS)' dcos_signal.go

build-linux:
	@echo "+$@"
	GOOS=linux go build -v -ldflags '$(LDFLAGS)' $(FILES)

install:
	@echo "+$@"
	go install -v -ldflags '$(LDFLAGS)' $(FILES)

integration:
	@cd scripts/mocklicensing && \
		make build && \
		make start && \
		GOCACHE=off go test -v -tags=integration $(FILES)
	@cd scripts/mocklicensing && \
		make build && \
		make stop

unit:
	@go test -v -cover -tags=unit $(FILES)

run:
	@echo "+$@"
	go run dcos_signal.go -v -anonymous-id-path $(ANON_PATH) -report-host $(HOST) -report-port 1050 -c $(CONFIG) $(EXTRA)
