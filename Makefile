# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
GOMOD=$(GOCMD) mod

# Binary directory
BINARY_DIR=bin

# Binary names
BINARIES=arxiv-mcp-local-server arxiv-taxonomy-scraper

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

.PHONY: arxiv-taxonomy-scraper
arxiv-taxonomy-scraper:
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/arxiv-taxonomy-scraper ./cmd/arxiv-taxonomy-scraper

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Format code
.PHONY: fmt
fmt:
	$(GOFMT) ./...

# Vet code for issues
.PHONY: vet
vet:
	$(GOVET) ./...

# Run linter (requires golangci-lint)
.PHONY: lint
lint:
	@which golangci-lint > /dev/null || echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
	@which golangci-lint > /dev/null && golangci-lint run ./... || true

# Tidy dependencies
.PHONY: tidy
tidy:
	$(GOMOD) tidy

# Download dependencies
.PHONY: deps
deps:
	$(GOMOD) download

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)
	rm -f coverage.out coverage.html

# Run the server (development)
.PHONY: run
run:
	$(GOCMD) run ./cmd/arxiv-mcp-local-server/main.go

# Install binaries to GOPATH/bin
.PHONY: install
install:
	$(GOCMD) install ./cmd/...

# Check code quality (fmt, vet, test)
.PHONY: check
check: fmt vet test

# Build for multiple platforms
.PHONY: build-all-platforms
build-all-platforms:
	@mkdir -p $(BINARY_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/arxiv-mcp-local-server-darwin-amd64 ./cmd/arxiv-mcp-local-server
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/arxiv-mcp-local-server-darwin-arm64 ./cmd/arxiv-mcp-local-server
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/arxiv-mcp-local-server-linux-amd64 ./cmd/arxiv-mcp-local-server
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/arxiv-mcp-local-server-linux-arm64 ./cmd/arxiv-mcp-local-server
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/arxiv-mcp-local-server-windows-amd64.exe ./cmd/arxiv-mcp-local-server

# Run the MCP inspector on local server
.PHONY: inspect
inspect:
	npx @modelcontextprotocol/inspector $(ARXIV_MCP_ROOT)/$(BINARY_DIR)/arxiv-mcp-local-server

# Add local server to claude code
.PHONY: cc-add-mcp
cc-add-mcp:
	claude mcp add arxiv-mcp-local-server --scope project -- $(ARXIV_MCP_ROOT)/$(BINARY_DIR)/arxiv-mcp-local-server

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all                   - Build all binaries (default)"
	@echo "  build                 - Build all binaries"
	@echo "  arxiv-mcp-local-server - Build arxiv-mcp-local-server binary"
	@echo "  arxiv-taxonomy-scraper - Build arxiv-taxonomy-scraper binary"
	@echo "  test                  - Run tests"
	@echo "  test-coverage         - Run tests with coverage report"
	@echo "  fmt                   - Format code"
	@echo "  vet                   - Vet code for issues"
	@echo "  lint                  - Run golangci-lint"
	@echo "  tidy                  - Tidy module dependencies"
	@echo "  deps                  - Download dependencies"
	@echo "  clean                 - Remove build artifacts"
	@echo "  run                   - Run the server in development mode"
	@echo "  install               - Install binaries to GOPATH/bin"
	@echo "  check                 - Run fmt, vet, and test"
	@echo "  build-all-platforms   - Build for multiple platforms"
	@echo "  inspect               - Run the MCP inspector on local server"
	@echo "  cc-add-mcp            - Add local server to claude code"
	@echo "  help                  - Show this help message"
