# Peloton Service Manager

## What This Is

A lightweight, self-contained workshop management system for a bicycle repair shop running on Debian Linux. It tracks customer contact details, organizes repair jobs through an active service queue, logs income and expenses, and calculates monthly net earnings, accessible via CLI and TUI locally, and via a responsive web interface across the local network on Android devices.

## Core Value

Keep the bike repair queue and monthly earnings crystal clear and up-to-date from both the Debian desktop workbench (CLI/TUI) and an Android phone on the shop's Wi-Fi.

## Current State

**Milestone:** v1.0 (implementation and planning reconciliation complete)
**Status:** Shipped scope reconciled; closeout archive/audit artifacts prepared, commit/tag pending
**Capabilities delivered:**
- Embedded SQLite storage with client, ticket, and financial transaction records
- CLI for client/service/finance workflows and monthly summaries
- Interactive Bubbletea/Lipgloss TUI workbench with queue columns, ticket status updates, and finance widget
- Embedded LAN web server (`peloton serve`) with mobile-responsive dashboard for intake, queue/status actions, tap-to-call, and monthly cashflow views

**Validation snapshot:**
- Phase 4 UAT: `.planning/phases/04-LAN-Web/04-UAT.md` complete (6/6 passed, 0 issues)
- v1.0 reconciliation audit: `.planning/v1.0-MILESTONE-AUDIT.md` (`passed_with_notes`)

## Requirements

### Validated

- ✓ **DATA-01**: SQLite database initializes automatically with clients, service tickets, and financial transactions — v1.0 (Phase 1)
- ✓ **DATA-02**: Client profiles store name, phone, bike details, and notes — v1.0 (Phase 1)
- ✓ **DATA-03**: Service tickets link to clients with description, priority, cost fields, and lifecycle status — v1.0 (Phase 1)
- ✓ **DATA-04**: Ticket status lifecycle supports Pending/In Progress/Waiting for Parts/Ready for Pickup/Completed — v1.0 (Phase 1)
- ✓ **CLI-01**: CLI client create/list/search commands — v1.0 (Phase 2)
- ✓ **CLI-02**: CLI service ticket create/list/status commands — v1.0 (Phase 2)
- ✓ **CLI-03**: CLI finance logging and monthly summary commands — v1.0 (Phase 2)
- ✓ **FIN-01**: Record income transactions (including ticket-linked service revenue) — v1.0 (Phase 2)
- ✓ **FIN-02**: Record expense transactions for parts/supplies — v1.0 (Phase 2)
- ✓ **FIN-03**: Calculate monthly net balance filterable by month/year — v1.0 (Phase 2)
- ✓ **TUI-01**: Interactive Bubbletea/Lipgloss queue dashboard with status columns — v1.0 (Phase 3)
- ✓ **TUI-02**: Keyboard navigation, ticket inspection, and in-place status updates — v1.0 (Phase 3)
- ✓ **TUI-03**: TUI monthly earnings/finance overview — v1.0 (Phase 3)
- ✓ **WEB-01**: Embedded HTTP server on configurable host/port serving self-contained assets — v1.0 (Phase 4)
- ✓ **WEB-02**: Mobile-responsive web dashboard for Android browsers on LAN — v1.0 (Phase 4)
- ✓ **WEB-03**: Web ticket intake, quick status updates, and tap-to-call client action — v1.0 (Phase 4)
- ✓ **WEB-04**: Mobile-accessible monthly cashflow view/page/API — v1.0 (Phase 4)

### Active

(None — define next milestone requirements with `/gsd-new-milestone`)

### Out of Scope

- Automated SMS/WhatsApp API integration (manual calls/messages only).
- Public cloud hosting / internet domain exposure (local LAN only).
- Multi-tenant shop user accounts and complex RBAC.
- Payment gateway/card processing integrations.

## Context

- Target OS: Debian Linux desktop.
- Network setup: Local shop Wi-Fi/LAN; Android phone browser access via local IP/port (e.g. `http://192.168.1.x:8080`).
- Architecture: Single compiled Go binary with SQLite, CLI/TUI, and embedded web assets.
- v1.0 closeout state: code and planning artifacts reconciled; uncommitted working-tree changes remain pending commit/tag decisions.

## Constraints

- **Tech stack**: Go, pure-Go SQLite driver, Bubbletea/Lipgloss, standard `net/http`, embedded assets.
- **Environment**: Debian desktop with minimal overhead; local-first operation.
- **Network**: LAN-only access, no public exposure required.

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Go single binary | Bundles CLI, TUI, web server, and embedded assets into one portable executable | ✓ Good |
| SQLite embedded database | Zero maintenance, single-file local storage for single-workshop scale | ✓ Good |
| Responsive web for phone | Enables Android browser access without custom app build/distribution | ✓ Good |
| Manual client notification | Avoids third-party messaging APIs/costs and keeps workflow local | ✓ Good |
| Reuse domain/store across interfaces | Keeps CLI/TUI/web behavior and financial totals consistent | ✓ Good |

## Next Milestone Goals

- Run fresh verification (`go test ./...`, `go build ./...`) and capture results in closeout commit notes.
- Perform real-device Android-on-LAN usability sanity check (forms, status updates, `tel:` links, finance view).
- Finalize git commit/tag for v1.0 release artifacts and implementation.
- Plan v1.1+ improvements (e.g. edit/delete flows, authentication/host restrictions hardening, reporting exports, backups).

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd-transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd-complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-11-11 after v1.0 strict closeout reconciliation*
