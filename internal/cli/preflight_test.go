package cli

import (
	"errors"
	"testing"

	"github.com/RevylAI/greenlight/internal/preflight"
)

// preflightExit only trips when --exit-code is set, and then on any CRITICAL or
// HIGH finding, or an incomplete scan.
func TestPreflightExit(t *testing.T) {
	defer func() { preflightExitCode = false }()

	crit := &preflight.Result{Summary: preflight.Summary{Critical: 1}}
	high := &preflight.Result{Summary: preflight.Summary{High: 1}}
	clean := &preflight.Result{Summary: preflight.Summary{Warns: 3}}
	// A crashed scanner surfaces only as WARN, but the scan is incomplete.
	incomplete := &preflight.Result{Summary: preflight.Summary{Warns: 1}, Incomplete: true}

	// Without the flag: never errors, whatever the findings.
	preflightExitCode = false
	for _, r := range []*preflight.Result{crit, high, clean, incomplete} {
		if err := preflightExit(r); err != nil {
			t.Errorf("no --exit-code: expected nil, got %v", err)
		}
	}

	// With the flag: critical, high, and an incomplete scan trip; clean does not.
	preflightExitCode = true
	if !errors.Is(preflightExit(crit), ErrThreshold) {
		t.Error("critical should trip --exit-code")
	}
	if !errors.Is(preflightExit(high), ErrThreshold) {
		t.Error("high should trip --exit-code")
	}
	if !errors.Is(preflightExit(incomplete), ErrThreshold) {
		t.Error("an incomplete scan (crashed scanner) should trip --exit-code")
	}
	if err := preflightExit(clean); err != nil {
		t.Errorf("clean static should not trip, got %v", err)
	}

}
