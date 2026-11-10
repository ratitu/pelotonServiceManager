package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config represents runtime configuration settings.
type Config struct {
	DBPath string
	Port   int
	Host   string
}

// DefaultDBPath resolves the database file path according to:
// 1. PELOTON_DB_PATH environment variable
// 2. $XDG_DATA_HOME/peloton/peloton.db (typically ~/.local/share/peloton/peloton.db)
// 3. Fallback to ./peloton.db
func DefaultDBPath() string {
	if envPath := os.Getenv("PELOTON_DB_PATH"); envPath != "" {
		return envPath
	}

	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, err := os.UserHomeDir()
		if err == nil && home != "" {
			dataHome = filepath.Join(home, ".local", "share")
		}
	}

	if dataHome != "" {
		return filepath.Join(dataHome, "peloton", "peloton.db")
	}

	return "peloton.db"
}

// EnsureDataDir ensures that the directory containing the database file exists.
func EnsureDataDir(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if dir == "" || dir == "." {
		return nil
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory %s: %w", dir, err)
	}
	return nil
}

// NewConfig initializes configuration with defaults and overrides.
func NewConfig(overrideDBPath string, port int, host string) (*Config, error) {
	dbPath := overrideDBPath
	if dbPath == "" {
		dbPath = DefaultDBPath()
	}

	if err := EnsureDataDir(dbPath); err != nil {
		return nil, err
	}

	if port <= 0 {
		port = 8080
	}
	if host == "" {
		host = "0.0.0.0"
	}

	return &Config{
		DBPath: dbPath,
		Port:   port,
		Host:   host,
	}, nil
}
