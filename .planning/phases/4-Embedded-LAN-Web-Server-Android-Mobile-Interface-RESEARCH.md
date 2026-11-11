# Phase 4: Embedded LAN Web Server & Android Mobile Interface - Research

**Researched:** 2026-11-10
**Domain:** Go Embedded Web Server, Mobile-First UI
**Confidence:** HIGH

## User Constraints

### Locked Decisions
- Go single binary with embedded assets (`go:embed`)
- SQLite embedded database
- LAN-only access (no public cloud/hosting)
- Manual notification (no SMS/WhatsApp API)

### the agent's Discretion
- Routing: Standard `net/http` (`ServeMux`) is recommended due to minimal requirements and performance/dependency trade-offs for this scale.
- CSS: Mobile-first responsive grid.
- LAN Discovery: Use `net.InterfaceAddrs()` to identify non-loopback IPv4 for server startup logging.

### Deferred Ideas (OUT OF SCOPE)
- Automated SMS/WhatsApp integration
- Public cloud hosting / port forwarding
- Multi-user authentication

## Architectural Responsibility Map

| Capability | Primary Tier | Secondary Tier | Rationale |
|------------|-------------|----------------|-----------|
| HTTP Server | Backend (Go) | — | Embedded in the single Go binary. |
| Web UI | Browser/Client | — | Static assets embedded and served via `embed.FS`. |
| LAN Discovery | Backend (Go) | — | Server startup logs the LAN IP for user access. |

## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Go stdlib | 1.19+ | Web Server, Routing | No external dependencies, minimal memory. |
| `embed` | stdlib | Asset Packaging | Canonical way to bundle assets into binary. |
| `html/template`| stdlib | UI Rendering | Built-in, safe, works with `fs.FS`. |

## Package Legitimacy Audit

| Package | Registry | Age | Downloads | Source Repo | Verdict | Disposition |
|---------|----------|-----|-----------|-------------|---------|-------------|
| `modernc.org/sqlite` | Go Proxy | 4 yrs | N/A | [modernc.org/sqlite](https://modernc.org/sqlite) | OK | Approved |

## Architecture Patterns

### Pattern 1: Embedded Static Assets
**What:** Use `go:embed` to package static assets (`html/`, `css/`, `js/`) into the binary.
**Example:**
```go
//go:embed assets/*
var assets embed.FS

func NewMux() *http.ServeMux {
    mux := http.NewServeMux()
    // StripPrefix handles path mapping
    fileServer := http.FileServer(http.FS(assets))
    mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
    return mux
}
```

### Pattern 2: LAN IP Discovery
**What:** Identify private network IP to guide users on mobile.
```go
func GetLocalIP() string {
    addrs, _ := net.InterfaceAddrs()
    for _, addr := range addrs {
        if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return "localhost"
}
```

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Routing | Custom regex matcher | `http.ServeMux` | Stdlib handles most cases securely. |
| Asset bundling | File-system-based assets | `go:embed` | Single binary is key for portability. |

## Common Pitfalls

- **Incorrect Path Stripping:** `http.FileServer` won't work correctly if the mount prefix isn't stripped via `http.StripPrefix`.
- **Private IP False Positives:** Loopback addresses (`127.0.0.1`) can appear in `net.InterfaceAddrs()`; always filter them out.
- **Mobile Click Targets:** Tap targets below 44px height are unusable for greasy mechanic hands. Use block buttons.

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Debian Desktop | Host OS | ✓ | 12 | — |
| Go Runtime | Build/Run | ✓ | 1.19.8 | — |

## Security Domain

### Applicable ASVS Categories
- **V5 Input Validation**: Ensure all ticket/client inputs are sanitized before SQL entry. SQLite parameters in `modernc.org/sqlite` provide protection against injection.

## Assumptions Log

*No claims found tagged `[ASSUMED]` — all claims are standard Go or verified patterns.*

## Open Questions

1. **Should the Android interface include dark mode?**
   - Recommendation: Use a system-aware CSS media query (`@media (prefers-color-scheme: dark)`) to ease eye strain in dark repair shop environments.

## Sources

### Primary (HIGH confidence)
- [pkg.go.dev/embed](https://pkg.go.dev/embed) - Standard library embedding patterns.
- [Ben Hoyt: Routing in Go](https://benhoyt.com/writings/go-routing/) - Standard library vs. router comparison.

### Secondary (MEDIUM confidence)
- [GoSamples: Local IP Discovery](https://gosamples.dev/local-ip-address/) - Verified pattern for LAN IP detection.

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH - Standard Go library.
- Architecture: HIGH - Proven Go deployment patterns.
- Pitfalls: MEDIUM - Mobile-specific UI constraints depend on mechanic feedback.

**Research date:** 2026-11-10
**Valid until:** 2026-12-10
