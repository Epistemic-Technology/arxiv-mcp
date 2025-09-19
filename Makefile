# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

# Binary directory
BINARY_DIR=bin

# Binary names
BINARIES=arxiv-mcp-local-server arxiv-mcp-http-server arxiv-taxonomy-scraper

# Build flags
LDFLAGS=-ldflags "-s -w"

export ARXIV_MCP_ROOT := $(CURDIR)

# Default target
.PHONY: all
all: build

# Build all binaries
.PHONY: build
build: $(BINARIES)

# Individual binary targets
.PHONY: arxiv-mcp-local-server
arxiv-mcp-local-server:
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/arxiv-mcp-local-server ./cmd/arxiv-mcp-local-server

.PHONY: arxiv-mcp-http-server
arxiv-mcp-http-server:
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/arxiv-mcp-http-server ./cmd/arxiv-mcp-http-server

.PHONY: arxiv-taxonomy-scraper
arxiv-taxonomy-scraper:
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/arxiv-taxonomy-scraper ./cmd/arxiv-taxonomy-scraper

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)

# Run the server (development)
.PHONY: run
run:
	$(GOCMD) run ./cmd/arxiv-mcp-local-server/main.go

# Install binaries to GOPATH/bin
.PHONY: install
install:
	$(GOCMD) install ./cmd/...

# Run the MCP inspector on local server
.PHONY: inspect
inspect:
	npx @modelcontextprotocol/inspector $(ARXIV_MCP_ROOT)/$(BINARY_DIR)/arxiv-mcp-local-server

# Add local server to claude code
.PHONY: cc-add-mcp
cc-add-mcp:
	-claude mcp remove arxiv-mcp-http-server --scope project
	claude mcp add arxiv-mcp-local-server --scope project -- $(ARXIV_MCP_ROOT)/$(BINARY_DIR)/arxiv-mcp-local-server

# Run http server and add to claude code
.PHONY: cc-add-http
cc-add-http:
	-claude mcp remove arxiv-mcp-local-server --scope project
	-claude mcp add -t http --scope=project arxiv-mcp-http-server http://localhost:8888
	$(ARXIV_MCP_ROOT)/$(BINARY_DIR)/arxiv-mcp-http-server

# Remove all mcp servers from claude code
.PHONY: cc-remove-mcp
cc-remove-mcp:
	-claude mcp remove arxiv-mcp-local-server --scope project
	-claude mcp remove arxiv-mcp-http-server --scope project

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all                   - Build all binaries (default)"
	@echo "  build                 - Build all binaries"
	@echo "  arxiv-mcp-local-server - Build arxiv-mcp-local-server binary"
	@echo "  arxiv-mcp-http-server  - Build arxiv-mcp-http-server binary"
	@echo "  arxiv-taxonomy-scraper - Build arxiv-taxonomy-scraper binary"
	@echo "  test                  - Run tests"
	@echo "  clean                 - Remove build artifacts"
	@echo "  run                   - Run the server in development mode"
	@echo "  install               - Install binaries to GOPATH/bin"
	@echo "  inspect               - Run the MCP inspector on local server"
	@echo "  cc-add-mcp            - Add local server to claude code"
	@echo "  help                  - Show this help message"
