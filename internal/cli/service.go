package cli

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"

	"peloton/internal/domain"
)

// RunService handles 'peloton service' subcommands.
func (app *App) RunService(args []string) int {
	if len(args) == 0 {
		app.printServiceHelp()
		return 0
	}

	switch args[0] {
	case "add":
		return app.serviceAdd(args[1:])
	case "list":
		return app.serviceList(args[1:])
	case "status":
		return app.serviceStatus(args[1:])
	case "--help", "-h", "help":
		app.printServiceHelp()
		return 0
	default:
		fmt.Fprintf(app.Stderr, "Unknown service command: '%s'. Run 'peloton service --help'.\n", args[0])
		return 1
	}
}

func (app *App) printServiceHelp() {
	help := `Usage: peloton service <subcommand> [flags]

Subcommands:
  add       Create a new service repair ticket
            --client <id>       Client ID (required)
            --desc <string>     Description of repair work (required)
            --priority <string> Normal | High | Urgent (default: Normal)
            --est <float>       Estimated cost (default: 0.0)

  list      List repair tickets in the workshop queue
            --status <status>   Filter by status (Pending, In Progress, Waiting for Parts, Ready for Pickup, Completed)

  status    Update ticket status and optionally set final cost
            Usage: peloton service status <ticket_id> <new_status> [--cost <final_cost>]
            Valid statuses:
              - "Pending"
              - "In Progress"
              - "Waiting for Parts"
              - "Ready for Pickup"
              - "Completed"
`
	fmt.Fprint(app.Stdout, help)
}

func (app *App) serviceAdd(args []string) int {
	fs := flag.NewFlagSet("service add", flag.ContinueOnError)
	fs.SetOutput(app.Stderr)

	clientID := fs.Int64("client", 0, "Client ID")
	desc := fs.String("desc", "", "Description of repair")
	priority := fs.String("priority", "Normal", "Priority level")
	estCost := fs.Float64("est", 0.0, "Estimated cost")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if *clientID <= 0 {
		fmt.Fprintln(app.Stderr, "Error: --client <id> is required.")
		return 1
	}
	if strings.TrimSpace(*desc) == "" {
		fmt.Fprintln(app.Stderr, "Error: --desc <text> is required.")
		return 1
	}

	// Verify client exists
	client, err := app.ClientRepo.GetByID(*clientID)
	if err != nil {
		fmt.Fprintf(app.Stderr, "Error: client #%d does not exist.\n", *clientID)
		return 1
	}

	ticket := &domain.ServiceTicket{
		ClientID:      *clientID,
		Description:   strings.TrimSpace(*desc),
		Status:        domain.StatusPending,
		Priority:      strings.TrimSpace(*priority),
		EstimatedCost: *estCost,
	}

	if err := app.TicketRepo.Create(ticket); err != nil {
		fmt.Fprintf(app.Stderr, "Error creating service ticket: %v\n", err)
		return 1
	}

	fmt.Fprintf(app.Stdout, "✓ Service Ticket #%d created for client %s (%s)!\n", ticket.ID, client.Name, client.Phone)
	fmt.Fprintf(app.Stdout, "  Status: %s | Priority: %s | Est: $%.2f\n", ticket.Status, ticket.Priority, ticket.EstimatedCost)
	return 0
}

func (app *App) serviceList(args []string) int {
	fs := flag.NewFlagSet("service list", flag.ContinueOnError)
	fs.SetOutput(app.Stderr)

	statusFilter := fs.String("status", "", "Filter by status")
	if err := fs.Parse(args); err != nil {
		return 1
	}

	var filter domain.TicketStatus
	if *statusFilter != "" {
		filter = domain.TicketStatus(strings.TrimSpace(*statusFilter))
		if !domain.IsValidStatus(filter) {
			fmt.Fprintf(app.Stderr, "Warning: '%s' is not a standard status. Filtering as provided.\n", *statusFilter)
		}
	}

	tickets, err := app.TicketRepo.List(filter)
	if err != nil {
		fmt.Fprintf(app.Stderr, "Error listing tickets: %v\n", err)
		return 1
	}

	if len(tickets) == 0 {
		if filter != "" {
			fmt.Fprintf(app.Stdout, "No service tickets found with status '%s'.\n", filter)
		} else {
			fmt.Fprintln(app.Stdout, "Workshop queue is empty! Use 'peloton service add' to log a job.")
		}
		return 0
	}

	app.renderTicketTable(tickets)
	return 0
}

func (app *App) serviceStatus(args []string) int {
	if len(args) < 2 {
		fmt.Fprintln(app.Stderr, "Error: ticket ID and new status are required.")
		fmt.Fprintln(app.Stderr, "Usage: peloton service status <ticket_id> <new_status> [--cost <final_cost>]")
		return 1
	}

	ticketID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		fmt.Fprintf(app.Stderr, "Invalid ticket ID: '%s'\n", args[0])
		return 1
	}

	newStatusStr := args[1]
	// Handle multi-word status like "Waiting for Parts", "In Progress", "Ready for Pickup"
	// Also parse optional --cost flag from remaining args
	var finalCost *float64
	var statusWords []string

	for i := 1; i < len(args); i++ {
		arg := args[i]
		if arg == "--cost" {
			if i+1 < len(args) {
				costVal, err := strconv.ParseFloat(args[i+1], 64)
				if err == nil {
					finalCost = &costVal
				}
				i++
				continue
			}
		} else if strings.HasPrefix(arg, "--cost=") {
			costVal, err := strconv.ParseFloat(strings.TrimPrefix(arg, "--cost="), 64)
			if err == nil {
				finalCost = &costVal
			}
			continue
		} else {
			statusWords = append(statusWords, arg)
		}
	}

	newStatusStr = strings.Join(statusWords, " ")

	// Normalize case-insensitive matching for convenience
	matchedStatus := domain.TicketStatus(newStatusStr)
	for _, s := range domain.AllTicketStatuses {
		if strings.EqualFold(string(s), newStatusStr) {
			matchedStatus = s
			break
		}
	}

	if !domain.IsValidStatus(matchedStatus) {
		fmt.Fprintf(app.Stderr, "Error: invalid status '%s'. Allowed values:\n", newStatusStr)
		for _, s := range domain.AllTicketStatuses {
			fmt.Fprintf(app.Stderr, "  - %s\n", s)
		}
		return 1
	}

	if err := app.TicketRepo.UpdateStatus(ticketID, matchedStatus, finalCost); err != nil {
		fmt.Fprintf(app.Stderr, "Error updating ticket #%d: %v\n", ticketID, err)
		return 1
	}

	fmt.Fprintf(app.Stdout, "✓ Service Ticket #%d updated to status '%s'!\n", ticketID, matchedStatus)
	if finalCost != nil {
		fmt.Fprintf(app.Stdout, "  Final labor/repair cost recorded: $%.2f\n", *finalCost)
	}
	return 0
}

func (app *App) renderTicketTable(tickets []domain.ServiceTicket) {
	w := tabwriter.NewWriter(app.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "TICKET\tCLIENT\tPHONE\tSTATUS\tPRIORITY\tEST/FINAL\tDESCRIPTION")
	fmt.Fprintln(w, "------\t------\t-----\t------\t--------\t---------\t-----------")

	for _, t := range tickets {
		costStr := fmt.Sprintf("$%.2f", t.EstimatedCost)
		if t.FinalCost > 0 {
			costStr = fmt.Sprintf("$%.2f (fin)", t.FinalCost)
		}
		desc := t.Description
		if len(desc) > 35 {
			desc = desc[:32] + "..."
		}
		fmt.Fprintf(w, "#%d\t%s\t%s\t%s\t%s\t%s\t%s\n",
			t.ID, t.ClientName, t.ClientPhone, t.Status, t.Priority, costStr, desc,
		)
	}
	w.Flush()
}
