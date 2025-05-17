# Makefile for kv project

# Go build command
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean

# Output directory for binaries
BIN_DIR := ./bin

# Binaries to build
BINARIES := kvController kvNode kvClient kvLoadBalancer

.PHONY: all
all: $(BINARIES)

.PHONY: kvController
kvController:
	$(GOBUILD) -o $(BIN_DIR)/kvController ./cmd/kvController

.PHONY: kvNode
kvNode:
	$(GOBUILD) -o $(BIN_DIR)/kvNode ./cmd/kvNode

.PHONY: kvClient
kvClient:
	$(GOBUILD) -o $(BIN_DIR)/kvClient ./cmd/kvClient

.PHONY: kvLoadBalancer
kvLoadBalancer:
	$(GOBUILD) -o $(BIN_DIR)/kvLoadBalancer ./cmd/kvLoadBalancer

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(BIN_DIR)
