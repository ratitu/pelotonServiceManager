package cli

import (
	"flag"
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	"peloton/internal/domain"
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
  income     Log revenue or customer service payment
             --amount <float>    Payment amount (required, > 0)
             --category <string> Category (default: "Service")
             --ticket <id>       Associated service ticket ID (optional)
             --notes <string>    Payment notes (e.g. "Cash", "Card")
             --date <YYYY-MM-DD> Transaction date (default: today)

  expense    Log replacement parts or workshop supplies purchase
             --amount <float>    Expense amount (required, > 0)
             --category <string> Category (default: "Parts")
             --ticket <id>       Associated service ticket ID (optional)
             --notes <string>    Expense notes / invoice details
             --date <YYYY-MM-DD> Transaction date (default: today)

  summary    Display monthly cashflow balance and net earnings
             --month <1-12>      Month number (default: current month)
             --year <YYYY>       Year (default: current year)

  list       List individual financial transaction records
             --from <YYYY-MM-DD> Start date (default: start of current month)
             --to <YYYY-MM-DD>   End date (default: end of current month)
`
	fmt.Fprint(app.Stdout, help)
}

func (app *App) financeIncome(args []string) int {
	fs := flag.NewFlagSet("finance income", flag.ContinueOnError)
	fs.SetOutput(app.Stderr)

	amount := fs.Float64("amount", 0.0, "Income amount")
	category := fs.String("category", "Service", "Category")
	ticketID := fs.Int64("ticket", 0, "Ticket ID")
	notes := fs.String("notes", "", "Notes")
	dateStr := fs.String("date", "", "Date (YYYY-MM-DD)")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if *amount <= 0 {
		fmt.Fprintln(app.Stderr, "Error: --amount must be greater than 0.")
		return 1
	}

	txDate := time.Now()
	if *dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", *dateStr)
		if err != nil {
			fmt.Fprintf(app.Stderr, "Error: invalid date format '%s'. Use YYYY-MM-DD.\n", *dateStr)
			return 1
		}
		txDate = parsedDate
	}

	var tID *int64
	if *ticketID > 0 {
		// Verify ticket exists
		_, err := app.TicketRepo.GetByID(*ticketID)
		if err != nil {
			fmt.Fprintf(app.Stderr, "Warning: ticket #%d not found in database.\n", *ticketID)
		}
		tID = ticketID
	}

	tx := &domain.FinancialTransaction{
		TicketID: tID,
		Type:     domain.TypeIncome,
		Category: strings.TrimSpace(*category),
		Amount:   *amount,
		Notes:    strings.TrimSpace(*notes),
		Date:     txDate,
	}

	if err := app.FinanceRepo.Create(tx); err != nil {
		fmt.Fprintf(app.Stderr, "Error saving income transaction: %v\n", err)
		return 1
	}

	fmt.Fprintf(app.Stdout, "✓ Income logged: +$%.2f [%s] on %s\n",
		tx.Amount, tx.Category, tx.Date.Format("2006-01-02"),
	)
	if tx.TicketID != nil {
		fmt.Fprintf(app.Stdout, "  Associated with Service Ticket #%d\n", *tx.TicketID)
	}
	return 0
}

func (app *App) financeExpense(args []string) int {
	fs := flag.NewFlagSet("finance expense", flag.ContinueOnError)
	fs.SetOutput(app.Stderr)

	amount := fs.Float64("amount", 0.0, "Expense amount")
	category := fs.String("category", "Parts", "Category")
	ticketID := fs.Int64("ticket", 0, "Ticket ID")
	notes := fs.String("notes", "", "Notes")
	dateStr := fs.String("date", "", "Date (YYYY-MM-DD)")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if *amount <= 0 {
		fmt.Fprintln(app.Stderr, "Error: --amount must be greater than 0.")
		return 1
	}

	txDate := time.Now()
	if *dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", *dateStr)
		if err != nil {
			fmt.Fprintf(app.Stderr, "Error: invalid date format '%s'. Use YYYY-MM-DD.\n", *dateStr)
			return 1
		}
		txDate = parsedDate
	}

	var tID *int64
	if *ticketID > 0 {
		tID = ticketID
	}

	tx := &domain.FinancialTransaction{
		TicketID: tID,
		Type:     domain.TypeExpense,
		Category: strings.TrimSpace(*category),
		Amount:   *amount,
		Notes:    strings.TrimSpace(*notes),
		Date:     txDate,
	}

	if err := app.FinanceRepo.Create(tx); err != nil {
		fmt.Fprintf(app.Stderr, "Error saving expense transaction: %v\n", err)
		return 1
	}

	fmt.Fprintf(app.Stdout, "✓ Expense logged: -$%.2f [%s] on %s\n",
		tx.Amount, tx.Category, tx.Date.Format("2006-01-02"),
	)
	return 0
}

func (app *App) financeSummary(args []string) int {
	now := time.Now()
	fs := flag.NewFlagSet("finance summary", flag.ContinueOnError)
	fs.SetOutput(app.Stderr)

	month := fs.Int("month", int(now.Month()), "Month (1-12)")
	year := fs.Int("year", now.Year(), "Year")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if *month < 1 || *month > 12 {
		fmt.Fprintf(app.Stderr, "Error: invalid month %d. Must be between 1 and 12.\n", *month)
		return 1
	}

	summary, err := app.FinanceRepo.GetMonthlySummary(*year, *month)
	if err != nil {
		fmt.Fprintf(app.Stderr, "Error generating monthly summary: %v\n", err)
		return 1
	}

	monthName := time.Month(*month).String()

	statusIndicator := "● PROFIT"
	if summary.NetEarnings < 0 {
		statusIndicator = "▼ DEFICIT"
	} else if summary.NetEarnings == 0 {
		statusIndicator = "■ BREAK-EVEN"
	}

	card := fmt.Sprintf(`
==================================================
  PELOTON FINANCIAL SUMMARY: %s %d
==================================================
  Completed Services : %d jobs
  Total Income (+)   : $%.2f
  Total Expense (-)  : $%.2f
--------------------------------------------------
  Net Balance        : $%.2f (%s)
==================================================
`,
		strings.ToUpper(monthName), summary.Year,
		summary.ServicesCount,
		summary.TotalIncome,
		summary.TotalExpense,
		summary.NetEarnings,
		statusIndicator,
	)

	fmt.Fprint(app.Stdout, card)
	return 0
}

func (app *App) financeList(args []string) int {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	fs := flag.NewFlagSet("finance list", flag.ContinueOnError)
	fs.SetOutput(app.Stderr)

	fromStr := fs.String("from", startOfMonth.Format("2006-01-02"), "Start date (YYYY-MM-DD)")
	toStr := fs.String("to", endOfMonth.Format("2006-01-02"), "End date (YYYY-MM-DD)")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	fromDate, err := time.Parse("2006-01-02", *fromStr)
	if err != nil {
		fmt.Fprintf(app.Stderr, "Error: invalid --from date '%s'. Use YYYY-MM-DD.\n", *fromStr)
		return 1
	}

	toDate, err := time.Parse("2006-01-02", *toStr)
	if err != nil {
		fmt.Fprintf(app.Stderr, "Error: invalid --to date '%s'. Use YYYY-MM-DD.\n", *toStr)
		return 1
	}
	// Extend to end of the selected day
	toDate = toDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	txs, err := app.FinanceRepo.ListTransactions(fromDate, toDate)
	if err != nil {
		fmt.Fprintf(app.Stderr, "Error listing transactions: %v\n", err)
		return 1
	}

	if len(txs) == 0 {
		fmt.Fprintf(app.Stdout, "No transactions recorded between %s and %s.\n", *fromStr, *toStr)
		return 0
	}

	w := tabwriter.NewWriter(app.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "DATE\tTYPE\tAMOUNT\tCATEGORY\tTICKET\tNOTES")
	fmt.Fprintln(w, "----\t----\t------\t--------\t------\t-----")

	for _, tx := range txs {
		ticketStr := "-"
		if tx.TicketID != nil {
			ticketStr = fmt.Sprintf("#%d", *tx.TicketID)
		}
		amountStr := fmt.Sprintf("+$%.2f", tx.Amount)
		if tx.Type == domain.TypeExpense {
			amountStr = fmt.Sprintf("-$%.2f", tx.Amount)
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			tx.Date.Format("2006-01-02"),
			tx.Type,
			amountStr,
			tx.Category,
			ticketStr,
			tx.Notes,
		)
	}
	w.Flush()
	return 0
}
