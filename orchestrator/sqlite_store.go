package orchestrator

import (
	"database/sql"
	"encoding/json"
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

	query := `
	CREATE TABLE IF NOT EXISTS memories (
		id TEXT PRIMARY KEY,
		tier TEXT,
		content TEXT,
		base_score REAL,
		timestamp DATETIME,
		tags TEXT
	);`
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) SaveMemory(tier string, entry MemoryEntry) error {
	tagsData, _ := json.Marshal(entry.Tags)
	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO memories (id, tier, content, base_score, timestamp, tags)
		VALUES (?, ?, ?, ?, ?, ?)`,
		entry.ID, tier, entry.Content, entry.BaseScore, entry.Timestamp, string(tagsData))
	return err
}

func (s *SQLiteStore) LoadMemories(tier string) ([]MemoryEntry, error) {
	rows, err := s.db.Query("SELECT id, content, base_score, timestamp, tags FROM memories WHERE tier = ?", tier)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []MemoryEntry
	for rows.Next() {
		var e MemoryEntry
		var tagsStr string
		err := rows.Scan(&e.ID, &e.Content, &e.BaseScore, &e.Timestamp, &tagsStr)
		if err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(tagsStr), &e.Tags)
		entries = append(entries, e)
	}
	return entries, nil
}
