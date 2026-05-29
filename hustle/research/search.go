package research

import (
	"fmt"
)

type SearchResult struct {
	URL     string
	Snippet string
	Title   string
}

type SearchInterface interface {
	Query(q string) ([]SearchResult, error)
}

type MockSearch struct{}

func (s *MockSearch) Query(q string) ([]SearchResult, error) {
	fmt.Printf("Mock searching for: %s\n", q)
	return []SearchResult{
		{URL: "https://example.com/result1", Title: "Result 1", Snippet: "This is a mock search result."},
	}, nil
}
