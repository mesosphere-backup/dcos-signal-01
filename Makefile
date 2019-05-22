DEFAULT_TARGET: build

VERSION := 0.4.0
COMMIT := $(shell git rev-parse --short HEAD)
LDFLAGS := -X github.com/dcos/dcos-signal/config.Version=$(VERSION) -X github.com/dcos/dcos-signal/config.Commit=$(COMMIT)

CURRENT_DIR=$(shell pwd)
BUILD_DIR=build
PKG_DIR=/dcos-signal
BINARY_NAME=dcos-signal
IMAGE_NAME=dcos-signal-dev
SRCS := $(shell find . -type f -name '*.go' -not -path './vendor/*')

.PHONY: docker
docker:
ifndef NO_DOCKER
	docker build -t $(IMAGE_NAME) .
endif

$(BUILD_DIR)/$(BINARY_NAME): $(SRCS)
	mkdir -p $(BUILD_DIR)
	$(call inDocker,go build -mod=vendor -v -ldflags '$(LDFLAGS)' -o $(BUILD_DIR)/$(BINARY_NAME))

.PHONY: build
build: docker $(BUILD_DIR)/$(BINARY_NAME)

# install does not run in a docker container to build for the correct OS
.PHONY: install
install:
	go install -mod=vendor -v -ldflags '$(LDFLAGS)'

.PHONY: test
test: docker
	$(call inDocker,bash -x -c './scripts/test.sh')

.PHONY: integration
integration:
	$(call inDocker,bash -x -c './scripts/integration.sh')

.PHONY: publish
publish:
	@echo "This project doesn't have any artifacts to be published"

.PHONY: clean
clean:
	rm -rf ./$(BUILD_DIR)

ifdef NO_DOCKER
  define inDocker
    $1
  endef
else
  define inDocker
    docker run \
      $(shell test -t 0 && echo "-t -i" || true) \
      -v $(CURRENT_DIR):$(PKG_DIR) \
      -w $(PKG_DIR) \
      --rm \
      $(IMAGE_NAME) \
      $1
  endef
endif
