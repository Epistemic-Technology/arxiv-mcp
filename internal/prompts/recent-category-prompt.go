package prompts

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var CategoryPrompt = mcp.Prompt{
	Name:        "recent-category",
	Description: "Get articles from the last week for a specific category",
	Arguments: []*mcp.PromptArgument{
		{
			Name:        "category",
			Description: "The category to get articles from",
			Required:    true,
		},
	},
}

func CategoryPromptHandler(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "Prompt to get articles from the last week for a specific category",
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: "Find the arXiv category for " + req.Params.Arguments["category"] + ". If the category matches a general subject like math or computer science, get the category for general articles within that field. Search for 50 articles from the last week in that category. If none are found, try expanding the time range to the last month, 6 months, or a year. Display them in a table with columns for title, first author, ID, and PDF URL."},
			},
		},
	}, nil
}
