package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderFinanceWidget() string {
	if m.MonthlySummary == nil {
		return ""
	}

	badge := "● PROFIT"
	badgeColor := lipgloss.Color("#98C379")
	if m.MonthlySummary.NetEarnings < 0 {
		badge = "▼ DEFICIT"
		badgeColor = lipgloss.Color("#E06C75")
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#61AFEF")).
		Padding(0, 1).
		MarginBottom(1)

	content := fmt.Sprintf(
		"Monthly Balance (%d/%d): Jobs Completed: %d  |  Total Income: +$%.2f  |  Total Expenses: -$%.2f  |  Net: $%.2f [%s]",
		m.MonthlySummary.Month, m.MonthlySummary.Year,
		m.MonthlySummary.ServicesCount,
		m.MonthlySummary.TotalIncome,
		m.MonthlySummary.TotalExpense,
		m.MonthlySummary.NetEarnings,
		lipgloss.NewStyle().Foreground(badgeColor).Bold(true).Render(badge),
	)

	return style.Render(content)
}
