package main

import (
	"fmt"
	"github.com/robertpelloni/hustle/hustle/research"
	"github.com/robertpelloni/hustle/orchestrator"
)

func main() {
	fmt.Println("=== Hustle Research Module ===")
	orch := orchestrator.NewOrchestrator()
	searcher := &research.ResearchSearch{
		ActiveProvider: research.Tavily,
		Orchestrator:   orch,
	}

	query := "Latest trends in AI agent orchestration"
	fmt.Printf("Running research for: %s\n", query)

	results, err := searcher.Query(query)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	report := &research.Report{Title: "AI Trends 2026"}
	report.Synthesize(orch, results)

	fmt.Println("Synthesis Complete.")
	fmt.Println(report.Content)
}
