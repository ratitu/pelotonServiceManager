# Peloton Service Manager

Peloton Service Manager is a bicycle repair workshop management tool built in Go. It runs as a single binary and provides:

- **CLI commands** for clients, service tickets, and finances
- **Interactive TUI dashboard** (Bubble Tea + Lip Gloss) for desktop workflow
- **Embedded LAN web server** with a mobile-friendly interface for Android/phone access on the local network
- **Embedded SQLite storage** (`modernc.org/sqlite`, pure Go)

## Features

- Client management (add, list, search)
- Service ticket queue with status tracking (`Pending`, `In Progress`, `Waiting for Parts`, `Ready for Pickup`, `Completed`)
- Financial tracking (income/expense transactions and monthly summary)
- Desktop TUI workbench
- Mobile-first web dashboard for LAN use:
  - Active repair queue view
  - New repair intake form
  - Quick ticket status updates
  - Tap-to-call client links (`tel:`)
  - Monthly cashflow summary page/API

## Requirements

- Go (project currently targets Go `1.19` per `go.mod`)
- Linux/macOS/Windows environment capable of building Go binaries (developed for Debian/Linux use)

No external database server is required; data is stored in a local SQLite file.

## Build

```bash
go build -o ./peloton ./cmd/peloton
```

## Database location

By default, Peloton resolves the SQLite database path in this order:

1. `--db <path>` (global CLI flag)
2. `PELOTON_DB_PATH` environment variable
3. `$XDG_DATA_HOME/peloton/peloton.db` (typically `~/.local/share/peloton/peloton.db`)
4. `./peloton.db` (fallback)

The data directory is created automatically if needed.

## Usage

### Show help

```bash
./peloton --help
./peloton client --help
./peloton service --help
./peloton finance --help
```

### Launch TUI (desktop dashboard)

```bash
./peloton
# or explicitly
./peloton tui
```

> When run in a terminal with no subcommand, Peloton launches the TUI by default.

### Start embedded LAN web server (mobile access)

```bash
./peloton serve --port 8080 --host 0.0.0.0
```

Options:
- `--port <int>`: HTTP port (default `8080`)
- `--host <addr>`: bind host/interface (default `0.0.0.0`)

On startup, the server logs local and LAN access URLs (when a non-loopback IPv4 interface is detected), for example:

- `http://localhost:8080`
- `http://<LAN-IP>:8080`

Open the LAN URL from an Android phone/browser on the same network.

### CLI examples

#### Clients

```bash
# Add a client
./peloton client add --name "Ana Silva" --phone "+55 11 99999-0000" --bike "Caloi Elite" --notes "Prefers WhatsApp"

# List clients
./peloton client list

# Search clients
./peloton client search "Ana"
```

#### Service tickets

```bash
# Create a repair ticket
./peloton service add --client 1 --desc "Brake adjustment and chain lubrication" --priority High --est 120.00

# List tickets (optionally filter by status)
./peloton service list
./peloton service list --status "In Progress"

# Update ticket status (optional final cost)
./peloton service status 1 "Ready for Pickup"
./peloton service status 1 Completed --cost 135.00
```

#### Finance

```bash
# Log income / expense
./peloton finance income --amount 135.00 --category "Repair" --notes "Ticket #1"
./peloton finance expense --amount 40.00 --category "Parts" --notes "Brake pads"

# Monthly summary and transaction list
./peloton finance summary
./peloton finance list
```

## Web interface (Phase 4)

The embedded web UI is designed for local workshop use (LAN-only):

- `GET /` — mobile dashboard (queue, intake form, finance summary)
- `GET /finance` — monthly cashflow page
- `GET /static/*` — embedded static assets (HTML/CSS)
- `GET /api/tickets` — list tickets (optional `?status=...`)
- `POST /api/tickets/{id}/status` — update ticket status (`status` via form or JSON)
- `GET /api/finance/summary` — monthly finance summary JSON

> Tip: Add clients via CLI first so they appear in the web intake form client selector.

## Project structure (high level)

- `cmd/peloton/` — application entrypoint
- `internal/cli/` — CLI command handlers (`client`, `service`, `finance`, `serve`, TUI launcher)
- `internal/tui/` — Bubble Tea/Lip Gloss desktop UI
- `internal/http/` — embedded LAN web server and handlers
- `internal/store/sqlite/` — SQLite connection, migrations, repositories
- `internal/domain/` — domain models and validation
- `assets/web/` — embedded web UI assets (HTML/CSS)
- `.planning/` — GSD planning artifacts (roadmap, plans, research/state)

## Development notes

- The web UI assets are embedded into the binary at build time (`go:embed`), so rebuilding is required after changing files in `assets/web/`.
- The project is intended for trusted LAN environments; there is no multi-user authentication layer by default.

## Testing

Run tests with:

```bash
go test ./...
```

(Individual packages can be tested as needed, e.g. `go test ./internal/store/sqlite ./internal/domain`.)

## License

No license file is included in this repository yet. Add one if you plan to distribute the project.