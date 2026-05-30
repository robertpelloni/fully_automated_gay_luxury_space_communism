package orchestrator

import (
	"encoding/json"
	"math"
	"os"
	"strings"
	"time"
)

// MemoryEntry represents a single unit of information
type MemoryEntry struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	BaseScore float64   `json:"base_score"`
	Timestamp time.Time `json:"timestamp"`
	Tags      []string  `json:"tags"`
	Vector    []float32 `json:"vector,omitempty"` // Placeholder for semantic embeddings
}

// Score calculates the current "heat" of the memory with temporal decay
func (e MemoryEntry) Score() float64 {
	elapsed := time.Since(e.Timestamp).Hours()
	// Simple exponential decay: Score = BaseScore * e^(-0.1 * hours)
	return e.BaseScore * math.Exp(-0.1*elapsed)
}

// L1Memory (Scratchpad) - Current session context
type L1Memory struct {
	Entries []MemoryEntry `json:"entries"`
}

func (m *L1Memory) Add(entry MemoryEntry) {
	m.Entries = append(m.Entries, entry)
}

func (m *L1Memory) Search(query string) []MemoryEntry {
	var results []MemoryEntry
	for _, entry := range m.Entries {
		if strings.Contains(strings.ToLower(entry.Content), strings.ToLower(query)) {
			results = append(results, entry)
		}
	}
	return results
}

// L2Memory (Vault) - Long-term historical patterns
type L2Memory struct {
	Entries []MemoryEntry `json:"entries"`
}

func (m *L2Memory) Add(entry MemoryEntry) {
	m.Entries = append(m.Entries, entry)
}

func (m *L2Memory) Search(query string) []MemoryEntry {
	var results []MemoryEntry
	for _, entry := range m.Entries {
		if strings.Contains(strings.ToLower(entry.Content), strings.ToLower(query)) {
			results = append(results, entry)
		}
	}
	return results
}

// L3Memory (Archive) - Archived decisions and cold storage
type L3Memory struct {
	Entries []MemoryEntry `json:"entries"`
}

func (m *L3Memory) Add(entry MemoryEntry) {
	m.Entries = append(m.Entries, entry)
}

func (m *L3Memory) Search(query string) []MemoryEntry {
	var results []MemoryEntry
	for _, entry := range m.Entries {
		if strings.Contains(strings.ToLower(entry.Content), strings.ToLower(query)) {
			results = append(results, entry)
		}
	}
	return results
}

// Orchestrator handles tiered memory orchestration
type Orchestrator struct {
	L1 L1Memory `json:"l1"`
	L2 L2Memory `json:"l2"`
	L3 L3Memory `json:"l3"`
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		L1: L1Memory{Entries: make([]MemoryEntry, 0)},
		L2: L2Memory{Entries: make([]MemoryEntry, 0)},
		L3: L3Memory{Entries: make([]MemoryEntry, 0)},
	}
}

func (o *Orchestrator) Save(filepath string) error {
	data, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, data, 0644)
}

func (o *Orchestrator) Load(filepath string) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil // No file yet, which is fine
	}
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, o)
}
