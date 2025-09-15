package main

import (
	"context"
	"log"

	"github.com/Epistemic-Technology/arxiv-mcp/internal/server"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	server := server.CreateServer()
	err := server.Run(context.Background(), &mcp.StdioTransport{})
	if err != nil {
		log.Fatal(err)
	}
}
