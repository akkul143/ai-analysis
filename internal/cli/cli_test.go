package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := Execute([]string{"version"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected exit code 0, got %d: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "phase1") {
		t.Fatalf("expected phase1 version, got %q", stdout.String())
	}
}

func TestRunRequiresDryRunInPhaseOne(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := Execute([]string{"run"}, &stdout, &stderr); code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !strings.Contains(stderr.String(), "Phase 2") {
		t.Fatalf("expected Phase 2 explanation, got %q", stderr.String())
	}
}
