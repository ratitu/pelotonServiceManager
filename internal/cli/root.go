package cli

import (
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"peloton/internal/config"
	"peloton/internal/store/sqlite"
	"peloton/internal/tui"
)

// App manages dependency injection and command execution for the CLI.
type App struct {
	Config      *config.Config
	DB          *sqlite.DB
	ClientRepo  *sqlite.ClientRepository
	TicketRepo  *sqlite.TicketRepository
	FinanceRepo *sqlite.FinanceRepository
	Stdout      io.Writer
	Stderr      io.Writer
}

// NewApp instantiates an App with given DB path and output writers.
func NewApp(dbPath string, stdout, stderr io.Writer) (*App, error) {
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}

	cfg, err := config.NewConfig(dbPath, 8080, "0.0.0.0")
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	db, err := sqlite.OpenDB(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("database initialization error: %w", err)
	}

	return &App{
		Config:      cfg,
		DB:          db,
		ClientRepo:  sqlite.NewClientRepository(db),
		TicketRepo:  sqlite.NewTicketRepository(db),
		FinanceRepo: sqlite.NewFinanceRepository(db),
		Stdout:      stdout,
		Stderr:      stderr,
	}, nil
}

// Close closes the underlying database connection.
func (app *App) Close() error {
	if app.DB != nil {
		return app.DB.Close()
	}
	return nil
}

// PrintHelp outputs the main usage manual.
func (app *App) PrintHelp() {
	help := `Peloton Service Manager - Bicycle Workshop Management System

USAGE:
  peloton [global options] <command> <subcommand> [flags] [arguments]

GLOBAL OPTIONS:
  --db <path>      Override SQLite database file path (defaults to XDG path or PELOTON_DB_PATH)
  --help, -h       Show help information
  --version, -v    Show software version

COMMANDS:
  tui              Launch interactive desktop workbench dashboard (default when no args)
  client           Manage bicycle owners and client profiles
    add            Register a new client
    list           List all clients in the system
    search         Search clients by name, phone, or bike

  service          Manage repair queue and service tickets
    add            Create a repair service ticket for a client
    list           List active or past service tickets
    status         Advance or update a service ticket's status

  finance          Manage workshop finances and cashflow
    income         Log income payment (cash in)
    expense        Log parts or supplies purchase (cash out)
    summary        View monthly balance and earnings summary
    list           List individual financial transaction records

Run 'peloton <command> --help' for details on a specific command.
`
	fmt.Fprint(app.Stdout, help)
}

// Run parses arguments and routes to the appropriate subcommand handler.
func (app *App) Run(args []string) int {
	if len(args) == 0 {
		app.PrintHelp()
		return 0
	}

	switch args[0] {
	case "--help", "-h", "help":
		app.PrintHelp()
		return 0
	case "--version", "-v", "version":
		fmt.Fprintln(app.Stdout, "Peloton Service Manager v1.0.0 (Debian/Go)")
		return 0
	case "tui":
		return app.RunTUI()
	case "client":
		return app.RunClient(args[1:])
	case "service":
		return app.RunService(args[1:])
	case "finance":
		return app.RunFinance(args[1:])
	default:
		fmt.Fprintf(app.Stderr, "Unknown command: '%s'. Run 'peloton --help' for usage.\n", args[0])
		return 1
	}
}

// RunTUI initializes and runs the interactive Charm Bubbletea TUI program.
func (app *App) RunTUI() int {
	model := tui.NewModel(app.TicketRepo, app.FinanceRepo, app.ClientRepo)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(app.Stderr, "Error running TUI dashboard: %v\n", err)
		return 1
	}
	return 0
}
