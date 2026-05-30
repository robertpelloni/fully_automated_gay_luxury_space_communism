package orchestrator

import (
	"fmt"
	"testing"
	"time"
)

func TestWeightedCouncil(t *testing.T) {
	orc := NewOrchestrator()
	council := NewCouncil(orc)

	proposal := "Test Weighted Voting"
	result, err := council.Debate(proposal)
	if err != nil {
		t.Fatalf("Debate failed: %v", err)
	}

	fmt.Printf("Weighted Score: %.2f\n", result.WeightedScore)
	if result.WeightedScore == 0 {
		t.Errorf("Expected non-zero weighted score")
	}
}

func TestRankedSearch(t *testing.T) {
	orch := NewOrchestrator()

	orch.L1.Add(MemoryEntry{ID: "old", Content: "Keyword", BaseScore: 10, Timestamp: time.Now().Add(-24 * time.Hour)})
	orch.L1.Add(MemoryEntry{ID: "fresh", Content: "Keyword", BaseScore: 10, Timestamp: time.Now()})

	results := orch.L1.RankedSearch("Keyword", orch.Embedder)
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	if results[0].ID != "fresh" {
		t.Errorf("RankedSearch failed. Expected 'fresh' to be first, got %s", results[0].ID)
	}
}
