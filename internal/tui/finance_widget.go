package tui

import (
	"fmt"
	"time"

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
	} else if m.MonthlySummary.NetEarnings == 0 {
		badge = "■ BREAK-EVEN"
		badgeColor = lipgloss.Color("#ABB2BF")
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#61AFEF")).
		Padding(0, 1).
		MarginBottom(1)

	monthName := time.Month(m.MonthlySummary.Month).String()

	content := fmt.Sprintf(
		"📊 %s %d EARNINGS | Completed Jobs: %d  |  Income: +$%.2f  |  Expenses: -$%.2f  |  Net: $%.2f [%s]",
		monthName, m.MonthlySummary.Year,
		m.MonthlySummary.ServicesCount,
		m.MonthlySummary.TotalIncome,
		m.MonthlySummary.TotalExpense,
		m.MonthlySummary.NetEarnings,
		lipgloss.NewStyle().Foreground(badgeColor).Bold(true).Render(badge),
	)

	return style.Render(content)
}
