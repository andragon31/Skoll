.PHONY: build test lint clean install run

BINARY_NAME=skoll
VERSION=$(shell git describe --tags --always 2>/dev/null || echo "v0.1.0")
GO_LDFLAGS=-ldflags "-X main.version=${VERSION}"

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod

BUILD_DIR=./bin

all: clean test build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(GO_LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/skoll

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

install: build
	@echo "Installing $(BINARY_NAME)..."
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)

init:
	@echo "Initializing Skoll..."
	./$(BUILD_DIR)/$(BINARY_NAME) init

mcp:
	./$(BUILD_DIR)/$(BINARY_NAME) mcp

tui:
	./$(BUILD_DIR)/$(BINARY_NAME) tui
