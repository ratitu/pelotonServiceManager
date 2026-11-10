package sqlite

import (
	"fmt"
	"time"

	"peloton/internal/domain"
)

// FinanceRepository handles storage of income and expense transactions.
type FinanceRepository struct {
	db *DB
}

// NewFinanceRepository constructs a new FinanceRepository.
func NewFinanceRepository(db *DB) *FinanceRepository {
	return &FinanceRepository{db: db}
}

// Create inserts a new financial transaction.
func (r *FinanceRepository) Create(t *domain.FinancialTransaction) error {
	if t.Amount <= 0 {
		return fmt.Errorf("transaction amount must be positive, got %.2f", t.Amount)
	}
	if t.Type != domain.TypeIncome && t.Type != domain.TypeExpense {
		return fmt.Errorf("invalid transaction type: %s", t.Type)
	}
	if t.Date.IsZero() {
		t.Date = time.Now()
	}

	query := `
		INSERT INTO financial_transactions (ticket_id, type, category, amount, notes, date)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	res, err := r.db.Exec(query, t.TicketID, string(t.Type), t.Category, t.Amount, t.Notes, t.Date)
	if err != nil {
		return fmt.Errorf("failed to insert financial transaction: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get transaction id: %w", err)
	}

	t.ID = id
	return nil
}

// ListTransactions retrieves transactions within a date range.
func (r *FinanceRepository) ListTransactions(startDate, endDate time.Time) ([]domain.FinancialTransaction, error) {
	query := `
		SELECT id, ticket_id, type, category, amount, notes, date
		FROM financial_transactions
		WHERE date >= ? AND date <= ?
		ORDER BY date DESC
	`

	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error querying transactions: %w", err)
	}
	defer rows.Close()

	var list []domain.FinancialTransaction
	for rows.Next() {
		var t domain.FinancialTransaction
		var typeStr string
		err := rows.Scan(&t.ID, &t.TicketID, &typeStr, &t.Category, &t.Amount, &t.Notes, &t.Date)
		if err != nil {
			return nil, fmt.Errorf("error scanning transaction row: %w", err)
		}
		t.Type = domain.TransactionType(typeStr)
		list = append(list, t)
	}

	return list, rows.Err()
}

// GetMonthlySummary calculates total income, total expenses, and completed services for a given month.
func (r *FinanceRepository) GetMonthlySummary(year int, month int) (*domain.MonthlyBalance, error) {
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	summary := &domain.MonthlyBalance{
		Year:  year,
		Month: month,
	}

	// 1. Calculate Income and Expense totals
	financeQuery := `
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'INCOME' THEN amount ELSE 0 END), 0.0) as total_income,
			COALESCE(SUM(CASE WHEN type = 'EXPENSE' THEN amount ELSE 0 END), 0.0) as total_expense
		FROM financial_transactions
		WHERE date >= ? AND date <= ?
	`

	err := r.db.QueryRow(financeQuery, startOfMonth, endOfMonth).Scan(&summary.TotalIncome, &summary.TotalExpense)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate monthly financial transactions: %w", err)
	}

	summary.NetEarnings = summary.TotalIncome - summary.TotalExpense

	// 2. Count completed services in that month
	ticketQuery := `
		SELECT COUNT(*)
		FROM service_tickets
		WHERE status = 'Completed' AND completed_at >= ? AND completed_at <= ?
	`

	err = r.db.QueryRow(ticketQuery, startOfMonth, endOfMonth).Scan(&summary.ServicesCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count completed monthly services: %w", err)
	}

	return summary, nil
}
