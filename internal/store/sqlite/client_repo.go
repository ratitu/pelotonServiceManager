package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"peloton/internal/domain"
)

// ClientRepository handles database operations on clients.
type ClientRepository struct {
	db *DB
}

// NewClientRepository constructs a new ClientRepository.
func NewClientRepository(db *DB) *ClientRepository {
	return &ClientRepository{db: db}
}

// Create inserts a new client into the database.
func (r *ClientRepository) Create(c *domain.Client) error {
	if err := c.Validate(); err != nil {
		return err
	}

	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	query := `
		INSERT INTO clients (name, phone, bike_details, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	res, err := r.db.Exec(query, c.Name, c.Phone, c.BikeDetails, c.Notes, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert client: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve client id: %w", err)
	}

	c.ID = id
	return nil
}

// GetByID finds a client by their integer ID.
func (r *ClientRepository) GetByID(id int64) (*domain.Client, error) {
	query := `
		SELECT id, name, phone, bike_details, notes, created_at, updated_at
		FROM clients
		WHERE id = ?
	`

	var c domain.Client
	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.Name, &c.Phone, &c.BikeDetails, &c.Notes, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("client with ID %d not found", id)
		}
		return nil, fmt.Errorf("error querying client: %w", err)
	}

	return &c, nil
}

// List returns all clients sorted by creation date descending.
func (r *ClientRepository) List() ([]domain.Client, error) {
	query := `
		SELECT id, name, phone, bike_details, notes, created_at, updated_at
		FROM clients
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error listing clients: %w", err)
	}
	defer rows.Close()

	var clients []domain.Client
	for rows.Next() {
		var c domain.Client
		if err := rows.Scan(&c.ID, &c.Name, &c.Phone, &c.BikeDetails, &c.Notes, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning client row: %w", err)
		}
		clients = append(clients, c)
	}

	return clients, rows.Err()
}

// Search searches clients matching by name, phone, or bike details.
func (r *ClientRepository) Search(term string) ([]domain.Client, error) {
	query := `
		SELECT id, name, phone, bike_details, notes, created_at, updated_at
		FROM clients
		WHERE name LIKE ? OR phone LIKE ? OR bike_details LIKE ?
		ORDER BY name ASC
	`
	wildcard := "%" + strings.TrimSpace(term) + "%"

	rows, err := r.db.Query(query, wildcard, wildcard, wildcard)
	if err != nil {
		return nil, fmt.Errorf("search clients failed: %w", err)
	}
	defer rows.Close()

	var clients []domain.Client
	for rows.Next() {
		var c domain.Client
		if err := rows.Scan(&c.ID, &c.Name, &c.Phone, &c.BikeDetails, &c.Notes, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning searched client: %w", err)
		}
		clients = append(clients, c)
	}

	return clients, rows.Err()
}
