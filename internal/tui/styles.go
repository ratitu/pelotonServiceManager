package tui

import (
	"github.com/charmbracelet/lipgloss"
	"peloton/internal/domain"
)

// Palette colors
var (
	ColorPending         = lipgloss.Color("#E5C07B") // Warm Amber
	ColorInProgress      = lipgloss.Color("#61AFEF") // Bright Cyan/Blue
	ColorWaitingForParts = lipgloss.Color("#D19A66") // Rust Orange
	ColorReadyForPickup  = lipgloss.Color("#98C379") // Lime Green
	ColorCompleted       = lipgloss.Color("#5C6370") // Slate Grey
	ColorAccent          = lipgloss.Color("#C678DD") // Purple
	ColorSelected        = lipgloss.Color("#E06C75") // Coral Red / Selection
	ColorCardBg          = lipgloss.Color("#21252B") // Dark Slate
	ColorActiveCardBg    = lipgloss.Color("#282C34") // Active Slate
	ColorBorder          = lipgloss.Color("#3E4451") // Muted Border
	ColorHighlightBorder = lipgloss.Color("#98C379") // Active Column Border
)

// Base UI Styles
var (
	StyleAppHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#282C34")).
			Padding(0, 1).
			MarginBottom(1)

	StyleStatusBar = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ABB2BF")).
			Background(lipgloss.Color("#1E2227")).
			Padding(0, 1).
			MarginTop(1)

	StyleColumn = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(0, 1)

	StyleActiveColumn = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorHighlightBorder).
				Padding(0, 1)

	StyleColumnHeader = lipgloss.NewStyle().
				Bold(true).
				Align(lipgloss.Center).
				MarginBottom(1)

	StyleCard = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorBorder).
			Background(ColorCardBg).
			Padding(0, 1).
			MarginBottom(1)

	StyleSelectedCard = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				BorderForeground(ColorSelected).
				Background(ColorActiveCardBg).
				Padding(0, 1).
				MarginBottom(1)

	StyleClientName = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF"))

	StylePriorityHigh = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#E06C75")).
				Bold(true)

	StylePriorityNormal = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ABB2BF"))

	StylePriorityUrgent = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#BE5046")).
				Bold(true).
				Blink(true)
)

// GetStatusColor maps ticket status to its designated theme color.
func GetStatusColor(status domain.TicketStatus) lipgloss.Color {
	switch status {
	case domain.StatusPending:
		return ColorPending
	case domain.StatusInProgress:
		return ColorInProgress
	case domain.StatusWaitingParts:
		return ColorWaitingForParts
	case domain.StatusReadyPickup:
		return ColorReadyForPickup
	case domain.StatusCompleted:
		return ColorCompleted
	default:
		return lipgloss.Color("#ABB2BF")
	}
}
