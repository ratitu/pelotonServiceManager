package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"peloton/internal/domain"
)

// View renders the TUI layout.
func (m Model) View() string {
	if m.Width == 0 {
		return "Initializing Peloton Workbench..."
	}

	var sections []string

	// 1. App Header
	header := StyleAppHeader.Width(m.Width - 2).Render(
		"🚴 PELOTON WORKSHOP MANAGER | Desktop Workbench Dashboard",
	)
	sections = append(sections, header)

	// 2. Financial Summary Widget (if enabled)
	if m.ShowFinance && m.MonthlySummary != nil {
		finWidget := m.renderFinanceWidget()
		sections = append(sections, finWidget)
	}

	// 3. Main Board: 5 Columns
	board := m.renderBoard()
	sections = append(sections, board)

	// 4. Modal Overlay (if inspection active)
	if m.ShowModal && m.ActiveTicket != nil {
		modalView := m.renderDetailModal()
		// If modal is active, display modal centered
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, modalView)
	}

	// 5. Status / Keybinding Legend
	legend := StyleStatusBar.Width(m.Width - 2).Render(
		" [←/→/h/l] Switch Column  [↑/↓/j/k] Select Ticket  [Enter] Inspect/Update  [m] Toggle Finances  [r] Refresh  [q] Quit",
	)
	sections = append(sections, legend)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) renderBoard() string {
	numCols := len(m.Columns)
	// Calculate column width dynamically based on terminal width
	colWidth := (m.Width - (numCols * 3) - 2) / numCols
	if colWidth < 22 {
		colWidth = 22
	}

	colViews := make([]string, numCols)

	for i, col := range m.Columns {
		isActive := i == m.ActiveColIdx
		statusColor := GetStatusColor(col.Status)

		// Column Header
		headerTitle := fmt.Sprintf("%s (%d)", strings.ToUpper(string(col.Status)), len(col.Tickets))
		headerStyle := StyleColumnHeader.Width(colWidth).Foreground(statusColor)
		colHeader := headerStyle.Render(headerTitle)

		// Column Body (Ticket Cards)
		var cards []string
		cards = append(cards, colHeader)

		if len(col.Tickets) == 0 {
			emptyNotice := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#5C6370")).
				Italic(true).
				Align(lipgloss.Center).
				Width(colWidth).
				Render("No tickets")
			cards = append(cards, emptyNotice)
		} else {
			for tIdx, ticket := range col.Tickets {
				isSelected := isActive && tIdx == col.CursorIndex
				cardStr := m.renderTicketCard(ticket, isSelected, colWidth-4)
				cards = append(cards, cardStr)
			}
		}

		bodyContent := strings.Join(cards, "\n")
		colStyle := StyleColumn
		if isActive {
			colStyle = StyleActiveColumn
		}

		colViews[i] = colStyle.Width(colWidth).Render(bodyContent)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, colViews...)
}

func (m Model) renderTicketCard(t domain.ServiceTicket, isSelected bool, width int) string {
	style := StyleCard
	if isSelected {
		style = StyleSelectedCard
	}

	// Priority badge
	var priorityStyle lipgloss.Style
	switch strings.ToLower(t.Priority) {
	case "urgent":
		priorityStyle = StylePriorityUrgent
	case "high":
		priorityStyle = StylePriorityHigh
	default:
		priorityStyle = StylePriorityNormal
	}
	priorityBadge := priorityStyle.Render("[" + t.Priority + "]")

	// Cost string
	cost := fmt.Sprintf("$%.2f", t.EstimatedCost)
	if t.FinalCost > 0 {
		cost = fmt.Sprintf("$%.2f (fin)", t.FinalCost)
	}

	// Client & Bike
	clientInfo := StyleClientName.Render(fmt.Sprintf("#%d %s", t.ID, t.ClientName))

	// Truncate description for card preview
	desc := t.Description
	if len(desc) > width-2 && width > 4 {
		desc = desc[:width-5] + "..."
	}

	content := fmt.Sprintf("%s  %s\n%s\nCost: %s", clientInfo, priorityBadge, desc, cost)
	return style.Width(width).Render(content)
}
