# Phase 1: Project Scaffolding & Core SQLite Storage - Context

## Executive Summary

Phase 1 establishes the Go project repository structure, module setup, and the core SQLite storage engine. It provides the domain models, database schema, auto-migration mechanisms, and data access repositories for Clients, Service Tickets, and Financial Transactions.

## Decisions

### 1. Database Location & Storage Path
- **Default Path**: `~/.local/share/peloton/peloton.db` following the XDG base directory specification on Linux.
- **Configurability**: Support `PELOTON_DB_PATH` environment variable and `--db` flag override, falling back to local `./peloton.db` if XDG path is unavailable.
- **Driver**: Pure Go SQLite (`modernc.org/sqlite`) to eliminate CGO dependencies and guarantee cross-compilation simplicity on Debian.

### 2. Schema & Identity Strategy
- **Primary Keys**: Integer auto-increment IDs for friendly human CLI reference (e.g. `Ticket #14`) alongside indexed created/updated timestamps.
- **Foreign Keys**: Enabled (`PRAGMA foreign_keys = ON;`) linking `service_tickets.client_id -> clients.id` and `financial_transactions.ticket_id -> service_tickets.id`.
- **Status Lifecycle Enum**:
  - `Pending` (Initial ticket created)
  - `In Progress` (Work underway)
  - `Waiting for Parts` (Blocked waiting for component delivery)
  - `Ready for Pickup` (Service completed, awaiting customer retrieval)
  - `Completed` (Paid and collected)

### 3. Domain Model Specifications
- **Client**:
  - `ID`: integer primary key
  - `Name`: text, required
  - `Phone`: text, required (clean normalized format for quick dialing)
  - `BikeDetails`: text (make, model, color, serial number)
  - `Notes`: text (special instructions or repair history notes)
  - `CreatedAt`, `UpdatedAt`: ISO8601 timestamps
- **Service Ticket**:
  - `ID`: integer primary key
  - `ClientID`: integer foreign key
  - `Description`: text (repair tasks to do)
  - `Status`: text (lifecycle enum)
  - `Priority`: text (`Normal`, `High`, `Urgent`)
  - `EstimatedCost`: integer or decimal (cents or currency units)
  - `FinalCost`: integer or decimal
  - `CreatedAt`, `UpdatedAt`, `CompletedAt`: timestamps
- **Financial Transaction**:
  - `ID`: integer primary key
  - `TicketID`: nullable foreign key to service ticket
  - `Type`: `INCOME` (labor/parts sold) vs `EXPENSE` (parts/supplies bought)
  - `Category`: text (e.g., `Labor`, `Parts`, `Shop Supplies`)
  - `Amount`: decimal / cents
  - `Notes`: text
  - `Date`: date/timestamp of payment or expense

### 4. Code Organization
- Standard Go layout:
  - `cmd/peloton/`: Entry point main package
  - `internal/domain/`: Core entities and interfaces
  - `internal/store/sqlite/`: SQLite connection, migrations, and repository implementations
  - `internal/config/`: Configuration and path resolution

## Deferred Ideas
- Automatic daily database backup export (consider in later maintenance milestone).
- Full text search (FTS5) extension indexing for repair notes (can be added if ticket count grows into thousands).

---
*Context generated: 2026-11-10*
