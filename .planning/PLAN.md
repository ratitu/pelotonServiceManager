---
requirements:
  - WEB-01
  - WEB-02
  - WEB-03
  - WEB-04
---

# PLAN.md - Phase 4: Embedded LAN Web Server & Android Mobile Interface

## Phase 4 Goal
Implement a mobile-responsive LAN web server that provides the same functionality as the TUI dashboard, with mobile-first design for Android smartphones.

## Tasks Breakdown

### Task 4.1: Embedded HTTP Server Setup (WEB-01)
- Implement `peloton serve` command with port and host flags
- Set up HTTP server with proper routing for API endpoints and static assets
- Handle concurrent connections and graceful shutdown
- Implement authentication for LAN access (optional but recommended)

### Task 4.2: Template Embedding & Asset Management (WEB-02)
- Use `go:embed` to embed HTML templates and CSS/JS assets
- Create responsive HTML templates with mobile-first design
- Implement asset serving with proper MIME types and caching headers
- Handle template parsing and rendering with proper escaping

### Task 4.3: Mobile-First UI Implementation (WEB-03)
- Design responsive layout for 320px-1080px viewport range
- Implement active queue view with touch-friendly elements
- Add tap-to-call client button functionality
- Implement quick ticket status update via modal dialogs
- Create new repair intake form with mobile-optimized fields

### Task 4.4: Mobile Cashflow Summary Page (WEB-04)
- Create dedicated page showing monthly financial metrics
- Implement interactive charts or summary cards
- Add print-friendly layout option
- Ensure data synchronization with TUI financial data

### Task 4.5: Verification & Testing
- Implement end-to-end testing for web interface
- Test mobile responsiveness on multiple device sizes
- Verify financial data consistency between TUI and web
- Test concurrent access scenarios

## Dependencies
- Completed Phase 3 (TUI Dashboard) - must share data models and business logic
- Completed CLI commands - for API endpoint consistency
- Completed financial system - for cashflow data integration

## Success Criteria
- `./peloton serve --port 8080 --host 0.0.0.0` starts web server
- Mobile browser access shows responsive interface
- All mobile UI elements are touch-friendly and functional
- Financial data is consistent between TUI and web interfaces
- No console errors in browser dev tools