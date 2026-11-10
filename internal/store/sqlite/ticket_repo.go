package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"peloton/internal/domain"
)

// TicketRepository handles database operations on service tickets.
type TicketRepository struct {
	db *DB
}

// NewTicketRepository constructs a new TicketRepository.
func NewTicketRepository(db *DB) *TicketRepository {
	return &TicketRepository{db: db}
}

// Create inserts a new service ticket into the database.
func (r *TicketRepository) Create(t *domain.ServiceTicket) error {
	if err := t.Validate(); err != nil {
		return err
	}

	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now

	query := `
		INSERT INTO service_tickets (
			client_id, description, status, priority,
			estimated_cost, final_cost, created_at, updated_at, completed_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := r.db.Exec(
		query,
		t.ClientID, t.Description, string(t.Status), t.Priority,
		t.EstimatedCost, t.FinalCost, t.CreatedAt, t.UpdatedAt, t.CompletedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert service ticket: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get ticket id: %w", err)
	}

	t.ID = id
	return nil
}

// GetByID finds a ticket with its client details joined.
func (r *TicketRepository) GetByID(id int64) (*domain.ServiceTicket, error) {
	query := `
		SELECT 
			t.id, t.client_id, c.name, c.phone, t.description, t.status, t.priority,
			t.estimated_cost, t.final_cost, t.created_at, t.updated_at, t.completed_at
		FROM service_tickets t
		JOIN clients c ON t.client_id = c.id
		WHERE t.id = ?
	`

	var t domain.ServiceTicket
	var statusStr string
	err := r.db.QueryRow(query, id).Scan(
		&t.ID, &t.ClientID, &t.ClientName, &t.ClientPhone, &t.Description, &statusStr, &t.Priority,
		&t.EstimatedCost, &t.FinalCost, &t.CreatedAt, &t.UpdatedAt, &t.CompletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("service ticket #%d not found", id)
		}
		return nil, fmt.Errorf("error querying service ticket: %w", err)
	}

	t.Status = domain.TicketStatus(statusStr)
	return &t, nil
}

// List returns tickets, optionally filtered by status.
func (r *TicketRepository) List(statusFilter domain.TicketStatus) ([]domain.ServiceTicket, error) {
	baseQuery := `
		SELECT 
			t.id, t.client_id, c.name, c.phone, t.description, t.status, t.priority,
			t.estimated_cost, t.final_cost, t.created_at, t.updated_at, t.completed_at
		FROM service_tickets t
		JOIN clients c ON t.client_id = c.id
	`
	var rows *sql.Rows
	var err error

	if statusFilter != "" {
		query := baseQuery + " WHERE t.status = ? ORDER BY t.created_at ASC"
		rows, err = r.db.Query(query, string(statusFilter))
	} else {
		query := baseQuery + " ORDER BY t.created_at DESC"
		rows, err = r.db.Query(query)
	}

	if err != nil {
		return nil, fmt.Errorf("error listing service tickets: %w", err)
	}
	defer rows.Close()

	var tickets []domain.ServiceTicket
	for rows.Next() {
		var t domain.ServiceTicket
		var statusStr string
		err := rows.Scan(
			&t.ID, &t.ClientID, &t.ClientName, &t.ClientPhone, &t.Description, &statusStr, &t.Priority,
			&t.EstimatedCost, &t.FinalCost, &t.CreatedAt, &t.UpdatedAt, &t.CompletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning ticket row: %w", err)
		}
		t.Status = domain.TicketStatus(statusStr)
		tickets = append(tickets, t)
	}

	return tickets, rows.Err()
}

// UpdateStatus changes a ticket's status and updates completion time if finished.
func (r *TicketRepository) UpdateStatus(id int64, newStatus domain.TicketStatus, finalCost *float64) error {
	if !domain.IsValidStatus(newStatus) {
		return fmt.Errorf("invalid ticket status: %s", newStatus)
	}

	now := time.Now()
	var completedAt *time.Time
	if newStatus == domain.StatusCompleted {
		completedAt = &now
	}

	query := `
		UPDATE service_tickets
		SET status = ?, updated_at = ?, completed_at = COALESCE(?, completed_at)
	`
	args := []interface{}{string(newStatus), now, completedAt}

	if finalCost != nil {
		query += ", final_cost = ?"
		args = append(args, *finalCost)
	}

	query += " WHERE id = ?"
	args = append(args, id)

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update ticket status: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("ticket #%d not found", id)
	}

	return nil
}
