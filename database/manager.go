package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite" // Import the pure-Go driver
)

type Manager struct {
	db *sql.DB
}

func NewManager(dbPath string) (*Manager, error) {
	// Initialize the connection
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	m := &Manager{db: db}
	if err := m.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	log.Printf("[DATABASE] Audit trail initialized at %s", dbPath)
	return m, nil
}

func (m *Manager) initSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS audit_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		target TEXT NOT NULL,
		provider TEXT NOT NULL,
		raw_data TEXT NOT NULL,
		risk_score INTEGER DEFAULT 0,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := m.db.Exec(query)
	return err
}

func (m *Manager) SaveAssessment(target string, provider string, rawData map[string]interface{}, score int) error {
	// Convert the dynamic map back to a JSON string for storage
	jsonData, err := json.Marshal(rawData)
	if err != nil {
		return fmt.Errorf("failed to marshal raw data: %w", err)
	}

	query := `INSERT INTO audit_logs (target, provider, raw_data, risk_score, timestamp) VALUES (?, ?, ?, ?, ?)`

	_, err = m.db.Exec(query, target, provider, string(jsonData), score, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}

	return nil
}
