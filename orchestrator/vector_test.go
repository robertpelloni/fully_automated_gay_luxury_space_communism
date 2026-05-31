package orchestrator

import (
	"os"
	"testing"
	"time"
)

func TestVectorSearch(t *testing.T) {
	dbFile := "test_vector.db"
	defer os.Remove(dbFile)

	store, err := NewSQLiteStore(dbFile)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	e1 := MemoryEntry{
		ID:        "1",
		Content:   "AI Orchestration",
		Timestamp: time.Now(),
		Vector:    []float32{1.0, 0.0, 0.0},
	}
	e2 := MemoryEntry{
		ID:        "2",
		Content:   "Banana Bread",
		Timestamp: time.Now(),
		Vector:    []float32{0.0, 1.0, 0.0},
	}

	store.SaveMemory("L1", e1)
	store.SaveMemory("L1", e2)

	// Search for something close to [1, 0, 0]
	results, err := store.VectorSearch("L1", []float32{0.9, 0.1, 0.0}, 1)
	if err != nil {
		t.Fatalf("VectorSearch failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].ID != "1" {
		t.Errorf("Expected ID '1' (AI), got '%s'", results[0].ID)
	}
}
