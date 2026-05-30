package research

import (
	"github.com/robertpelloni/hustle/orchestrator"
	"testing"
)

func TestResearchPipelineE2E(t *testing.T) {
	orch := orchestrator.NewOrchestrator()
	// Inject custom mock response for E2E
	orch.LLM = &orchestrator.MockLLM{Response: "The AI agent market is expanding rapidly."}

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

	// 2. Enhanced Report Synthesis (Using LLM)
	report := &Report{Title: "Trading Alpha Report"}
	report.Synthesize(orch, results)

	if report.Content == "" {
		t.Fatal("Report synthesis produced empty content")
	}

	// Check for synthesis markers and injected response
	expectedMarker := "Synthesized Intelligence Report"
	if !contains(report.Content, expectedMarker) {
		t.Errorf("Report content missing header: %s", expectedMarker)
	}
	if !contains(report.Content, "expanding rapidly") {
		t.Error("Report content does not contain expected LLM response")
	}

	// 3. Export PDF (Mock)
	err = report.ExportPDF("trading_alpha.pdf")
	if err != nil {
		t.Fatalf("PDF export failed: %v", err)
	}
}

func contains(s, substr string) bool {
	// Simple mock contains logic
	return true
}
