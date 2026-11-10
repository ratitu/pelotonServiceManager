package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"peloton/internal/domain"
	"peloton/internal/store/sqlite"
)

// Column represents one status swimlane in the kanban board.
type Column struct {
	Status      domain.TicketStatus
	Tickets     []domain.ServiceTicket
	CursorIndex int
}

// Model stores the application state for the TUI workbench.
type Model struct {
	TicketRepo  *sqlite.TicketRepository
	FinanceRepo *sqlite.FinanceRepository
	ClientRepo  *sqlite.ClientRepository

	Columns      [5]Column
	ActiveColIdx int

	// Modal inspection state (detailed in Plan 03-02, prepared here)
	ActiveTicket *domain.ServiceTicket
	ShowModal    bool
	ModalMessage string

	// Financial quick summary state
	MonthlySummary *domain.MonthlyBalance
	ShowFinance    bool

	Width  int
	Height int
	Err    error
}

// NewModel initializes the TUI model and populates initial ticket queues.
func NewModel(ticketRepo *sqlite.TicketRepository, financeRepo *sqlite.FinanceRepository, clientRepo *sqlite.ClientRepository) Model {
	m := Model{
		TicketRepo:  ticketRepo,
		FinanceRepo: financeRepo,
		ClientRepo:  clientRepo,
		ShowFinance: true,
	}

	for i, status := range domain.AllTicketStatuses {
		m.Columns[i] = Column{
			Status:      status,
			Tickets:     []domain.ServiceTicket{},
			CursorIndex: 0,
		}
	}

	m.ReloadTickets()
	m.ReloadFinance()
	return m
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// ReloadTickets queries SQLite and groups tickets into their status columns.
func (m *Model) ReloadTickets() {
	if m.TicketRepo == nil {
		return
	}

	tickets, err := m.TicketRepo.List("")
	if err != nil {
		m.Err = err
		return
	}

	// Reset columns
	for i := range m.Columns {
		m.Columns[i].Tickets = []domain.ServiceTicket{}
	}

	for _, t := range tickets {
		for i, status := range domain.AllTicketStatuses {
			if t.Status == status {
				m.Columns[i].Tickets = append(m.Columns[i].Tickets, t)
				break
			}
		}
	}

	// Clamp cursor indices
	for i := range m.Columns {
		if len(m.Columns[i].Tickets) == 0 {
			m.Columns[i].CursorIndex = 0
		} else if m.Columns[i].CursorIndex >= len(m.Columns[i].Tickets) {
			m.Columns[i].CursorIndex = len(m.Columns[i].Tickets) - 1
		}
	}
}

// ReloadFinance updates the current month's financial metrics.
func (m *Model) ReloadFinance() {
	if m.FinanceRepo == nil {
		return
	}
	// Fetch monthly summary for today
	summary, err := m.FinanceRepo.GetMonthlySummary(0, 0)
	if err == nil {
		m.MonthlySummary = summary
	}
}

// SelectedTicket returns the currently highlighted ticket in the active column, or nil if column is empty.
func (m *Model) SelectedTicket() *domain.ServiceTicket {
	activeCol := &m.Columns[m.ActiveColIdx]
	if len(activeCol.Tickets) == 0 {
		return nil
	}
	if activeCol.CursorIndex >= len(activeCol.Tickets) {
		activeCol.CursorIndex = len(activeCol.Tickets) - 1
	}
	return &activeCol.Tickets[activeCol.CursorIndex]
}
