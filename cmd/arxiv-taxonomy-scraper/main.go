package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Field struct {
	Title      string     `json:"title"`
	Categories []Category `json:"categories"`
}

type Category struct {
	Tag         string `json:"tag"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

func main() {
	outputFile := "arxiv-taxonomy.json"
	if len(os.Args) > 1 {
		outputFile = os.Args[1]
	}

	log.Printf("Fetching arXiv category taxonomy from https://arxiv.org/category_taxonomy")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get("https://arxiv.org/category_taxonomy")
	if err != nil {
		log.Fatalf("Error fetching taxonomy: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("HTTP error: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	fieldHeadingSelections := doc.Find("h2")
	fieldHeadings := make([]string, 0)
	for i, node := range fieldHeadingSelections.Nodes {
		if i < 3 {
			continue
		}
		fieldHeadings = append(fieldHeadings, node.FirstChild.Data)
	}

	fields := make([]Field, 0)
	fieldSelections := doc.Find(".accordion-body").EachIter()
	for i, fieldSelection := range fieldSelections {
		field := Field{
			Title:      fieldHeadings[i],
			Categories: make([]Category, 0),
		}

		categorySelections := fieldSelection.Find(".columns.divided").EachIter()
		for _, categorySelection := range categorySelections {
			tagLine := categorySelection.Find("h4").First().Text()
			re := regexp.MustCompile(`^([^\s(]+)\s+\(([^)]+)\)`)
			matches := re.FindStringSubmatch(tagLine)
			var tag string
			if len(matches) >= 2 {
				tag = matches[1]
			}
			var label string
			if len(matches) == 3 {
				label = matches[2]
			}
			description := categorySelection.Find(".column:not(.is-one-fifth)").Find("p").First().Text()
			category := Category{
				Tag:         tag,
				Label:       label,
				Description: description,
			}
			field.Categories = append(field.Categories, category)
		}

		fields = append(fields, field)
	}

	jsonData, err := json.Marshal(fields)
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
	}

	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	log.Printf("Successfully wrote taxonomy to %s", outputFile)
}
