package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultDBPath_EnvOverride(t *testing.T) {
	expected := "/tmp/test-peloton-custom.db"
	os.Setenv("PELOTON_DB_PATH", expected)
	defer os.Unsetenv("PELOTON_DB_PATH")

	if got := DefaultDBPath(); got != expected {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}

func TestDefaultDBPath_XDGDefault(t *testing.T) {
	os.Unsetenv("PELOTON_DB_PATH")
	tempHome := t.TempDir()
	os.Setenv("HOME", tempHome)
	os.Unsetenv("XDG_DATA_HOME")

	got := DefaultDBPath()
	expected := filepath.Join(tempHome, ".local", "share", "peloton", "peloton.db")
	if got != expected {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}

func TestNewConfig_CreatesDirectory(t *testing.T) {
	tempDir := t.TempDir()
	targetDB := filepath.Join(tempDir, "nested", "storage", "peloton.db")

	cfg, err := NewConfig(targetDB, 9090, "127.0.0.1")
	if err != nil {
		t.Fatalf("NewConfig failed: %v", err)
	}

	if cfg.DBPath != targetDB {
		t.Errorf("expected DBPath %s, got %s", targetDB, cfg.DBPath)
	}
	if cfg.Port != 9090 {
		t.Errorf("expected Port 9090, got %d", cfg.Port)
	}
	if cfg.Host != "127.0.0.1" {
		t.Errorf("expected Host 127.0.0.1, got %s", cfg.Host)
	}

	dirInfo, err := os.Stat(filepath.Dir(targetDB))
	if err != nil || !dirInfo.IsDir() {
		t.Fatalf("expected directory to exist: %v", err)
	}
}
