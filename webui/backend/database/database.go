package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps the sql.DB with additional methods
type DB struct {
	*sql.DB
}

// Initialize creates and initializes the database
func Initialize(dbPath string) (*DB, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{sqlDB}

	// Run migrations
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// migrate runs database migrations
func (db *DB) migrate() error {
	migrations := []string{
		// Scans table
		`CREATE TABLE IF NOT EXISTS scans (
			id TEXT PRIMARY KEY,
			timestamp DATETIME NOT NULL,
			project_path TEXT NOT NULL,
			commit_hash TEXT,
			branch TEXT,
			duration INTEGER,
			dast_enabled BOOLEAN DEFAULT 0,
			summary_json TEXT,
			status TEXT DEFAULT 'completed',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Findings table
		`CREATE TABLE IF NOT EXISTS findings (
			id TEXT PRIMARY KEY,
			scan_id TEXT NOT NULL,
			title TEXT NOT NULL,
			description TEXT,
			severity TEXT NOT NULL,
			cwe TEXT,
			category TEXT,
			file_path TEXT,
			line_number INTEGER,
			line_content TEXT,
			remediation TEXT,
			confidence TEXT,
			first_detected DATETIME,
			last_seen DATETIME,
			status TEXT DEFAULT 'new',
			FOREIGN KEY (scan_id) REFERENCES scans(id) ON DELETE CASCADE
		)`,

		// DAST Findings table
		`CREATE TABLE IF NOT EXISTS dast_findings (
			id TEXT PRIMARY KEY,
			scan_id TEXT NOT NULL,
			endpoint TEXT NOT NULL,
			method TEXT NOT NULL,
			title TEXT NOT NULL,
			description TEXT,
			severity TEXT NOT NULL,
			request_payload TEXT,
			response_code INTEGER,
			response_body TEXT,
			evidence TEXT,
			screenshot TEXT,
			remediation TEXT,
			timestamp DATETIME,
			FOREIGN KEY (scan_id) REFERENCES scans(id) ON DELETE CASCADE
		)`,

		// API Tokens table
		`CREATE TABLE IF NOT EXISTS api_tokens (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			token TEXT UNIQUE NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			expires_at DATETIME,
			last_used DATETIME
		)`,

		// Users table (for authentication)
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			email TEXT,
			role TEXT DEFAULT 'user',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Indexes for better query performance
		`CREATE INDEX IF NOT EXISTS idx_scans_timestamp ON scans(timestamp DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_scans_status ON scans(status)`,
		`CREATE INDEX IF NOT EXISTS idx_findings_scan_id ON findings(scan_id)`,
		`CREATE INDEX IF NOT EXISTS idx_findings_severity ON findings(severity)`,
		`CREATE INDEX IF NOT EXISTS idx_findings_category ON findings(category)`,
		`CREATE INDEX IF NOT EXISTS idx_dast_findings_scan_id ON dast_findings(scan_id)`,
		`CREATE INDEX IF NOT EXISTS idx_api_tokens_token ON api_tokens(token)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}
