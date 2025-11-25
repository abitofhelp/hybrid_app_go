// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.

//go:build e2e

// Package e2e provides end-to-end black-box tests.
//
// E2E tests execute the actual built binary and verify behavior from
// a user's perspective. They test:
// - CLI argument parsing
// - Exit codes
// - stdout/stderr output
//
// Run with: go test -v -tags=e2e ./test/e2e/...
//
// Prerequisites: Binary must be built first (make build)
package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/abitofhelp/hybrid_app_go/domain/test"
)

// findBinary locates the greeter binary.
// Searches in cmd/greeter/ relative to project root.
func findBinary() (string, error) {
	// Try to find project root by looking for go.work
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up to find project root
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.work")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}

	binary := filepath.Join(dir, "cmd", "greeter", "greeter")
	if _, err := os.Stat(binary); err != nil {
		return "", err
	}
	return binary, nil
}

// TestCLIGreeter tests the greeter CLI binary end-to-end.
func TestCLIGreeter(t *testing.T) {
	tf := test.New("E2E.CLI.Greeter")

	binary, err := findBinary()
	if err != nil {
		t.Skipf("Binary not found (run 'make build' first): %v", err)
	}

	// ========================================================================
	// Test: Valid name produces greeting
	// ========================================================================

	cmd := exec.Command(binary, "Alice")
	output, err := cmd.Output()
	tf.RunTest("Valid name - command succeeds", err == nil)
	tf.RunTest("Valid name - output contains greeting",
		strings.Contains(string(output), "Hello, Alice!"))

	// Check exit code
	tf.RunTest("Valid name - exit code is 0", cmd.ProcessState.ExitCode() == 0)

	// ========================================================================
	// Test: Name with spaces (quoted argument)
	// ========================================================================

	cmd2 := exec.Command(binary, "Bob Smith")
	output2, err2 := cmd2.Output()
	tf.RunTest("Name with spaces - command succeeds", err2 == nil)
	tf.RunTest("Name with spaces - output contains full name",
		strings.Contains(string(output2), "Hello, Bob Smith!"))

	// ========================================================================
	// Test: No arguments shows usage
	// ========================================================================

	cmd3 := exec.Command(binary)
	output3, _ := cmd3.CombinedOutput()
	tf.RunTest("No args - exit code is 1", cmd3.ProcessState.ExitCode() == 1)
	tf.RunTest("No args - output contains 'Usage'",
		strings.Contains(string(output3), "Usage"))

	// ========================================================================
	// Test: Empty name argument produces error
	// ========================================================================

	cmd4 := exec.Command(binary, "")
	output4, _ := cmd4.CombinedOutput()
	tf.RunTest("Empty name - exit code is 1", cmd4.ProcessState.ExitCode() == 1)
	tf.RunTest("Empty name - output contains error message",
		strings.Contains(string(output4), "Error") ||
			strings.Contains(string(output4), "empty"))

	// ========================================================================
	// Test: Multiple names (CLI joins them or uses first - behavior varies)
	// Note: This test verifies the CLI handles multiple args gracefully
	// ========================================================================

	cmd5 := exec.Command(binary, "First", "Second")
	output5, _ := cmd5.CombinedOutput()
	// CLI may join args, use first, or reject - just verify it runs
	tf.RunTest("Multiple args - produces some output", len(output5) > 0)

	// Print summary
	tf.Summary(t)
}
