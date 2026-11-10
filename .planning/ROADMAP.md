# Roadmap: Peloton Service Manager

## Overview

A 4-phase execution roadmap delivering a full-featured bicycle repair workshop management system in a single compiled Go binary. Covers embedded SQLite data layer, CLI commands, interactive TUI workbench, and mobile-responsive LAN web app for Android phones.

## Phase List

- [ ] **Phase 1: Project Scaffolding & Core SQLite Storage**
  - Initialize Go project structure and modules.
  - Implement SQLite database migration and connection management (`modernc.org/sqlite`).
  - Create repositories and domain models for Clients, Service Tickets (with full status enum), and Financial Transactions.
  - Unit tests verifying data integrity and status transitions.
  - *Requirements covered: DATA-01, DATA-02, DATA-03, DATA-04*

- [ ] **Phase 2: CLI Commands & Financial Balancing**
  - Implement CLI subcommands (`peloton client`, `peloton service`, `peloton finance`) using Cobra or standard flag routing.
  - Support adding, searching, and filtering clients and repair jobs.
  - Implement cashflow logging (income/expenses) and monthly calculation summary command.
  - CLI integration tests.
  - *Requirements covered: CLI-01, CLI-02, CLI-03, FIN-01, FIN-02, FIN-03*

- [ ] **Phase 3: Interactive TUI Dashboard (Desktop Workbench)**
  - Implement Bubbletea/Lipgloss interactive TUI dashboard.
  - Multi-column visual repair queue (Pending, In Progress, Waiting for Parts, Ready for Pickup, Completed).
  - Keyboard navigation, detailed ticket inspection modal, and status advancement.
  - Monthly financial balance card embedded in TUI view.
  - *Requirements covered: TUI-01, TUI-02, TUI-03*

- [ ] **Phase 4: Embedded LAN Web Server & Android Mobile Interface**
  - Implement embedded HTTP server (`peloton serve --port 8080 --host 0.0.0.0`).
  - Embed HTML templates and responsive CSS assets (`go:embed`).
  - Mobile-first layout for Android smartphone browsers: active queue view, tap-to-call client button, quick ticket status update, and new repair intake form.
  - Mobile monthly cashflow summary page.
  - *Requirements covered: WEB-01, WEB-02, WEB-03, WEB-04*

## Progress Tracking

| Phase | Status | Requirements | Completed Date |
|-------|--------|--------------|----------------|
| Phase 1: Project Scaffolding & Core SQLite Storage | Ready | DATA-01, DATA-02, DATA-03, DATA-04 | — |
| Phase 2: CLI Commands & Financial Balancing | Pending | CLI-01, CLI-02, CLI-03, FIN-01, FIN-02, FIN-03 | — |
| Phase 3: Interactive TUI Dashboard | Pending | TUI-01, TUI-02, TUI-03 | — |
| Phase 4: Embedded LAN Web Server & Android Access | Pending | WEB-01, WEB-02, WEB-03, WEB-04 | — |

---
*Roadmap defined: 2026-11-10*
