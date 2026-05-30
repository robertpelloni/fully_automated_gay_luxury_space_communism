package orchestrator

import (
	"os"
	"testing"
	"time"
)

func TestMemoryScoring(t *testing.T) {
	entry := MemoryEntry{
		ID:        "test1",
		Content:   "Highly relevant info",
		BaseScore: 100.0,
		Timestamp: time.Now(),
	}

	score1 := entry.Score()
	if score1 <= 0 {
		t.Errorf("Expected positive score, got %v", score1)
	}

	// Move timestamp back by 24 hours
	entry.Timestamp = time.Now().Add(-24 * time.Hour)
	score2 := entry.Score()

	if score2 >= score1 {
		t.Errorf("Score should decay over time. Initial: %v, After 24h: %v", score1, score2)
	}
}

func TestMemoryTiers(t *testing.T) {
	orch := NewOrchestrator()
	entry := MemoryEntry{ID: "m1", Content: "L1 Data", Timestamp: time.Now()}

	orch.L1.Add(entry)
	results := orch.L1.Search("L1")

	if len(results) != 1 {
		t.Errorf("Expected 1 result in L1, got %d", len(results))
	}
}

func TestMemoryPersistence(t *testing.T) {
	tmpFile := "test_memory.json"
	defer os.Remove(tmpFile)

	orch := NewOrchestrator()
	orch.L1.Add(MemoryEntry{ID: "p1", Content: "Persist Me", Timestamp: time.Now()})

	err := orch.Save(tmpFile)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	orch2 := NewOrchestrator()
	err = orch2.Load(tmpFile)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(orch2.L1.Entries) != 1 || orch2.L1.Entries[0].ID != "p1" {
		t.Errorf("Persistence failed. Expected 1 entry with ID p1, got %d entries", len(orch2.L1.Entries))
	}
}
