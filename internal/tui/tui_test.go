package tui

import (
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"peloton/internal/domain"
	"peloton/internal/store/sqlite"
)

func setupTestTUI(t *testing.T) (Model, *sqlite.DB, func()) {
	t.Helper()
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_tui.db")

	db, err := sqlite.OpenDB(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	clientRepo := sqlite.NewClientRepository(db)
	ticketRepo := sqlite.NewTicketRepository(db)
	financeRepo := sqlite.NewFinanceRepository(db)

	// Seed client and tickets
	client := &domain.Client{Name: "Jan Kowalski", Phone: "555-100", BikeDetails: "Trek Domane"}
	if err := clientRepo.Create(client); err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	t1 := &domain.ServiceTicket{
		ClientID:      client.ID,
		Description:   "Brake pads replacement",
		Status:        domain.StatusPending,
		Priority:      "High",
		EstimatedCost: 80.0,
	}
	t2 := &domain.ServiceTicket{
		ClientID:      client.ID,
		Description:   "Drivetrain cleaning",
		Status:        domain.StatusPending,
		Priority:      "Normal",
		EstimatedCost: 60.0,
	}
	t3 := &domain.ServiceTicket{
		ClientID:      client.ID,
		Description:   "Tire swap",
		Status:        domain.StatusInProgress,
		Priority:      "Urgent",
		EstimatedCost: 40.0,
	}
	_ = ticketRepo.Create(t1)
	_ = ticketRepo.Create(t2)
	_ = ticketRepo.Create(t3)

	model := NewModel(ticketRepo, financeRepo, clientRepo)

	cleanup := func() {
		db.Close()
	}

	return model, db, cleanup
}

func TestTUI_Navigation(t *testing.T) {
	model, _, cleanup := setupTestTUI(t)
	defer cleanup()

	// Initial state: Column 0 (Pending) should have 2 tickets, Column 1 (In Progress) has 1
	if len(model.Columns[0].Tickets) != 2 {
		t.Fatalf("expected 2 pending tickets, got %d", len(model.Columns[0].Tickets))
	}
	if len(model.Columns[1].Tickets) != 1 {
		t.Fatalf("expected 1 in-progress ticket, got %d", len(model.Columns[1].Tickets))
	}

	// 1. Move Cursor Down in Column 0
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	m := updatedModel.(Model)
	if m.Columns[0].CursorIndex != 1 {
		t.Errorf("expected cursor index 1 after down arrow, got %d", m.Columns[0].CursorIndex)
	}

	// 2. Move Right (Column 0 -> Column 1)
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = updatedModel.(Model)
	if m.ActiveColIdx != 1 {
		t.Errorf("expected active column index 1 after right arrow, got %d", m.ActiveColIdx)
	}

	// 3. Tab to cycle columns
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updatedModel.(Model)
	if m.ActiveColIdx != 2 {
		t.Errorf("expected active column index 2 after tab, got %d", m.ActiveColIdx)
	}

	// 4. Test View rendering
	m.Width = 120
	m.Height = 30
	viewStr := m.View()
	if !strings.Contains(viewStr, "PELOTON WORKSHOP MANAGER") {
		t.Errorf("view missing expected header")
	}
	if !strings.Contains(viewStr, "PENDING (2)") {
		t.Errorf("view missing column header: PENDING (2)")
	}
	if !strings.Contains(viewStr, "Brake pads") {
		t.Errorf("view missing ticket card content")
	}
}
