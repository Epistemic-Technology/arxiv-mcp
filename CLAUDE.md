# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Model Context Protocol (MCP) server for arXiv, written in Go. It provides tools for searching academic papers on arXiv through the MCP interface.

## Architecture

The codebase follows a standard Go project structure:

- **cmd/arxiv-mcp-local-server/**: Entry point for the MCP server that runs via stdio transport
- **internal/server/**: Server initialization and configuration logic that creates the MCP server instance
- **internal/tools/**: MCP tool implementations
  - `search-tool.go`: Implements the `arxiv-search` tool that queries arXiv API with support for filtering by title, author, abstract, category, and date ranges

The server uses:
- `github.com/modelcontextprotocol/go-sdk` for MCP functionality
- `github.com/Epistemic-Technology/arxiv` as the arXiv API client

## Common Commands

### Build
```bash
go build ./cmd/arxiv-mcp-local-server
```

### Run tests
```bash
go test ./...
```

### Run the server
```bash
go run ./cmd/arxiv-mcp-local-server/main.go
```

### Format code
```bash
go fmt ./...
```

### Check for issues
```bash
go vet ./...
```

### Update dependencies
```bash
go mod tidy
```

## Key Implementation Details

The MCP server exposes an `arxiv-search` tool that accepts SearchQuery parameters including:
- Title, Author, Abstract, SubjectCategory filters
- Date range filtering (SubmittedSince/SubmittedBefore)
- All-fields search
- MaxResults limit (defaults to 20)

The search tool constructs arXiv API queries and returns structured search results through the MCP protocol.