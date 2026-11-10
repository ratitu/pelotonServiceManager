package sqlite

import (
	"path/filepath"
	"testing"
	"time"

	"peloton/internal/domain"
)

func setupTestStore(t *testing.T) (*DB, func()) {
	t.Helper()
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := OpenDB(dbPath)
	if err != nil {
		t.Fatalf("OpenDB failed: %v", err)
	}

	return db, func() {
		db.Close()
	}
}

func TestStore_CompleteLifecycle(t *testing.T) {
	db, cleanup := setupTestStore(t)
	defer cleanup()

	clientRepo := NewClientRepository(db)
	ticketRepo := NewTicketRepository(db)
	financeRepo := NewFinanceRepository(db)

	// 1. Create client
	client := &domain.Client{
		Name:        "Marek Nowak",
		Phone:       "+48 600 700 800",
		BikeDetails: "Trek Domane AL 4 Disc (Blue)",
		Notes:       "Prefers morning pickup",
	}

	if err := clientRepo.Create(client); err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	if client.ID == 0 {
		t.Fatal("expected non-zero client id")
	}

	// Read back client
	retrievedClient, err := clientRepo.GetByID(client.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if retrievedClient.Name != client.Name {
		t.Errorf("expected client name %s, got %s", client.Name, retrievedClient.Name)
	}

	// Search client
	searchResults, err := clientRepo.Search("Trek")
	if err != nil || len(searchResults) != 1 {
		t.Fatalf("Search failed, results: %v (err: %v)", searchResults, err)
	}

	// 2. Create service ticket
	ticket := &domain.ServiceTicket{
		ClientID:      client.ID,
		Description:   "Hydraulic brakes bleeding, cassette replacement",
		Status:        domain.StatusPending,
		Priority:      "High",
		EstimatedCost: 240.0,
	}

	if err := ticketRepo.Create(ticket); err != nil {
		t.Fatalf("failed to create ticket: %v", err)
	}

	// 3. Status progression: Pending -> In Progress -> Waiting for Parts -> Ready for Pickup -> Completed
	statuses := []domain.TicketStatus{
		domain.StatusInProgress,
		domain.StatusWaitingParts,
		domain.StatusReadyPickup,
	}

	for _, s := range statuses {
		if err := ticketRepo.UpdateStatus(ticket.ID, s, nil); err != nil {
			t.Fatalf("UpdateStatus to %s failed: %v", s, err)
		}
	}

	finalCost := 260.0
	if err := ticketRepo.UpdateStatus(ticket.ID, domain.StatusCompleted, &finalCost); err != nil {
		t.Fatalf("UpdateStatus to Completed failed: %v", err)
	}

	// Verify completion
	completedTicket, err := ticketRepo.GetByID(ticket.ID)
	if err != nil {
		t.Fatalf("failed to get completed ticket: %v", err)
	}
	if completedTicket.Status != domain.StatusCompleted {
		t.Errorf("expected status Completed, got %s", completedTicket.Status)
	}
	if completedTicket.FinalCost != 260.0 {
		t.Errorf("expected final cost 260.0, got %.2f", completedTicket.FinalCost)
	}
	if completedTicket.CompletedAt == nil {
		t.Error("expected completed_at timestamp to be set")
	}

	// 4. Financial tracking for current month
	now := time.Now()
	// Parts expense
	partExpense := &domain.FinancialTransaction{
		TicketID: &ticket.ID,
		Type:     domain.TypeExpense,
		Category: "Parts",
		Amount:   110.0,
		Notes:    "Shimano 11-speed cassette",
		Date:     now,
	}
	if err := financeRepo.Create(partExpense); err != nil {
		t.Fatalf("failed to log expense: %v", err)
	}

	// Service Income
	serviceIncome := &domain.FinancialTransaction{
		TicketID: &ticket.ID,
		Type:     domain.TypeIncome,
		Category: "Service",
		Amount:   260.0,
		Notes:    "Customer paid in cash",
		Date:     now,
	}
	if err := financeRepo.Create(serviceIncome); err != nil {
		t.Fatalf("failed to log income: %v", err)
	}

	// 5. Monthly balance check
	summary, err := financeRepo.GetMonthlySummary(now.Year(), int(now.Month()))
	if err != nil {
		t.Fatalf("GetMonthlySummary failed: %v", err)
	}

	if summary.TotalIncome != 260.0 {
		t.Errorf("expected total income 260.0, got %.2f", summary.TotalIncome)
	}
	if summary.TotalExpense != 110.0 {
		t.Errorf("expected total expense 110.0, got %.2f", summary.TotalExpense)
	}
	if summary.NetEarnings != 150.0 {
		t.Errorf("expected net earnings 150.0, got %.2f", summary.NetEarnings)
	}
	if summary.ServicesCount != 1 {
		t.Errorf("expected 1 completed service, got %d", summary.ServicesCount)
	}
}

func TestStore_ForeignKeysEnforced(t *testing.T) {
	db, cleanup := setupTestStore(t)
	defer cleanup()

	ticketRepo := NewTicketRepository(db)

	// Attempt to create ticket for non-existent client ID 9999
	ticket := &domain.ServiceTicket{
		ClientID:    9999,
		Description: "Ghost bike repair",
		Status:      domain.StatusPending,
	}

	err := ticketRepo.Create(ticket)
	if err == nil {
		t.Fatal("expected foreign key constraint error, got nil")
	}
}
