package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestApp(t *testing.T) (*App, *bytes.Buffer, *bytes.Buffer, func()) {
	t.Helper()
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_cli.db")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	app, err := NewApp(dbPath, stdout, stderr)
	if err != nil {
		t.Fatalf("failed to create test app: %v", err)
	}

	cleanup := func() {
		app.Close()
	}

	return app, stdout, stderr, cleanup
}

func TestClientCLI(t *testing.T) {
	app, stdout, stderr, cleanup := setupTestApp(t)
	defer cleanup()

	// 1. Add client
	code := app.Run([]string{"client", "add", "--name", "Piotr Zielinski", "--phone", "+48 501 234 567", "--bike", "Canyon Endurace CF 7", "--notes", "Tubeless setup"})
	if code != 0 {
		t.Fatalf("client add failed with exit code %d. Stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Client registered successfully") {
		t.Errorf("expected success message in stdout, got: %s", stdout.String())
	}

	// 2. List clients
	stdout.Reset()
	stderr.Reset()
	code = app.Run([]string{"client", "list"})
	if code != 0 {
		t.Fatalf("client list failed with code %d", code)
	}
	if !strings.Contains(stdout.String(), "Piotr Zielinski") || !strings.Contains(stdout.String(), "Canyon Endurace") {
		t.Errorf("client list missing registered client. Output:\n%s", stdout.String())
	}

	// 3. Search client
	stdout.Reset()
	stderr.Reset()
	code = app.Run([]string{"client", "search", "Endurace"})
	if code != 0 {
		t.Fatalf("client search failed with code %d", code)
	}
	if !strings.Contains(stdout.String(), "Piotr Zielinski") {
		t.Errorf("search results missing expected client. Output:\n%s", stdout.String())
	}
}

func TestServiceCLI(t *testing.T) {
	app, stdout, stderr, cleanup := setupTestApp(t)
	defer cleanup()

	// Setup client
	app.Run([]string{"client", "add", "--name", "Anna Lewandowska", "--phone", "555-1234", "--bike", "Specialized Roubaix"})

	// 1. Add service ticket
	stdout.Reset()
	stderr.Reset()
	code := app.Run([]string{"service", "add", "--client", "1", "--desc", "Full transmission tune-up and chain replacement", "--priority", "High", "--est", "180.0"})
	if code != 0 {
		t.Fatalf("service add failed: %s", stderr.String())
	}
	if !strings.Contains(stdout.String(), "Service Ticket #1 created") {
		t.Errorf("unexpected stdout: %s", stdout.String())
	}

	// 2. List tickets
	stdout.Reset()
	stderr.Reset()
	code = app.Run([]string{"service", "list"})
	if code != 0 {
		t.Fatalf("service list failed: %s", stderr.String())
	}
	if !strings.Contains(stdout.String(), "Anna Lewandowska") || !strings.Contains(stdout.String(), "Pending") {
		t.Errorf("service list output missing expected content:\n%s", stdout.String())
	}

	// 3. Status progression
	stdout.Reset()
	stderr.Reset()
	code = app.Run([]string{"service", "status", "1", "In Progress"})
	if code != 0 {
		t.Fatalf("service status transition failed: %s", stderr.String())
	}

	// 4. Complete with final cost
	stdout.Reset()
	stderr.Reset()
	code = app.Run([]string{"service", "status", "1", "Completed", "--cost", "195.50"})
	if code != 0 {
		t.Fatalf("service complete failed: %s", stderr.String())
	}
	if !strings.Contains(stdout.String(), "Completed") || !strings.Contains(stdout.String(), "195.50") {
		t.Errorf("unexpected completion stdout:\n%s", stdout.String())
	}

	// 5. Verify filter list
	stdout.Reset()
	stderr.Reset()
	code = app.Run([]string{"service", "list", "--status", "Completed"})
	if code != 0 {
		t.Fatalf("service list --status failed: %s", stderr.String())
	}
	if !strings.Contains(stdout.String(), "Completed") {
		t.Errorf("filtered list missing completed ticket:\n%s", stdout.String())
	}
}

func TestCLI_HelpAndVersion(t *testing.T) {
	app, stdout, _, cleanup := setupTestApp(t)
	defer cleanup()

	code := app.Run([]string{"--version"})
	if code != 0 || !strings.Contains(stdout.String(), "Peloton Service Manager") {
		t.Errorf("version output check failed")
	}

	stdout.Reset()
	code = app.Run([]string{"--help"})
	if code != 0 || !strings.Contains(stdout.String(), "USAGE:") {
		t.Errorf("help output check failed")
	}
}
