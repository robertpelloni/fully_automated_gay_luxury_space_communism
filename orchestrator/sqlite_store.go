package orchestrator

import (
	"database/sql"
	"encoding/json"
	"math"
	"sort"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(filepath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	// Standard table for metadata
	query := `
	CREATE TABLE IF NOT EXISTS memories (
		id TEXT PRIMARY KEY,
		tier TEXT,
		content TEXT,
		base_score REAL,
		timestamp DATETIME,
		tags TEXT,
		embedding BLOB
	);`
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) SaveMemory(tier string, entry MemoryEntry) error {
	tagsData, err := json.Marshal(entry.Tags)
	if err != nil {
		return err
	}

	var embeddingData []byte
	if len(entry.Vector) > 0 {
		embeddingData, err = json.Marshal(entry.Vector)
		if err != nil {
			return err
		}
	}

	_, err = s.db.Exec(`
		INSERT OR REPLACE INTO memories (id, tier, content, base_score, timestamp, tags, embedding)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		entry.ID, tier, entry.Content, entry.BaseScore, entry.Timestamp, string(tagsData), embeddingData)
	return err
}

func (s *SQLiteStore) LoadMemories(tier string) ([]MemoryEntry, error) {
	rows, err := s.db.Query("SELECT id, content, base_score, timestamp, tags, embedding FROM memories WHERE tier = ?", tier)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []MemoryEntry
	for rows.Next() {
		var e MemoryEntry
		var tagsStr string
		var embeddingData []byte
		err := rows.Scan(&e.ID, &e.Content, &e.BaseScore, &e.Timestamp, &tagsStr, &embeddingData)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(tagsStr), &e.Tags); err != nil {
			continue
		}
		if len(embeddingData) > 0 {
			if err := json.Unmarshal(embeddingData, &e.Vector); err != nil {
				continue
			}
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// VectorSearch finds similar memories using cosine similarity.
// In the future, this will utilize the sqlite-vec extension for SQL-level performance.
func (s *SQLiteStore) VectorSearch(tier string, target []float32, limit int) ([]MemoryEntry, error) {
	memories, err := s.LoadMemories(tier)
	if err != nil {
		return nil, err
	}

	type match struct {
		entry MemoryEntry
		score float64
	}
	var matches []match

	for _, m := range memories {
		if len(m.Vector) == 0 || len(m.Vector) != len(target) {
			continue
		}
		score := cosineSimilarity(m.Vector, target)
		matches = append(matches, match{m, score})
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].score > matches[j].score
	})

	var results []MemoryEntry
	for i := 0; i < len(matches) && i < limit; i++ {
		results = append(results, matches[i].entry)
	}
	return results, nil
}

func cosineSimilarity(a, b []float32) float64 {
	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
