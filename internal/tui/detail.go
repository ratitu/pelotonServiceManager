package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) renderDetailModal() string {
	// Stub to be implemented fully in Wave 2
	return "Ticket Detail Modal"
}

func (m Model) updateModal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Stub to be implemented fully in Wave 2
	if msg.String() == "esc" || msg.String() == "q" {
		m.ShowModal = false
	}
	return m, nil
}
