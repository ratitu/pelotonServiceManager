# Requirements: Peloton Service Manager

**Defined:** 2026-11-10
**Core Value:** Keep the bike repair queue and monthly earnings crystal clear and up-to-date from both the Debian desktop workbench (CLI/TUI) and an Android phone on the shop's Wi-Fi.

## v1.0 Requirements

Requirements for milestone v1.0. Each maps to roadmap phases and is reconciled against implemented/validated scope.

### Core Data & Storage

- [x] **DATA-01**: SQLite database initializes automatically on startup with tables for clients, service tickets, and financial transactions.
- [x] **DATA-02**: Store client profiles with name, phone number, bike brand/model, and notes.
- [x] **DATA-03**: Create service tickets linked to clients with description, priority, estimated completion, and lifecycle status.
- [x] **DATA-04**: Track ticket status flow: `Pending` -> `In Progress` -> `Waiting for Parts` -> `Ready for Pickup` -> `Completed`.

### CLI Interface

- [x] **CLI-01**: CLI command to create, search, and list clients (`peloton client add/list/search`).
- [x] **CLI-02**: CLI command to create, update status, and view service queue tickets (`peloton service add/list/status`).
- [x] **CLI-03**: CLI command to log cash in/out and display monthly financial summary (`peloton finance income/expense/summary`).

### TUI Dashboard

- [x] **TUI-01**: Interactive terminal dashboard powered by Bubbletea/Lipgloss showing categorized active queue columns.
- [x] **TUI-02**: Navigate tickets via keyboard, inspect ticket details, and update status in-place.
- [x] **TUI-03**: TUI monthly earnings overview panel showing total income, expenses, and net profit.

### Local Web & Android Access

- [x] **WEB-01**: Embedded HTTP server listening on configurable host/port (e.g. `0.0.0.0:8080`) serving self-contained HTML/CSS/assets.
- [x] **WEB-02**: Mobile-responsive web dashboard tailored for Android mobile browsers on local Wi-Fi.
- [x] **WEB-03**: Web interface supports adding new repair jobs, quick client phone call action, and advancing ticket statuses.
- [x] **WEB-04**: Monthly cashflow view accessible from phone browser (dashboard section, `/finance`, `/api/finance/summary`).

### Monthly Cashflow

- [x] **FIN-01**: Record revenue entries associated with completed service tickets and general shop sales.
- [x] **FIN-02**: Record expense entries for replacement parts and shop supplies.
- [x] **FIN-03**: Calculate monthly net balance (total cash in minus total cash out) filterable by month/year.

## Out of Scope

- Automated SMS/WhatsApp API messaging gateway (manual phone calls/messages only).
- Cloud syncing / multi-tenant authentication.
- Payment gateway card processing (cash/terminal receipts tracked manually).

## Traceability

| Requirement | Phase | Status | Evidence |
|-------------|-------|--------|----------|
| DATA-01 | Phase 1 | Complete | `01-01/01-02` summaries, SQLite store/migrations implementation |
| DATA-02 | Phase 1 | Complete | `01-01/01-02` summaries, client model/repository implementation |
| DATA-03 | Phase 1 | Complete | `01-01/01-02` summaries, service ticket model/repository implementation |
| DATA-04 | Phase 1 | Complete | `01-01/01-02` summaries, status lifecycle/domain + repo updates |
| CLI-01 | Phase 2 | Complete | `02-01` summary, CLI client commands/tests |
| CLI-02 | Phase 2 | Complete | `02-01` summary, CLI service commands/tests |
| CLI-03 | Phase 2 | Complete | `02-02` summary, CLI finance commands/tests |
| FIN-01 | Phase 2 | Complete | `02-02` summary, finance transaction logging (income) |
| FIN-02 | Phase 2 | Complete | `02-02` summary, finance transaction logging (expense) |
| FIN-03 | Phase 2 | Complete | `02-02` summary, monthly summary calculation/reporting |
| TUI-01 | Phase 3 | Complete | `03-01` summary, Bubbletea/Lipgloss dashboard columns |
| TUI-02 | Phase 3 | Complete | `03-01/03-02` summaries, navigation + detail modal/status updates |
| TUI-03 | Phase 3 | Complete | `03-02` summary, finance widget/panel in TUI |
| WEB-01 | Phase 4 | Complete | `04-01` summary, `serve` command + embedded HTTP server |
| WEB-02 | Phase 4 | Complete | `04-01` summary, embedded responsive mobile UI assets |
| WEB-03 | Phase 4 | Complete | `04-02` summary, intake/status/call actions in web UI |
| WEB-04 | Phase 4 | Complete | `04-02` summary, `/finance` + `/api/finance/summary` + dashboard cashflow; UAT passed |

**Coverage:**
- v1.0 requirements: 17 total
- Implemented/reconciled: 17
- Pending/unmapped: 0

## Verification Notes

- Phase 4 user acceptance testing recorded in `.planning/phases/04-LAN-Web/04-UAT.md` (`status: complete`, 6/6 passed, 0 issues).
- Phase 1–3 verification is reconciled from committed implementation history, tests present in repo, and created phase summaries; no separate formal UAT/audit artifacts were found.
- Recommended before final archive/tag: run a fresh full test/build pass in this environment and perform a real Android-on-LAN device sanity check.

---
*Requirements reconciled: 2026-11-11 for v1.0 strict closeout*
