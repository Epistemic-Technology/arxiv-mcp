package server

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/Epistemic-Technology/arxiv-mcp/internal/resources"
	"github.com/Epistemic-Technology/arxiv-mcp/internal/tools"
)

func CreateServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "arxiv-mcp", Version: "v0.0.1"}, nil)
	mcp.AddTool(server, tools.SearchTool(), tools.SearchHandler)
	server.AddResource(&resources.TaxonomyResource, resources.TaxonomyResourceHandler)
	return server
}
