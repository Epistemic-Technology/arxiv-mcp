package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Epistemic-Technology/arxiv-mcp/internal/server"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}
	httpHandler := mcp.NewStreamableHTTPHandler(getServerForRequest, nil)
	if err := http.ListenAndServe(":"+port, httpHandler); err != nil {
		log.Fatal(err)
	}
}

func getServerForRequest(r *http.Request) *mcp.Server {
	return server.CreateServer()
}
