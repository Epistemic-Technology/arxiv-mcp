package resources

import (
	"context"
	_ "embed"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var TaxonomyResource = mcp.Resource{
	Name:        "category-taxonomy",
	Description: "A JSON representation of the arXiv category taxonomy, showing all category tags and their descriptions.",
	Title:       "Category Taxonomy",
	URI:         "file://arxiv/taxonomy.json",
}

//go:embed arxiv-taxonomy.json
var taxonomyData string

func TaxonomyResourceHandler(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     taxonomyData,
			},
		},
	}, nil
}
