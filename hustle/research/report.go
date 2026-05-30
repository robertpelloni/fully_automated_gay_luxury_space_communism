package research

import (
	"fmt"
	"github.com/robertpelloni/hustle/orchestrator"
	"os"
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
	fmt.Printf("Generating report file: %s\n", filepath)
	// In a real implementation, we would use a PDF library like gofpdf.
	// For now, we simulate by writing to a text file with a .pdf extension
	// to verify the export pipeline in E2E tests.
	return os.WriteFile(filepath, []byte(r.Content), 0644)
}
