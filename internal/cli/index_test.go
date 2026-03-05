package cli

import (
	"strings"
	"testing"
)

func TestReportServerErrors_NoErrors(t *testing.T) {
	if err := reportServerErrors(nil); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if err := reportServerErrors([]string{}); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestReportServerErrors_ReturnsNonZeroSignal(t *testing.T) {
	err := reportServerErrors([]string{"extract error src/a.py"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "server-reported errors") {
		t.Fatalf("unexpected error: %v", err)
	}
}
