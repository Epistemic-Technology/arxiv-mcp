package tools

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Epistemic-Technology/arxiv/arxiv"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SearchQuery struct {
	Title             string   `json:"title,omitempty"`
	Author            string   `json:"author,omitempty"`
	Abstract          string   `json:"abstract,omitempty"`
	SubjectCategory   string   `json:"subject_category,omitempty" jsonschema:"subject category, using arXiv category taxonomy"`
	SubmittedSince    string   `json:"submitted_since,omitempty" pattern:"\\d{4}-\\d{2}-\\d{2}" jsonschema:"date in YYYY-MM-DD"`
	SubmittedBefore   string   `json:"submitted_before,omitempty" pattern:"\\d{4}-\\d{2}-\\d{2}" jsonschema:"date in YYYY-MM-DD"`
	SubmittedRelative string   `json:"submitted_relative,omitempty" pattern:"[0-9]+ (days|weeks|months|years)" jsonschema:"relative date in days, weeks, months, or years from today"`
	All               string   `json:"all,omitempty" jsonschema:"search within title, author, abstract, subject"`
	IdList            []string `json:"id_list,omitempty" jsonschema:"array of arXiv IDs to search within. Can be passed alone to retrieve specific papers"`
	MaxResults        int      `json:"max,omitempty"`
	ReturnFields      []string `json:"return_fields,omitempty" jsonschema:"array of fields to return. Returns all if empty"`
}

type SearchResults struct {
	Entries []EntryView `json:"entries,omitempty"`
}

type EntryView struct {
	ID               *string           `json:"id,omitempty"`
	Title            *string           `json:"title,omitempty"`
	Published        *time.Time        `json:"published,omitempty"`
	Updated          *time.Time        `json:"updated,omitempty"`
	Summary          *string           `json:"summary,omitempty"`
	Authors          *[]arxiv.Author   `json:"authors,omitempty"`
	Categories       *[]arxiv.Category `json:"categories,omitempty"`
	PrimaryCategory  *arxiv.Category   `json:"primaryCategory,omitempty"`
	Links            *[]arxiv.Link     `json:"links,omitempty"`
	Comment          *string           `json:"comment,omitempty"`
	JournalReference *string           `json:"journalReference,omitempty"`
	DOI              *string           `json:"doi,omitempty"`
	AbstractUrl      *string           `json:"abstractUrl,omitempty"`
	PDFUrl           *string           `json:"pdfUrl,omitempty"`
}

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
	if len(query.IdList) > 0 {
		params.IdList = query.IdList
	}
	arxivClient := arxiv.NewClient()
	results, err := arxivClient.Search(ctx, params)
	if err != nil {
		return nil, SearchResults{}, err
	}

	// Filter to only requested fields
	filteredEntries := make([]EntryView, len(results.Entries))
	for i, entry := range results.Entries {
		filteredEntries[i] = filterEntry(entry, query.ReturnFields)
	}
	searchResults := SearchResults{
		Entries: filteredEntries,
	}

	return &mcp.CallToolResult{}, searchResults, nil
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

	// Handle relative date if provided and no explicit dates are set
	if query.SubmittedRelative != "" && query.SubmittedSince == "" && query.SubmittedBefore == "" {
		since, err := parseRelativeDate(query.SubmittedRelative)
		if err != nil {
			return *arxivQuery, err
		}
		before := time.Now()
		arxivQuery = arxivQuery.SubmittedBetween(since, before)
	} else if query.SubmittedSince != "" || query.SubmittedBefore != "" {
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

func filterEntry(entry arxiv.EntryMetadata, fields []string) EntryView {
	view := EntryView{}

	if len(fields) == 0 {
		view = EntryView{
			ID:               &entry.ID,
			Title:            &entry.Title,
			Published:        &entry.Published,
			Updated:          &entry.Updated,
			Summary:          &entry.Summary,
			Authors:          &entry.Authors,
			Categories:       &entry.Categories,
			PrimaryCategory:  &entry.PrimaryCategory,
			Links:            &entry.Links,
			Comment:          &entry.Comment,
			JournalReference: &entry.JournalReference,
			DOI:              &entry.DOI,
		}
		return view
	}

	for _, field := range fields {
		switch strings.ToLower(field) {
		case "id":
			view.ID = &entry.ID
		case "title":
			view.Title = &entry.Title
		case "published":
			view.Published = &entry.Published
		case "updated":
			view.Updated = &entry.Updated
		case "summary", "abstract":
			view.Summary = &entry.Summary
		case "authors", "author":
			authors := entry.Authors
			view.Authors = &authors
		case "categories", "category":
			categories := entry.Categories
			view.Categories = &categories
		case "primarycategory", "primary_category":
			view.PrimaryCategory = &entry.PrimaryCategory
		case "links", "link":
			links := entry.Links
			view.Links = &links
		case "comment":
			view.Comment = &entry.Comment
		case "journalreference", "journal_reference", "journal":
			view.JournalReference = &entry.JournalReference
		case "doi":
			view.DOI = &entry.DOI
		case "abstracturl", "abstract_url":
			view.AbstractUrl = &entry.AbstractUrl
		case "pdfurl", "pdf_url", "pdf":
			view.PDFUrl = &entry.PDFUrl
		}
	}

	return view
}

func parseRelativeDate(relative string) (time.Time, error) {
	parts := strings.Fields(relative)
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid relative date format: %s", relative)
	}

	num, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid number in relative date: %s", parts[0])
	}

	unit := strings.ToLower(parts[1])
	now := time.Now()

	switch unit {
	case "day", "days":
		return now.AddDate(0, 0, -num), nil
	case "week", "weeks":
		return now.AddDate(0, 0, -num*7), nil
	case "month", "months":
		return now.AddDate(0, -num, 0), nil
	case "year", "years":
		return now.AddDate(-num, 0, 0), nil
	default:
		return time.Time{}, fmt.Errorf("invalid time unit in relative date: %s", unit)
	}
}
