package server

import (
	"github.com/Epistemic-Technology/arxiv-mcp/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func CreateServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "arxiv-mcp", Version: "v0.0.1"}, nil)
	mcp.AddTool(server, tools.SearchTool(), tools.SearchHandler)

	return server
}
