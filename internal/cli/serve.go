package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	web "peloton/internal/http"
)

// RunServe starts the embedded LAN web server.
func (app *App) RunServe(args []string) int {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(app.Stderr)

	port := fs.Int("port", app.Config.Port, "Port to listen on")
	host := fs.String("host", app.Config.Host, "Host/interface to bind")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	// Allow simple positional override? keep flags primary.
	srv := web.NewServer(*host, *port, app.ClientRepo, app.TicketRepo, app.FinanceRepo)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		fmt.Fprintf(app.Stdout, "\nReceived %s, shutting down web server...\n", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			fmt.Fprintf(app.Stderr, "Shutdown error: %v\n", err)
			return 1
		}
		return 0
	case err := <-errCh:
		if err != nil {
			fmt.Fprintf(app.Stderr, "Web server error: %v\n", err)
			return 1
		}
		return 0
	}
}
