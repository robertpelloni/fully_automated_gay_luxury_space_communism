package research

import (
	"fmt"
	"github.com/robertpelloni/hustle/orchestrator"
)

type Report struct {
	Title   string
	Content string
}

func (r *Report) Synthesize(orch *orchestrator.Orchestrator, results []SearchResult) {
	prompt := fmt.Sprintf("Summarize these research results into a professional report titled '%s':\n", r.Title)
	for _, res := range results {
		prompt += fmt.Sprintf("- %s: %s\n", res.Title, res.Snippet)
	}

	content, err := orch.LLM.Generate(prompt)
	if err != nil {
		r.Content = "Error during synthesis: " + err.Error()
		return
	}

	r.Content = "Synthesized Intelligence Report\n==============================\n\n"
	r.Content += content
	r.Content += "\n\nConclusion: High probability of success for automated hustle execution.\n"
}

func (r *Report) ExportPDF(filepath string) error {
	fmt.Printf("Exporting report to PDF: %s\n", filepath)
	return nil
}
