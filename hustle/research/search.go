package research

import (
	"fmt"
	"github.com/robertpelloni/hustle/orchestrator"
	"time"
)

type SearchResult struct {
	URL      string
	Snippet  string
	Title    string
	Provider string
}

type SearchInterface interface {
	Query(q string) ([]SearchResult, error)
}

type Provider string

const (
	Tavily Provider = "Tavily"
	Brave  Provider = "Brave"
	Google Provider = "Google"
)

type ResearchSearch struct {
	ActiveProvider Provider
	Orchestrator   *orchestrator.Orchestrator
}

func (s *ResearchSearch) Query(q string) ([]SearchResult, error) {
	fmt.Printf("Searching via %s for: %s\n", s.ActiveProvider, q)

	results := []SearchResult{
		{URL: "https://hustle.com/info", Title: "Hustle Strategy", Snippet: "Key insights for automated revenue.", Provider: string(s.ActiveProvider)},
	}

	// Integration: Store in Orchestrator Memory if available
	if s.Orchestrator != nil {
		for _, res := range results {
			entry := orchestrator.MemoryEntry{
				ID:        res.URL,
				Content:   fmt.Sprintf("%s: %s", res.Title, res.Snippet),
				BaseScore: 50.0, // Initial relevance score
				Timestamp: time.Now(),
				Tags:      []string{"research", string(s.ActiveProvider)},
			}
			s.Orchestrator.L1.Add(entry)
		}
	}

	return results, nil
}

type MockSearch struct {
	Orchestrator *orchestrator.Orchestrator
}

func (s *MockSearch) Query(q string) ([]SearchResult, error) {
	results := []SearchResult{
		{URL: "https://example.com/mock", Title: "Mock Result", Snippet: "Mock data.", Provider: "Mock"},
	}

	if s.Orchestrator != nil {
		for _, res := range results {
			entry := orchestrator.MemoryEntry{
				ID:        res.URL,
				Content:   res.Snippet,
				Timestamp: time.Now(),
			}
			s.Orchestrator.L1.Add(entry)
		}
	}

	return results, nil
}
