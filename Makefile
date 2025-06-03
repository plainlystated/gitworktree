APP_NAME := gwt
VERSION := $(shell git describe --tags --always --dirty)
BUILD_DIR := ./bin
BIN := $(BUILD_DIR)/$(APP_NAME)
DOCKER_IMAGE := gwt-builder
GOFILES := $(shell find . -name '*.go' -type f)

.PHONY: all build clean install docker-build

all: build

build: $(BIN)

$(BIN): $(GOFILES)
	mkdir -p $(BUILD_DIR)
	go build -ldflags "-X main.version=$(VERSION)" -o $(BIN) ./cmd/tui

install: build
	install -m 755 $(BIN) /usr/local/bin/$(APP_NAME)

clean:
	rm -rf $(BUILD_DIR)

docker-build:
	docker build -f Dockerfile.build -t $(DOCKER_IMAGE) .
	docker run --rm -v $(PWD):/app -w /app $(DOCKER_IMAGE) make build

