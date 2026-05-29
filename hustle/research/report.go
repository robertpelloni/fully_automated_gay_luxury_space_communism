package research

import (
	"fmt"
)

type Report struct {
	Title   string
	Content string
}

func (r *Report) Synthesize(results []SearchResult) {
	r.Content = "Synthesized report based on search results:\n"
	for _, res := range results {
		r.Content += fmt.Sprintf("- %s: %s\n", res.Title, res.Snippet)
	}
}

func (r *Report) ExportPDF(filepath string) error {
	fmt.Printf("Exporting report to PDF: %s\n", filepath)
	return nil
}
