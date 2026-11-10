# Peloton Service Manager

## What This Is

A lightweight, self-contained workshop management system for a bicycle repair shop running on Debian Linux. It tracks customer contact details, organizes repair jobs through an active service queue, logs income and expenses, and calculates monthly net earnings, accessible via CLI and TUI locally, and via a responsive web interface across the local network on Android devices.

## Core Value

Keep the bike repair queue and monthly earnings crystal clear and up-to-date from both the Debian desktop workbench (CLI/TUI) and an Android phone on the shop's Wi-Fi.

## Business Context

- **Customer**: Solo bicycle mechanic / shop owner.
- **Revenue model**: Direct service labor and parts sales logged in shop.
- **Success metric**: Quick job lookup, zero forgotten repair tickets, and instant monthly cashflow balance.

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] **CLI interface**: Commands to quickly add, list, find, and update jobs, clients, and cash transactions.
- [ ] **TUI interface**: Terminal dashboard (Bubbletea/Charm) on Debian desktop displaying queue status columns and monthly summary.
- [ ] **LAN Web server**: Embedded HTTP server listening on local network providing a mobile-friendly web UI for Android browser access.
- [ ] **Client management**: Store customer name, phone number, notes, and bike model/make.
- [ ] **Service queue tracking**: Manage tickets with lifecycle: Pending -> In Progress -> Waiting for Parts -> Ready for Pickup -> Completed / Paid.
- [ ] **Earnings balance**: Record income (labor & parts) and costs (supplies/parts) with monthly summary calculations (cash in vs cash out, net earnings).
- [ ] **SQLite storage**: Single-file database with zero external dependencies and automatic schema setup.

### Out of Scope

- Automated SMS/WhatsApp API integration (shop owner calls/messages customers manually from phone).
- Public cloud hosting / internet domain exposure (operates strictly on Debian host and local LAN).
- Multi-tenant shop user accounts and complex role-based access control.

## Context

- Target OS: Debian Linux desktop.
- Network setup: Local shop Wi-Fi / LAN, enabling Android phone browser access via local IP:port (e.g., `http://192.168.1.X:8080`).
- Preferred architecture: Single compiled Go binary containing embedded assets (web templates/CSS/JS) and SQLite engine.

## Constraints

- **Tech stack**: Go single binary with embedded SQLite (`modernc.org/sqlite` or pure Go driver to avoid external C compiler issues if desired), Bubbletea/Lipgloss for TUI, standard net/http or lightweight router for local web.
- **Environment**: Debian desktop with minimal overhead, capable of running as a systemd user service or background daemon.
- **Network**: Local area network only, no public port forwards required.

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Go single binary | Bundles CLI, TUI, web server, and embedded assets into one portable executable | — Pending |
| SQLite embedded database | Zero maintenance, single file backup, fast and reliable for single-workshop scale | — Pending |
| Responsive web for phone | Allows seamless access from Android browser without needing custom APK compilation | — Pending |
| Manual client notification | Avoids third-party messaging costs and API credentials; keeps system offline-friendly | — Pending |

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
*Last updated: 2026-11-10 after initialization*
