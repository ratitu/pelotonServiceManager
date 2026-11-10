package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"peloton/internal/domain"
)

var (
	StyleModalBorder = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(lipgloss.Color("#61AFEF")).
				Background(lipgloss.Color("#1E2227")).
				Padding(1, 2)

	StyleModalTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#61AFEF")).
			Padding(0, 1)

	StyleModalLabel = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#E5C07B"))

	StyleModalValue = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ABB2BF"))

	StyleActionKey = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#98C379"))

	StyleModalFooter = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#5C6370")).
				Italic(true).
				MarginTop(1)
)

func (m Model) renderDetailModal() string {
	t := m.ActiveTicket
	if t == nil {
		return ""
	}

	modalWidth := 65
	if m.Width > 0 && m.Width-10 < modalWidth {
		modalWidth = m.Width - 10
	}

	title := StyleModalTitle.Render(fmt.Sprintf(" SERVICE TICKET #%d ", t.ID))

	statusColor := GetStatusColor(t.Status)
	statusBadge := lipgloss.NewStyle().
		Bold(true).
		Foreground(statusColor).
		Render(string(t.Status))

	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("%s %s", StyleModalLabel.Render("Client Name: "), StyleModalValue.Render(t.ClientName)))
	lines = append(lines, fmt.Sprintf("%s %s", StyleModalLabel.Render("Phone:       "), StyleModalValue.Render(t.ClientPhone)))
	lines = append(lines, fmt.Sprintf("%s %s", StyleModalLabel.Render("Status:      "), statusBadge))
	lines = append(lines, fmt.Sprintf("%s %s", StyleModalLabel.Render("Priority:    "), StyleModalValue.Render(t.Priority)))
	lines = append(lines, fmt.Sprintf("%s $%.2f", StyleModalLabel.Render("Estimated:   "), t.EstimatedCost))

	finalCostStr := "-"
	if t.FinalCost > 0 {
		finalCostStr = fmt.Sprintf("$%.2f", t.FinalCost)
	}
	lines = append(lines, fmt.Sprintf("%s %s", StyleModalLabel.Render("Final Cost:  "), StyleModalValue.Render(finalCostStr)))
	lines = append(lines, "")

	lines = append(lines, StyleModalLabel.Render("Problem Description:"))
	descBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#3E4451")).
		Width(modalWidth - 6).
		Padding(0, 1).
		Render(t.Description)
	lines = append(lines, descBox)
	lines = append(lines, "")

	if m.ModalMessage != "" {
		msgStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#98C379")).Bold(true)
		lines = append(lines, msgStyle.Render(m.ModalMessage))
		lines = append(lines, "")
	}

	// Hotkey Actions for in-place status updates
	actions := fmt.Sprintf(
		"Advance Status:  %s Pending  %s In Progress  %s Waiting Parts  %s Ready  %s Complete",
		StyleActionKey.Render("[1]"),
		StyleActionKey.Render("[2]"),
		StyleActionKey.Render("[3]"),
		StyleActionKey.Render("[4]"),
		StyleActionKey.Render("[5]"),
	)
	lines = append(lines, actions)

	footer := StyleModalFooter.Render("Press [Esc] or [q] to return to board")
	lines = append(lines, footer)

	content := strings.Join(lines, "\n")
	return StyleModalBorder.Width(modalWidth).Render(content)
}

func (m Model) updateModal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q", "enter":
		m.ShowModal = false
		m.ActiveTicket = nil
		m.ModalMessage = ""
		return m, nil

	case "1":
		m.changeActiveTicketStatus(domain.StatusPending, nil)
	case "2":
		m.changeActiveTicketStatus(domain.StatusInProgress, nil)
	case "3":
		m.changeActiveTicketStatus(domain.StatusWaitingParts, nil)
	case "4":
		m.changeActiveTicketStatus(domain.StatusReadyPickup, nil)
	case "5":
		// When completing ticket, if final cost is not set, use estimated cost as default
		finalCost := m.ActiveTicket.EstimatedCost
		if m.ActiveTicket.FinalCost > 0 {
			finalCost = m.ActiveTicket.FinalCost
		}
		m.changeActiveTicketStatus(domain.StatusCompleted, &finalCost)
	}

	return m, nil
}

func (m *Model) changeActiveTicketStatus(newStatus domain.TicketStatus, finalCost *float64) {
	if m.ActiveTicket == nil || m.TicketRepo == nil {
		return
	}

	err := m.TicketRepo.UpdateStatus(m.ActiveTicket.ID, newStatus, finalCost)
	if err != nil {
		m.ModalMessage = fmt.Sprintf("Error updating status: %v", err)
		return
	}

	m.ActiveTicket.Status = newStatus
	if finalCost != nil {
		m.ActiveTicket.FinalCost = *finalCost
	}

	// Reload board & finance data immediately
	m.ReloadTickets()
	m.ReloadFinance()

	m.ModalMessage = fmt.Sprintf("✓ Status updated to '%s'!", newStatus)
}
