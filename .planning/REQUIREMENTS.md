# Requirements: Peloton Service Manager

**Defined:** 2026-11-10
**Core Value:** Keep the bike repair queue and monthly earnings crystal clear and up-to-date from both the Debian desktop workbench (CLI/TUI) and an Android phone on the shop's Wi-Fi.

## v1 Requirements

Requirements for initial release. Each maps to roadmap phases.

### Core Data & Storage

- [ ] **DATA-01**: SQLite database initializes automatically on startup with tables for clients, service tickets, and financial transactions.
- [ ] **DATA-02**: Store client profiles with name, phone number, bike brand/model, and notes.
- [ ] **DATA-03**: Create service tickets linked to clients with description, priority, estimated completion, and lifecycle status.
- [ ] **DATA-04**: Track ticket status flow: `Pending` -> `In Progress` -> `Waiting for Parts` -> `Ready for Pickup` -> `Completed`.

### CLI Interface

- [ ] **CLI-01**: CLI command to create, search, and list clients (`peloton client add/list/search`).
- [ ] **CLI-02**: CLI command to create, update status, and view service queue tickets (`peloton service add/list/status`).
- [ ] **CLI-03**: CLI command to log cash in/out and display monthly financial summary (`peloton finance add/summary`).

### TUI Dashboard

- [ ] **TUI-01**: Interactive terminal dashboard powered by Bubbletea/Lipgloss showing categorized active queue columns.
- [ ] **TUI-02**: Navigate tickets via keyboard, inspect ticket details, and update status in-place.
- [ ] **TUI-03**: TUI monthly earnings overview panel showing total income, expenses, and net profit.

### Local Web & Android Access

- [ ] **WEB-01**: Embedded HTTP server listening on configurable host/port (e.g. `0.0.0.0:8080`) serving self-contained HTML/CSS.
- [ ] **WEB-02**: Mobile-responsive web dashboard tailored for Android mobile browsers on local Wi-Fi.
- [ ] **WEB-03**: Web interface supports adding new repair jobs, quick client phone call/copy, and advancing ticket statuses.
- [ ] **WEB-04**: Monthly cashflow view accessible from phone browser.

### Monthly Cashflow

- [ ] **FIN-01**: Record revenue entries associated with completed service tickets and general shop sales.
- [ ] **FIN-02**: Record expense entries for replacement parts and shop supplies.
- [ ] **FIN-03**: Calculate monthly net balance (total cash in minus total cash out) filterable by month/year.

## Out of Scope

- Automated SMS/WhatsApp API messaging gateway (manual phone calls/messages only).
- Cloud syncing / multi-tenant authentication.
- Payment gateway card processing (cash/terminal receipts tracked manually).

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| DATA-01 | Phase 1 | Pending |
| DATA-02 | Phase 1 | Pending |
| DATA-03 | Phase 1 | Pending |
| DATA-04 | Phase 1 | Pending |
| CLI-01 | Phase 2 | Pending |
| CLI-02 | Phase 2 | Pending |
| CLI-03 | Phase 2 | Pending |
| FIN-01 | Phase 2 | Pending |
| FIN-02 | Phase 2 | Pending |
| FIN-03 | Phase 2 | Pending |
| TUI-01 | Phase 3 | Pending |
| TUI-02 | Phase 3 | Pending |
| TUI-03 | Phase 3 | Pending |
| WEB-01 | Phase 4 | Pending |
| WEB-02 | Phase 4 | Pending |
| WEB-03 | Phase 4 | Pending |
| WEB-04 | Phase 4 | Pending |

**Coverage:**
- v1 requirements: 17 total
- Mapped to phases: 17
- Unmapped: 0 ✓

---
*Requirements defined: 2026-11-10*
*Last updated: 2026-11-10 after initial definition*
