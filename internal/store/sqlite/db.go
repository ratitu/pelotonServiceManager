package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// DB wraps the standard *sql.DB connection.
type DB struct {
	*sql.DB
}

// OpenDB opens a connection to the SQLite database and executes baseline PRAGMAs.
func OpenDB(path string) (*DB, error) {
	// Enable WAL mode, busy timeout, and foreign keys via DSN parameters
	dsn := fmt.Sprintf("%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)", path)
	if path == ":memory:" {
		dsn = ":memory:"
	}

	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	// SQLite operates best with a single writer connection
	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)
	conn.SetConnMaxLifetime(time.Hour)

	// Ensure foreign keys are strictly active
	if _, err := conn.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to enforce foreign keys: %w", err)
	}

	db := &DB{DB: conn}

	if err := db.Migrate(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("database migration failed: %w", err)
	}

	return db, nil
}
