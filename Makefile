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
	GO111MODULE=on go build -mod=vendor -v -o signal_'$(VERSION)' -ldflags '$(LDFLAGS)' dcos_signal.go

linux:
	@echo "+$@"
	GO111MODULE=on GOOS=linux go build -mod=vendor -v -o signal_'$(VERSION)'_linux -ldflags '$(LDFLAGS)' dcos_signal.go

build-linux:
	@echo "+$@"
	GO111MODULE=on GOOS=linux go build -mod=vendor -v -ldflags '$(LDFLAGS)' $(FILES)

install:
	@echo "+$@"
	GO111MODULE=on go install -mod=vendor -v -ldflags '$(LDFLAGS)' $(FILES)

integration:
	@-cd scripts/mocklicensing && \
		make build && \
		make start && \
		go test -v -count=1 -tags=integration $(FILES)
	@cd scripts/mocklicensing && \
		make stop && make clean

unit:
	@GO111MODULE=on go test -v -cover -mod=vendor -tags=unit $(FILES)

run:
	@echo "+$@"
	GO111MODULE=on go run -mod=vendor  dcos_signal.go -v -anonymous-id-path $(ANON_PATH) -report-host $(HOST) -report-port 1050 -c $(CONFIG) $(EXTRA)
