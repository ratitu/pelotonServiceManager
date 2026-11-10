# Peloton Service Manager - Agent Guidelines

This project uses the GSD (Get Stuff Done) workflow for development.

## Project Structure
- `.planning/`: Project memory, roadmap, and state tracking.
  - `PROJECT.md`: Vision, constraints, and core decisions.
  - `REQUIREMENTS.md`: v1 requirement definitions with traceability.
  - `ROADMAP.md`: Phased execution roadmap.
  - `STATE.md`: Real-time phase state and progress tracking.

## Technology Stack
- **Language**: Go
- **Database**: SQLite (`modernc.org/sqlite` pure-Go driver)
- **TUI**: Charm Bubbletea + Lipgloss
- **Web**: Standard Go `net/http` with embedded HTML/CSS via `embed.FS`
- **Target OS**: Debian Linux & Android browser (via local LAN)

## Workflow Commands
- Start next phase: `/gsd-plan-phase 1` or `/gsd-discuss-phase 1`
- Check status: `/gsd-progress`
