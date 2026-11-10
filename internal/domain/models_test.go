package domain

import (
	"testing"
)

func TestTicketStatus_Validation(t *testing.T) {
	valid := []TicketStatus{
		StatusPending,
		StatusInProgress,
		StatusWaitingParts,
		StatusReadyPickup,
		StatusCompleted,
	}

	for _, s := range valid {
		if !IsValidStatus(s) {
			t.Errorf("expected %s to be valid", s)
		}
	}

	invalid := []TicketStatus{
		"Unknown",
		"",
		"DONE",
	}

	for _, s := range invalid {
		if IsValidStatus(s) {
			t.Errorf("expected %s to be invalid", s)
		}
	}
}

func TestClient_Validation(t *testing.T) {
	c := &Client{}
	if err := c.Validate(); err == nil {
		t.Error("expected error for empty client")
	}

	c.Name = "Jan Kowalski"
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing phone")
	}

	c.Phone = "+48 123 456 789"
	if err := c.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestServiceTicket_Validation(t *testing.T) {
	ticket := &ServiceTicket{
		ClientID:    1,
		Description: "Brake bleeding and cable replacement",
		Status:      StatusPending,
	}

	if err := ticket.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if ticket.Priority != "Normal" {
		t.Errorf("expected default priority Normal, got %s", ticket.Priority)
	}

	ticket.Status = "InvalidStatus"
	if err := ticket.Validate(); err == nil {
		t.Error("expected error for invalid status")
	}
}
