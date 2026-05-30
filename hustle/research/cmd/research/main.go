package main

import (
	"fmt"
	"github.com/robertpelloni/hustle/hustle/research"
)

func main() {
	fmt.Println("=== Hustle Research Module ===")
	searcher := &research.ResearchSearch{ActiveProvider: research.Tavily}

	query := "Latest trends in AI agent orchestration"
	fmt.Printf("Running research for: %s\n", query)

	results, err := searcher.Query(query)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	report := &research.Report{Title: "AI Trends 2026"}
	report.Synthesize(results)

	fmt.Println("Synthesis Complete.")
	fmt.Println(report.Content)
}
