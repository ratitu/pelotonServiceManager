package sqlite

const schemaMigrations = `
CREATE TABLE IF NOT EXISTS clients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    phone TEXT NOT NULL,
    bike_details TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_clients_phone ON clients(phone);
CREATE INDEX IF NOT EXISTS idx_clients_name ON clients(name);

CREATE TABLE IF NOT EXISTS service_tickets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id INTEGER NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    status TEXT NOT NULL,
    priority TEXT NOT NULL DEFAULT 'Normal',
    estimated_cost REAL NOT NULL DEFAULT 0.0,
    final_cost REAL NOT NULL DEFAULT 0.0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    completed_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_service_tickets_status ON service_tickets(status);
CREATE INDEX IF NOT EXISTS idx_service_tickets_client_id ON service_tickets(client_id);
CREATE INDEX IF NOT EXISTS idx_service_tickets_created_at ON service_tickets(created_at);

CREATE TABLE IF NOT EXISTS financial_transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ticket_id INTEGER REFERENCES service_tickets(id) ON DELETE SET NULL,
    type TEXT NOT NULL, -- 'INCOME' or 'EXPENSE'
    category TEXT NOT NULL,
    amount REAL NOT NULL,
    notes TEXT NOT NULL DEFAULT '',
    date DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_financial_transactions_date ON financial_transactions(date);
CREATE INDEX IF NOT EXISTS idx_financial_transactions_type ON financial_transactions(type);
`

// Migrate runs the schema definition against the database.
func (db *DB) Migrate() error {
	_, err := db.Exec(schemaMigrations)
	return err
}
