package tools

import (
	"context"
	"time"

	"github.com/Epistemic-Technology/arxiv/arxiv"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SearchQuery struct {
	Title           string `json:"title,omitempty"`
	Author          string `json:"author,omitempty"`
	Abstract        string `json:"abstract,omitempty"`
	SubjectCategory string `json:"subject_category,omitempty"`
	SubmittedSince  string `json:"submitted_since,omitempty" pattern:"\\d{4}-\\d{2}-\\d{2}" jsonschema:"date in YYYY-MM-DD"`
	SubmittedBefore string `json:"submitted_before,omitempty" pattern:"\\d{4}-\\d{2}-\\d{2}" jsonschema:"date in YYYY-MM-DD"`
	All             string `json:"all,omitempty"`
	MaxResults      int    `json:"max,omitempty"`
}

type SearchResults arxiv.SearchResults

func SearchTool() *mcp.Tool {
	inputSchema, err := jsonschema.For[SearchQuery](nil)
	if err != nil {
		panic(err)
	}

	searchTool := mcp.Tool{
		Name:        "arxiv-search",
		Description: "Searches for papers on arXiv",
		InputSchema: inputSchema,
	}
	return &searchTool
}

func SearchHandler(ctx context.Context, req *mcp.CallToolRequest, query SearchQuery) (*mcp.CallToolResult, SearchResults, error) {
	arxivQuery, err := buildSearchQuery(query)
	if err != nil {
		return nil, SearchResults{}, err
	}
	max := query.MaxResults
	if max == 0 {
		max = 20
	}
	params := arxiv.SearchParams{
		Query:      arxivQuery.String(),
		MaxResults: max,
		SortBy:     arxiv.SortByRelevance,
		SortOrder:  arxiv.SortOrderDescending,
	}
	arxivClient := arxiv.NewClient()
	results, err := arxivClient.Search(ctx, params)
	if err != nil {
		return nil, SearchResults{}, err
	}
	return &mcp.CallToolResult{}, SearchResults(results), nil
}

func buildSearchQuery(query SearchQuery) (arxiv.SearchQuery, error) {
	arxivQuery := arxiv.NewSearchQuery()
	if query.Title != "" {
		arxivQuery = arxivQuery.Title(query.Title)
	}

	if query.Author != "" {
		arxivQuery = arxivQuery.Author(query.Author)
	}

	if query.Abstract != "" {
		arxivQuery = arxivQuery.Abstract(query.Abstract)
	}

	if query.SubjectCategory != "" {
		arxivQuery = arxivQuery.Category(query.SubjectCategory)
	}

	if query.All != "" {
		arxivQuery = arxivQuery.All(query.All)
	}

	if query.SubmittedSince != "" || query.SubmittedBefore != "" {
		var since, before time.Time
		if query.SubmittedBefore != "" {
			var err error
			before, err = time.Parse("2006-01-02", query.SubmittedBefore)
			if err != nil {
				return *arxivQuery, err
			}
		} else {
			before = time.Now()
		}
		if query.SubmittedSince != "" {
			var err error
			since, err = time.Parse("2006-01-02", query.SubmittedSince)
			if err != nil {
				return *arxivQuery, err
			}
		} else {
			since = time.Time{}
		}
		arxivQuery = arxivQuery.SubmittedBetween(since, before)
	}

	return *arxivQuery, nil
}
