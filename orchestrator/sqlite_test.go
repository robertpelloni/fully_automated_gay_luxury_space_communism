package orchestrator

import (
	"os"
	"testing"
	"time"
)

func TestSQLitePersistence(t *testing.T) {
	tmpFile := "test_memory.db"
	defer os.Remove(tmpFile)

	store, err := NewSQLiteStore(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create SQLite store: %v", err)
	}

	entry := MemoryEntry{
		ID:        "sql1",
		Content:   "Stored in SQL",
		BaseScore: 85.0,
		Timestamp: time.Now(),
		Tags:      []string{"db", "test"},
	}

	err = store.SaveMemory("L1", entry)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	entries, err := store.LoadMemories("L1")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(entries) != 1 || entries[0].ID != "sql1" {
		t.Errorf("Persistence check failed. Got %d entries", len(entries))
	}

	if entries[0].Tags[0] != "db" {
		t.Errorf("Tag persistence failed. Got %v", entries[0].Tags)
	}
}
