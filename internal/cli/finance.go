package cli

import (
	"fmt"
)

// RunFinance handles 'peloton finance' subcommands.
func (app *App) RunFinance(args []string) int {
	if len(args) == 0 {
		app.printFinanceHelp()
		return 0
	}

	switch args[0] {
	case "income":
		return app.financeIncome(args[1:])
	case "expense":
		return app.financeExpense(args[1:])
	case "summary":
		return app.financeSummary(args[1:])
	case "list":
		return app.financeList(args[1:])
	case "--help", "-h", "help":
		app.printFinanceHelp()
		return 0
	default:
		fmt.Fprintf(app.Stderr, "Unknown finance command: '%s'. Run 'peloton finance --help'.\n", args[0])
		return 1
	}
}

func (app *App) printFinanceHelp() {
	help := `Usage: peloton finance <subcommand> [flags]

Subcommands:
  income    Log revenue or client service payment
  expense   Log replacement parts or shop supplies purchase
  summary   Display monthly financial summary and net profit
  list      List individual logged transactions
`
	fmt.Fprint(app.Stdout, help)
}

func (app *App) financeIncome(args []string) int {
	// Stub to be completed in Wave 2
	return 0
}

func (app *App) financeExpense(args []string) int {
	// Stub to be completed in Wave 2
	return 0
}

func (app *App) financeSummary(args []string) int {
	// Stub to be completed in Wave 2
	return 0
}

func (app *App) financeList(args []string) int {
	// Stub to be completed in Wave 2
	return 0
}
