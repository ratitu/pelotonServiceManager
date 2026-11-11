package http

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	webassets "peloton/assets/web"
	"peloton/internal/domain"
)

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.renderDashboard(w, r, "", "")
		return
	case http.MethodPost:
		s.handleCreateTicket(w, r)
		return
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (s *Server) handleCreateTicket(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.renderDashboard(w, r, "", "Invalid form submission")
		return
	}

	clientID, _ := strconv.ParseInt(strings.TrimSpace(r.FormValue("client_id")), 10, 64)
	desc := strings.TrimSpace(r.FormValue("description"))
	priority := strings.TrimSpace(r.FormValue("priority"))
	est, _ := strconv.ParseFloat(strings.TrimSpace(r.FormValue("estimated_cost")), 64)
	if priority == "" {
		priority = "Normal"
	}

	ticket := &domain.ServiceTicket{
		ClientID:      clientID,
		Description:   desc,
		Status:        domain.StatusPending,
		Priority:      priority,
		EstimatedCost: est,
	}

	if _, err := s.clientRepo.GetByID(clientID); err != nil {
		s.renderDashboard(w, r, "", fmt.Sprintf("Client #%d not found", clientID))
		return
	}
	if err := s.ticketRepo.Create(ticket); err != nil {
		s.renderDashboard(w, r, "", "Failed to create ticket: "+err.Error())
		return
	}

	http.Redirect(w, r, "/?created=1", http.StatusSeeOther)
}

func (s *Server) handleFinancePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	year, month, _ := time.Now().Date()
	summary, err := s.financeRepo.GetMonthlySummary(year, int(month))
	if err != nil {
		http.Error(w, "failed to load finance summary", http.StatusInternalServerError)
		return
	}

	tpl := template.Must(template.New("finance").Funcs(template.FuncMap{
		"money": func(v float64) string { return fmt.Sprintf("$%.2f", v) },
	}).Parse(financePageHTML))
	_ = tpl.Execute(w, map[string]any{"Summary": summary})
}

func (s *Server) handleTicketsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	var filter domain.TicketStatus
	if status != "" {
		filter = domain.TicketStatus(status)
	}
	tickets, err := s.ticketRepo.List(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, tickets)
}

func (s *Server) handleTicketActionAPI(w http.ResponseWriter, r *http.Request) {
	// POST /api/tickets/{id}/status
	path := strings.TrimPrefix(r.URL.Path, "/api/tickets/")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 2 || parts[1] != "status" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid ticket id", http.StatusBadRequest)
		return
	}

	statusStr := ""
	ct := r.Header.Get("Content-Type")
	if strings.Contains(ct, "application/json") {
		var body struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
			statusStr = body.Status
		}
	} else {
		_ = r.ParseForm()
		statusStr = r.FormValue("status")
	}
	statusStr = strings.TrimSpace(statusStr)
	if statusStr == "" {
		http.Error(w, "status is required", http.StatusBadRequest)
		return
	}

	matched := domain.TicketStatus(statusStr)
	for _, st := range domain.AllTicketStatuses {
		if strings.EqualFold(string(st), statusStr) {
			matched = st
			break
		}
	}
	if !domain.IsValidStatus(matched) {
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}

	if err := s.ticketRepo.UpdateStatus(id, matched, nil); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// For browser form posts, redirect back to dashboard.
	if !strings.Contains(ct, "application/json") {
		http.Redirect(w, r, "/?updated=1", http.StatusSeeOther)
		return
	}
	writeJSON(w, map[string]any{"ok": true, "id": id, "status": matched})
}

func (s *Server) handleFinanceSummaryAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	year, month, _ := time.Now().Date()
	summary, err := s.financeRepo.GetMonthlySummary(year, int(month))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, summary)
}

func (s *Server) renderDashboard(w http.ResponseWriter, r *http.Request, successMsg, errMsg string) {
	if r.URL.Query().Get("created") == "1" {
		successMsg = "Repair ticket created."
	}
	if r.URL.Query().Get("updated") == "1" {
		successMsg = "Ticket status updated."
	}

	tickets, err := s.ticketRepo.List("")
	if err != nil {
		http.Error(w, "failed to load tickets", http.StatusInternalServerError)
		return
	}
	clients, err := s.clientRepo.List()
	if err != nil {
		http.Error(w, "failed to load clients", http.StatusInternalServerError)
		return
	}
	year, month, _ := time.Now().Date()
	summary, _ := s.financeRepo.GetMonthlySummary(year, int(month))

	tpl, err := template.New("index").Funcs(template.FuncMap{
		"statusOptions": func() []domain.TicketStatus { return domain.AllTicketStatuses },
		"money":         func(v float64) string { return fmt.Sprintf("$%.2f", v) },
		"tel":           func(p string) string { return "tel:" + strings.ReplaceAll(p, " ", "") },
		"eqStatus":      func(a domain.TicketStatus, b string) bool { return string(a) == b },
	}).Parse(indexPageHTML)
	if err != nil {
		// Fallback to embedded static file if template parse fails
		b, readErr := webassets.Files.ReadFile("index.html")
		if readErr != nil {
			http.Error(w, "template error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(b)
		return
	}

	data := map[string]any{
		"Tickets":   tickets,
		"Clients":   clients,
		"Summary":   summary,
		"Success":   successMsg,
		"Error":     errMsg,
		"Statuses":  domain.AllTicketStatuses,
		"Year":      year,
		"Month":     int(month),
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tpl.Execute(w, data); err != nil {
		http.Error(w, "render error", http.StatusInternalServerError)
		return
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	_ = enc.Encode(v)
}

const indexPageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Peloton Service Manager</title>
  <link rel="stylesheet" href="/static/style.css">
</head>
<body>
  <header class="topbar">
    <h1>Peloton Service Manager</h1>
    <p class="subtitle">Mobile repair queue & cashflow (LAN)</p>
  </header>
  <main class="container">
    {{if .Success}}<section class="card"><p class="success">{{.Success}}</p></section>{{end}}
    {{if .Error}}<section class="card"><p class="error">{{.Error}}</p></section>{{end}}

    <section class="card">
      <h2>New Repair Intake</h2>
      <form method="post" action="/">
        <div class="row">
          <label for="client_id">Client</label>
          <select id="client_id" name="client_id" required>
            <option value="">Select a client</option>
            {{range .Clients}}<option value="{{.ID}}">#{{.ID}} {{.Name}} ({{.Phone}})</option>{{end}}
          </select>
        </div>
        <div class="row">
          <label for="description">Repair description</label>
          <textarea id="description" name="description" required placeholder="What needs to be fixed?"></textarea>
        </div>
        <div class="grid two">
          <div class="row">
            <label for="priority">Priority</label>
            <select id="priority" name="priority">
              <option>Normal</option><option>High</option><option>Urgent</option>
            </select>
          </div>
          <div class="row">
            <label for="estimated_cost">Estimated cost</label>
            <input id="estimated_cost" name="estimated_cost" type="number" step="0.01" min="0" placeholder="0.00">
          </div>
        </div>
        <button class="primary" type="submit">Create Ticket</button>
      </form>
      <p class="muted">Tip: add clients via CLI first if none appear.</p>
    </section>

    <section class="card">
      <h2>Active Queue</h2>
      {{if not .Tickets}}<p class="muted">No tickets yet.</p>{{end}}
      <div class="grid">
        {{range .Tickets}}
        <article class="ticket">
          <h3>#{{.ID}} {{.ClientName}}</h3>
          <p class="meta"><strong>Status:</strong> {{.Status}} | <strong>Priority:</strong> {{.Priority}}</p>
          <p class="meta"><strong>Phone:</strong> {{.ClientPhone}}</p>
          <p>{{.Description}}</p>
          <div class="actions">
            <a class="btn" href="{{tel .ClientPhone}}">Call Client</a>
          </div>
          <form method="post" action="/api/tickets/{{.ID}}/status" class="actions">
            <label for="status-{{.ID}}" class="muted">Quick status update</label>
            <select id="status-{{.ID}}" name="status">
              {{range statusOptions}}
              <option value="{{.}}" {{if eqStatus $.Status (printf "%s" .)}}selected{{end}}>{{.}}</option>
              {{end}}
            </select>
            <button class="primary" type="submit">Update Status</button>
          </form>
        </article>
        {{end}}
      </div>
    </section>

    <section class="card">
      <h2>Monthly Cashflow</h2>
      {{if .Summary}}
      <table>
        <tr><th>Period</th><td>{{.Summary.Month}}/{{.Summary.Year}}</td></tr>
        <tr><th>Total Income</th><td>{{money .Summary.TotalIncome}}</td></tr>
        <tr><th>Total Expense</th><td>{{money .Summary.TotalExpense}}</td></tr>
        <tr><th>Net Earnings</th><td>{{money .Summary.NetEarnings}}</td></tr>
        <tr><th>Completed Services</th><td>{{.Summary.ServicesCount}}</td></tr>
      </table>
      <div class="actions" style="margin-top:12px"><a class="btn" href="/finance">Open Finance Page</a></div>
      {{else}}
      <p class="muted">No finance data available.</p>
      {{end}}
    </section>
  </main>
</body>
</html>`

const financePageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Peloton Finance Summary</title>
  <link rel="stylesheet" href="/static/style.css">
</head>
<body>
  <header class="topbar">
    <h1>Monthly Cashflow</h1>
    <p class="subtitle">Phone-friendly financial summary</p>
  </header>
  <main class="container">
    <section class="card">
      {{if .Summary}}
      <table>
        <tr><th>Period</th><td>{{.Summary.Month}}/{{.Summary.Year}}</td></tr>
        <tr><th>Total Income</th><td>{{money .Summary.TotalIncome}}</td></tr>
        <tr><th>Total Expense</th><td>{{money .Summary.TotalExpense}}</td></tr>
        <tr><th>Net Earnings</th><td>{{money .Summary.NetEarnings}}</td></tr>
        <tr><th>Completed Services</th><td>{{.Summary.ServicesCount}}</td></tr>
      </table>
      {{else}}<p class="muted">No data.</p>{{end}}
      <div class="actions" style="margin-top:12px"><a class="btn" href="/">Back to Queue</a></div>
    </section>
  </main>
</body>
</html>`
