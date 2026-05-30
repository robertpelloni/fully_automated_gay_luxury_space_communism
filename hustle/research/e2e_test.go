package research

import (
	"github.com/robertpelloni/hustle/orchestrator"
	"testing"
)

func TestResearchPipelineE2E(t *testing.T) {
	orch := orchestrator.NewOrchestrator()

	// 1. Multi-Provider Search with Orchestrator Integration
	searcher := &ResearchSearch{
		ActiveProvider: Tavily,
		Orchestrator:   orch,
	}
	results, err := searcher.Query("Advanced automated trading strategies")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("Expected at least one search result")
	}

	// Verify integration with Orchestrator Memory
	if len(orch.L1.Entries) == 0 {
		t.Fatal("Search results were not persisted in Orchestrator memory")
	}

	// 2. Enhanced Report Synthesis
	report := &Report{Title: "Trading Alpha Report"}
	report.Synthesize(results)
	if report.Content == "" {
		t.Fatal("Report synthesis produced empty content")
	}

	// 3. Export PDF (Mock)
	err = report.ExportPDF("trading_alpha.pdf")
	if err != nil {
		t.Fatalf("PDF export failed: %v", err)
	}
}
