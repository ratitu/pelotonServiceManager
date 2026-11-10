package domain

import (
	"errors"
	"time"
)

// TicketStatus represents the lifecycle state of a service repair ticket.
type TicketStatus string

const (
	StatusPending      TicketStatus = "Pending"
	StatusInProgress   TicketStatus = "In Progress"
	StatusWaitingParts TicketStatus = "Waiting for Parts"
	StatusReadyPickup  TicketStatus = "Ready for Pickup"
	StatusCompleted    TicketStatus = "Completed"
)

// AllTicketStatuses lists every valid ticket state in logical order.
var AllTicketStatuses = []TicketStatus{
	StatusPending,
	StatusInProgress,
	StatusWaitingParts,
	StatusReadyPickup,
	StatusCompleted,
}

// IsValidStatus checks if a given status string is recognized.
func IsValidStatus(status TicketStatus) bool {
	switch status {
	case StatusPending, StatusInProgress, StatusWaitingParts, StatusReadyPickup, StatusCompleted:
		return true
	default:
		return false
	}
}

// Client represents a bicycle owner.
type Client struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	BikeDetails string    `json:"bike_details"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate checks client required fields.
func (c *Client) Validate() error {
	if c.Name == "" {
		return errors.New("client name is required")
	}
	if c.Phone == "" {
		return errors.New("client phone is required")
	}
	return nil
}

// ServiceTicket represents an active or past repair job.
type ServiceTicket struct {
	ID            int64        `json:"id"`
	ClientID      int64        `json:"client_id"`
	ClientName    string       `json:"client_name,omitempty"`
	ClientPhone   string       `json:"client_phone,omitempty"`
	Description   string       `json:"description"`
	Status        TicketStatus `json:"status"`
	Priority      string       `json:"priority"` // Normal, High, Urgent
	EstimatedCost float64      `json:"estimated_cost"`
	FinalCost     float64      `json:"final_cost"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	CompletedAt   *time.Time   `json:"completed_at,omitempty"`
}

// Validate checks service ticket validity.
func (t *ServiceTicket) Validate() error {
	if t.ClientID <= 0 {
		return errors.New("client_id is required")
	}
	if t.Description == "" {
		return errors.New("ticket description is required")
	}
	if !IsValidStatus(t.Status) {
		return errors.New("invalid ticket status")
	}
	if t.Priority == "" {
		t.Priority = "Normal"
	}
	return nil
}

// TransactionType indicates cashflow direction.
type TransactionType string

const (
	TypeIncome  TransactionType = "INCOME"
	TypeExpense TransactionType = "EXPENSE"
)

// FinancialTransaction represents a logged income or expense entry.
type FinancialTransaction struct {
	ID        int64           `json:"id"`
	TicketID  *int64          `json:"ticket_id,omitempty"`
	Type      TransactionType `json:"type"`
	Category  string          `json:"category"`
	Amount    float64         `json:"amount"`
	Notes     string          `json:"notes"`
	Date      time.Time       `json:"date"`
}

// MonthlyBalance represents cash in vs cash out for a given month.
type MonthlyBalance struct {
	Year          int     `json:"year"`
	Month         int     `json:"month"`
	TotalIncome   float64 `json:"total_income"`
	TotalExpense  float64 `json:"total_expense"`
	NetEarnings   float64 `json:"net_earnings"`
	ServicesCount int     `json:"services_count"`
}
