package http

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	webassets "peloton/assets/web"
	"peloton/internal/store/sqlite"
)

// Server serves the embedded web UI and API over LAN.
type Server struct {
	addr        string
	httpServer  *http.Server
	clientRepo  *sqlite.ClientRepository
	ticketRepo  *sqlite.TicketRepository
	financeRepo *sqlite.FinanceRepository
}

// NewServer creates a configured web server.
func NewServer(host string, port int, clientRepo *sqlite.ClientRepository, ticketRepo *sqlite.TicketRepository, financeRepo *sqlite.FinanceRepository) *Server {
	if strings.TrimSpace(host) == "" {
		host = "0.0.0.0"
	}
	if port <= 0 {
		port = 8080
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	s := &Server{
		addr:        addr,
		clientRepo:  clientRepo,
		ticketRepo:  ticketRepo,
		financeRepo: financeRepo,
	}
	s.httpServer = &http.Server{
		Addr:              addr,
		Handler:           s.routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	return s
}

// Addr returns the listen address.
func (s *Server) Addr() string { return s.addr }

// Start runs the HTTP server and logs local access URLs.
func (s *Server) Start() error {
	ip := GetLocalIP()
	log.Printf("Peloton web server starting on http://%s", s.addr)
	if ip != "" {
		_, port, _ := net.SplitHostPort(s.addr)
		if port == "" {
			port = "8080"
		}
		log.Printf("LAN access: http://%s:%s", ip, port)
	}
	log.Printf("Local access: http://localhost:%s", portFromAddr(s.addr))
	err := s.httpServer.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()

	// Static assets embedded in binary.
	staticFS, err := fs.Sub(webassets.Files, ".")
	if err != nil {
		staticFS = webassets.Files
	}
	fileServer := http.FileServer(http.FS(staticFS))
	mux.HandleFunc("/static/index.html", func(w http.ResponseWriter, r *http.Request) {
		b, err := fs.ReadFile(staticFS, "index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(b)
	})
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Pages and API
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/finance", s.handleFinancePage)
	mux.HandleFunc("/api/tickets", s.handleTicketsAPI)
	mux.HandleFunc("/api/tickets/", s.handleTicketActionAPI)
	mux.HandleFunc("/api/finance/summary", s.handleFinanceSummaryAPI)

	return mux
}

// GetLocalIP returns a non-loopback IPv4 address for LAN access hints.
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok || ipNet.IP == nil || ipNet.IP.IsLoopback() {
			continue
		}
		ip := ipNet.IP.To4()
		if ip == nil {
			continue
		}
		return ip.String()
	}
	return ""
}

func portFromAddr(addr string) string {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "8080"
	}
	return port
}
