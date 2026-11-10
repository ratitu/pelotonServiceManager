package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles incoming messages and keyboard events.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// If modal is open, delegate to modal keyhandler
		if m.ShowModal {
			return m.updateModal(msg)
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "left", "h":
			if m.ActiveColIdx > 0 {
				m.ActiveColIdx--
			} else {
				m.ActiveColIdx = len(m.Columns) - 1 // Wrap around
			}
			return m, nil

		case "right", "l":
			if m.ActiveColIdx < len(m.Columns)-1 {
				m.ActiveColIdx++
			} else {
				m.ActiveColIdx = 0 // Wrap around
			}
			return m, nil

		case "tab":
			m.ActiveColIdx = (m.ActiveColIdx + 1) % len(m.Columns)
			return m, nil

		case "shift+tab":
			m.ActiveColIdx = (m.ActiveColIdx - 1 + len(m.Columns)) % len(m.Columns)
			return m, nil

		case "up", "k":
			col := &m.Columns[m.ActiveColIdx]
			if col.CursorIndex > 0 {
				col.CursorIndex--
			}
			return m, nil

		case "down", "j":
			col := &m.Columns[m.ActiveColIdx]
			if col.CursorIndex < len(col.Tickets)-1 {
				col.CursorIndex++
			}
			return m, nil

		case "r":
			m.ReloadTickets()
			m.ReloadFinance()
			return m, nil

		case "m", "f":
			m.ShowFinance = !m.ShowFinance
			return m, nil

		case "enter":
			if t := m.SelectedTicket(); t != nil {
				m.ActiveTicket = t
				m.ShowModal = true
				m.ModalMessage = ""
			}
			return m, nil
		}
	}

	return m, nil
}
