package cli

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestFinanceLogging(t *testing.T) {
	app, stdout, stderr, cleanup := setupTestApp(t)
	defer cleanup()

	// 1. Log Income
	code := app.Run([]string{"finance", "income", "--amount", "240.00", "--category", "Service", "--notes", "Full overhaul payment"})
	if code != 0 {
		t.Fatalf("finance income failed with exit code %d. Stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Income logged: +$240.00") {
		t.Errorf("expected income logged message, got: %s", stdout.String())
	}

	// 2. Reject zero or negative income
	stdout.Reset()
	stderr.Reset()
	code = app.Run([]string{"finance", "income", "--amount", "-10.00"})
	if code == 0 {
		t.Fatalf("expected error for negative amount, got code 0")
	}

	// 3. Log Expense
	stdout.Reset()
	stderr.Reset()
	code = app.Run([]string{"finance", "expense", "--amount", "95.50", "--category", "Parts", "--notes", "Shimano Ultegra chain"})
	if code != 0 {
		t.Fatalf("finance expense failed: %s", stderr.String())
	}
	if !strings.Contains(stdout.String(), "Expense logged: -$95.50") {
		t.Errorf("expected expense logged message, got: %s", stdout.String())
	}
}

func TestFinanceSummary(t *testing.T) {
	app, stdout, stderr, cleanup := setupTestApp(t)
	defer cleanup()

	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	// Register client and completed service
	app.Run([]string{"client", "add", "--name", "Tomasz Majewski", "--phone", "600-111-222"})
	app.Run([]string{"service", "add", "--client", "1", "--desc", "Bottom bracket replacement", "--est", "150.0"})
	app.Run([]string{"service", "status", "1", "Completed", "--cost", "150.0"})

	// Log Income and Expense
	app.Run([]string{"finance", "income", "--amount", "150.00", "--category", "Service", "--ticket", "1"})
	app.Run([]string{"finance", "expense", "--amount", "60.00", "--category", "Parts", "--notes", "Shimano BB-BBR60"})

	// Run summary
	stdout.Reset()
	stderr.Reset()
	code := app.Run([]string{"finance", "summary", "--month", fmt.Sprintf("%d", month), "--year", fmt.Sprintf("%d", year)})
	if code != 0 {
		t.Fatalf("finance summary failed: %s", stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "PELOTON FINANCIAL SUMMARY") {
		t.Errorf("summary missing title in:\n%s", output)
	}
	if !strings.Contains(output, "Completed Services : 1 jobs") {
		t.Errorf("summary missing completed service count in:\n%s", output)
	}
	if !strings.Contains(output, "Total Income (+)   : $150.00") {
		t.Errorf("summary missing income total in:\n%s", output)
	}
	if !strings.Contains(output, "Total Expense (-)  : $60.00") {
		t.Errorf("summary missing expense total in:\n%s", output)
	}
	if !strings.Contains(output, "Net Balance        : $90.00 (● PROFIT)") {
		t.Errorf("summary missing net profit in:\n%s", output)
	}

	// Test Transaction List
	stdout.Reset()
	stderr.Reset()
	code = app.Run([]string{"finance", "list"})
	if code != 0 {
		t.Fatalf("finance list failed: %s", stderr.String())
	}
	listOutput := stdout.String()
	if !strings.Contains(listOutput, "+$150.00") || !strings.Contains(listOutput, "-$60.00") {
		t.Errorf("transaction list missing items in:\n%s", listOutput)
	}
}
