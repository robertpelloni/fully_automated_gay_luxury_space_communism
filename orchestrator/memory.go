package orchestrator

import (
	"time"
)

// MemoryEntry represents a single unit of information
type MemoryEntry struct {
	ID        string
	Content   string
	Score     float64
	Timestamp time.Time
	Tags      []string
}

// L1Memory (Scratchpad) - Current session context
type L1Memory struct {
	Entries []MemoryEntry
}

// L2Memory (Vault) - Long-term historical patterns
type L2Memory struct {
	Entries []MemoryEntry
}

// L3Memory (Archive) - Archived decisions and cold storage
type L3Memory struct {
	Entries []MemoryEntry
}

// Orchestrator handles tiered memory orchestration
type Orchestrator struct {
	L1 L1Memory
	L2 L2Memory
	L3 L3Memory
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		L1: L1Memory{Entries: make([]MemoryEntry, 0)},
		L2: L2Memory{Entries: make([]MemoryEntry, 0)},
		L3: L3Memory{Entries: make([]MemoryEntry, 0)},
	}
}
