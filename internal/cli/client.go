package cli

import (
	"flag"
	"fmt"
	"strings"
	"text/tabwriter"

	"peloton/internal/domain"
)

// RunClient handles 'peloton client' subcommands.
func (app *App) RunClient(args []string) int {
	if len(args) == 0 {
		app.printClientHelp()
		return 0
	}

	switch args[0] {
	case "add":
		return app.clientAdd(args[1:])
	case "list":
		return app.clientList(args[1:])
	case "search":
		return app.clientSearch(args[1:])
	case "--help", "-h", "help":
		app.printClientHelp()
		return 0
	default:
		fmt.Fprintf(app.Stderr, "Unknown client command: '%s'. Run 'peloton client --help'.\n", args[0])
		return 1
	}
}

func (app *App) printClientHelp() {
	help := `Usage: peloton client <subcommand> [flags]

Subcommands:
  add       Register a new client
            --name <string>     Client name (required)
            --phone <string>    Client phone number (required)
            --bike <string>     Bicycle make/model/color (optional)
            --notes <string>    Customer or service notes (optional)

  list      List all clients in the system

  search    Search clients by query string
            Usage: peloton client search <query>
`
	fmt.Fprint(app.Stdout, help)
}

func (app *App) clientAdd(args []string) int {
	fs := flag.NewFlagSet("client add", flag.ContinueOnError)
	fs.SetOutput(app.Stderr)

	name := fs.String("name", "", "Client name")
	phone := fs.String("phone", "", "Phone number")
	bike := fs.String("bike", "", "Bike details")
	notes := fs.String("notes", "", "Notes")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	client := &domain.Client{
		Name:        strings.TrimSpace(*name),
		Phone:       strings.TrimSpace(*phone),
		BikeDetails: strings.TrimSpace(*bike),
		Notes:       strings.TrimSpace(*notes),
	}

	if err := app.ClientRepo.Create(client); err != nil {
		fmt.Fprintf(app.Stderr, "Error creating client: %v\n", err)
		return 1
	}

	fmt.Fprintf(app.Stdout, "✓ Client registered successfully! (ID: #%d, Name: %s)\n", client.ID, client.Name)
	return 0
}

func (app *App) clientList(args []string) int {
	clients, err := app.ClientRepo.List()
	if err != nil {
		fmt.Fprintf(app.Stderr, "Error retrieving clients: %v\n", err)
		return 1
	}

	if len(clients) == 0 {
		fmt.Fprintln(app.Stdout, "No clients registered yet. Use 'peloton client add' to create one.")
		return 0
	}

	app.renderClientTable(clients)
	return 0
}

func (app *App) clientSearch(args []string) int {
	if len(args) == 0 || strings.TrimSpace(args[0]) == "" {
		fmt.Fprintln(app.Stderr, "Error: search query cannot be empty. Usage: peloton client search <query>")
		return 1
	}

	query := strings.Join(args, " ")
	clients, err := app.ClientRepo.Search(query)
	if err != nil {
		fmt.Fprintf(app.Stderr, "Error searching clients: %v\n", err)
		return 1
	}

	if len(clients) == 0 {
		fmt.Fprintf(app.Stdout, "No clients found matching '%s'.\n", query)
		return 0
	}

	app.renderClientTable(clients)
	return 0
}

func (app *App) renderClientTable(clients []domain.Client) {
	w := tabwriter.NewWriter(app.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tPHONE\tBIKE DETAILS\tREGISTERED")
	fmt.Fprintln(w, "--\t----\t-----\t------------\t----------")

	for _, c := range clients {
		bike := c.BikeDetails
		if bike == "" {
			bike = "-"
		}
		dateStr := c.CreatedAt.Format("2006-01-02 15:04")
		fmt.Fprintf(w, "#%d\t%s\t%s\t%s\t%s\n", c.ID, c.Name, c.Phone, bike, dateStr)
	}
	w.Flush()
}
