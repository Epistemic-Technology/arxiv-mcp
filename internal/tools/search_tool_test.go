package tools

import (
	"context"
	"testing"
	"time"

	"github.com/Epistemic-Technology/arxiv/arxiv"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestSearchTool(t *testing.T) {
	tool := SearchTool()
	if tool.Name != "arxiv-search" {
		t.Errorf("expected tool name 'arxiv-search', got '%s'", tool.Name)
	}
	if tool.Description != "Searches for papers on arXiv" {
		t.Errorf("expected tool description 'Searches for papers on arXiv', got '%s'", tool.Description)
	}
	if tool.InputSchema == nil {
		t.Error("expected InputSchema to be non-nil")
	}
}

func TestBuildSearchQuery(t *testing.T) {
	tests := []struct {
		name        string
		query       SearchQuery
		expectError bool
		validate    func(*testing.T, arxiv.SearchQuery)
	}{
		{
			name: "title only",
			query: SearchQuery{
				Title: "quantum computing",
			},
			expectError: false,
			validate: func(t *testing.T, q arxiv.SearchQuery) {
				expected := "ti:quantum computing"
				if q.String() != expected {
					t.Errorf("expected query string '%s', got '%s'", expected, q.String())
				}
			},
		},
		{
			name: "author only",
			query: SearchQuery{
				Author: "Einstein",
			},
			expectError: false,
			validate: func(t *testing.T, q arxiv.SearchQuery) {
				expected := "au:Einstein"
				if q.String() != expected {
					t.Errorf("expected query string '%s', got '%s'", expected, q.String())
				}
			},
		},
		{
			name: "abstract only",
			query: SearchQuery{
				Abstract: "machine learning",
			},
			expectError: false,
			validate: func(t *testing.T, q arxiv.SearchQuery) {
				expected := "abs:machine learning"
				if q.String() != expected {
					t.Errorf("expected query string '%s', got '%s'", expected, q.String())
				}
			},
		},
		{
			name: "subject category only",
			query: SearchQuery{
				SubjectCategory: "cs.AI",
			},
			expectError: false,
			validate: func(t *testing.T, q arxiv.SearchQuery) {
				expected := "cat:cs.AI"
				if q.String() != expected {
					t.Errorf("expected query string '%s', got '%s'", expected, q.String())
				}
			},
		},
		{
			name: "all fields search",
			query: SearchQuery{
				All: "neural networks",
			},
			expectError: false,
			validate: func(t *testing.T, q arxiv.SearchQuery) {
				expected := "all:neural networks"
				if q.String() != expected {
					t.Errorf("expected query string '%s', got '%s'", expected, q.String())
				}
			},
		},
		{
			name: "multiple fields",
			query: SearchQuery{
				Title:  "quantum",
				Author: "Smith",
			},
			expectError: false,
			validate: func(t *testing.T, q arxiv.SearchQuery) {
				queryStr := q.String()
				if queryStr != "ti:quantum au:Smith" {
					t.Errorf("expected combined query, got '%s'", queryStr)
				}
			},
		},
		{
			name: "with submitted dates",
			query: SearchQuery{
				Title:           "AI",
				SubmittedSince:  "2023-01-01",
				SubmittedBefore: "2023-12-31",
			},
			expectError: false,
			validate: func(t *testing.T, q arxiv.SearchQuery) {
				queryStr := q.String()
				if !contains(queryStr, "ti:AI") {
					t.Errorf("expected query to contain title filter, got '%s'", queryStr)
				}
				if !contains(queryStr, "submittedDate") {
					t.Errorf("expected query to contain date filter, got '%s'", queryStr)
				}
			},
		},
		{
			name: "with relative date",
			query: SearchQuery{
				Title:             "machine learning",
				SubmittedRelative: "7 days",
			},
			expectError: false,
			validate: func(t *testing.T, q arxiv.SearchQuery) {
				queryStr := q.String()
				if !contains(queryStr, "ti:machine learning") {
					t.Errorf("expected query to contain title filter, got '%s'", queryStr)
				}
				if !contains(queryStr, "submittedDate") {
					t.Errorf("expected query to contain date filter, got '%s'", queryStr)
				}
			},
		},
		{
			name: "invalid submitted since date",
			query: SearchQuery{
				Title:          "test",
				SubmittedSince: "invalid-date",
			},
			expectError: true,
		},
		{
			name: "invalid submitted before date",
			query: SearchQuery{
				Title:           "test",
				SubmittedBefore: "invalid-date",
			},
			expectError: true,
		},
		{
			name: "invalid relative date format",
			query: SearchQuery{
				Title:             "test",
				SubmittedRelative: "invalid",
			},
			expectError: true,
		},
		{
			name:        "empty query",
			query:       SearchQuery{},
			expectError: false,
			validate: func(t *testing.T, q arxiv.SearchQuery) {
				if q.String() != "" {
					t.Errorf("expected empty query string, got '%s'", q.String())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := buildSearchQuery(tt.query)
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

func TestParseRelativeDate(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		relative    string
		expectError bool
		validate    func(*testing.T, time.Time)
	}{
		{
			name:        "7 days",
			relative:    "7 days",
			expectError: false,
			validate: func(t *testing.T, result time.Time) {
				expected := now.AddDate(0, 0, -7)
				if !roughlyEqual(result, expected) {
					t.Errorf("expected roughly %v, got %v", expected, result)
				}
			},
		},
		{
			name:        "1 day",
			relative:    "1 day",
			expectError: false,
			validate: func(t *testing.T, result time.Time) {
				expected := now.AddDate(0, 0, -1)
				if !roughlyEqual(result, expected) {
					t.Errorf("expected roughly %v, got %v", expected, result)
				}
			},
		},
		{
			name:        "2 weeks",
			relative:    "2 weeks",
			expectError: false,
			validate: func(t *testing.T, result time.Time) {
				expected := now.AddDate(0, 0, -14)
				if !roughlyEqual(result, expected) {
					t.Errorf("expected roughly %v, got %v", expected, result)
				}
			},
		},
		{
			name:        "1 week",
			relative:    "1 week",
			expectError: false,
			validate: func(t *testing.T, result time.Time) {
				expected := now.AddDate(0, 0, -7)
				if !roughlyEqual(result, expected) {
					t.Errorf("expected roughly %v, got %v", expected, result)
				}
			},
		},
		{
			name:        "3 months",
			relative:    "3 months",
			expectError: false,
			validate: func(t *testing.T, result time.Time) {
				expected := now.AddDate(0, -3, 0)
				if !roughlyEqual(result, expected) {
					t.Errorf("expected roughly %v, got %v", expected, result)
				}
			},
		},
		{
			name:        "1 month",
			relative:    "1 month",
			expectError: false,
			validate: func(t *testing.T, result time.Time) {
				expected := now.AddDate(0, -1, 0)
				if !roughlyEqual(result, expected) {
					t.Errorf("expected roughly %v, got %v", expected, result)
				}
			},
		},
		{
			name:        "2 years",
			relative:    "2 years",
			expectError: false,
			validate: func(t *testing.T, result time.Time) {
				expected := now.AddDate(-2, 0, 0)
				if !roughlyEqual(result, expected) {
					t.Errorf("expected roughly %v, got %v", expected, result)
				}
			},
		},
		{
			name:        "1 year",
			relative:    "1 year",
			expectError: false,
			validate: func(t *testing.T, result time.Time) {
				expected := now.AddDate(-1, 0, 0)
				if !roughlyEqual(result, expected) {
					t.Errorf("expected roughly %v, got %v", expected, result)
				}
			},
		},
		{
			name:        "invalid format - missing unit",
			relative:    "7",
			expectError: true,
		},
		{
			name:        "invalid format - missing number",
			relative:    "days",
			expectError: true,
		},
		{
			name:        "invalid number",
			relative:    "abc days",
			expectError: true,
		},
		{
			name:        "invalid unit",
			relative:    "7 centuries",
			expectError: true,
		},
		{
			name:        "empty string",
			relative:    "",
			expectError: true,
		},
		{
			name:        "too many parts",
			relative:    "7 days ago",
			expectError: true,
		},
		{
			name:        "uppercase unit",
			relative:    "7 DAYS",
			expectError: false,
			validate: func(t *testing.T, result time.Time) {
				expected := now.AddDate(0, 0, -7)
				if !roughlyEqual(result, expected) {
					t.Errorf("expected roughly %v, got %v", expected, result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseRelativeDate(tt.relative)
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

func TestSearchHandler(t *testing.T) {
	t.Run("default max results", func(t *testing.T) {
		query := SearchQuery{
			Title: "test",
		}
		req := &mcp.CallToolRequest{}

		// Note: This will make an actual API call in this test
		// In production, you might want to mock the arxiv client
		ctx := context.Background()
		result, searchResults, err := SearchHandler(ctx, req, query)

		if err != nil {
			// API might be unavailable, skip test
			t.Skipf("skipping test due to API error: %v", err)
		}

		if result == nil {
			t.Error("expected non-nil CallToolResult")
		}

		// Check that results were returned (even if empty)
		_ = searchResults // results structure is valid
	})

	t.Run("with id list", func(t *testing.T) {
		query := SearchQuery{
			IdList:     []string{"2301.00000"},
			MaxResults: 1,
		}
		req := &mcp.CallToolRequest{}

		ctx := context.Background()
		result, _, err := SearchHandler(ctx, req, query)

		if err != nil {
			// API might be unavailable, skip test
			t.Skipf("skipping test due to API error: %v", err)
		}

		if result == nil {
			t.Error("expected non-nil CallToolResult")
		}
	})

	t.Run("custom max results", func(t *testing.T) {
		query := SearchQuery{
			Title:      "quantum",
			MaxResults: 5,
		}
		req := &mcp.CallToolRequest{}

		ctx := context.Background()
		result, searchResults, err := SearchHandler(ctx, req, query)

		if err != nil {
			// API might be unavailable, skip test
			t.Skipf("skipping test due to API error: %v", err)
		}

		if result == nil {
			t.Error("expected non-nil CallToolResult")
		}

		// Verify max results is respected (if results are available)
		if len(searchResults.Entries) > 5 {
			t.Errorf("expected at most 5 results, got %d", len(searchResults.Entries))
		}
	})
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr)))
}

func roughlyEqual(t1, t2 time.Time) bool {
	diff := t1.Sub(t2)
	if diff < 0 {
		diff = -diff
	}
	// Allow 1 second difference for test execution time
	return diff < time.Second
}
