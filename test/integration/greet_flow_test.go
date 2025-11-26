// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.

//go:build integration

// Package integration provides CLI integration tests.
//
// Integration tests verify the complete application flow by running
// the actual greeter binary and checking stdout, stderr, and exit codes.
//
// Run with: go test -v -tags=integration ./test/integration/...
//
// Prerequisites:
//   - Build the greeter binary first: go build -o greeter ./cmd/greeter
package integration

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/abitofhelp/hybrid_app_go/domain/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test tracking for summary.
var (
	testCount   int32
	passedCount int32
)

// registerTest tracks a test and its outcome for the final summary.
// Call at the start of each test function or subtest.
func registerTest(t *testing.T) {
	atomic.AddInt32(&testCount, 1)
	t.Cleanup(func() {
		if !t.Failed() {
			atomic.AddInt32(&passedCount, 1)
		}
	})
}

// greeterPath is the path to the greeter binary.
// Set during TestMain.
var greeterPath string

// TestMain builds the greeter binary before running tests.
func TestMain(m *testing.M) {
	// Build the greeter binary
	projectRoot := findProjectRoot()
	greeterPath = filepath.Join(projectRoot, "greeter_test_binary")

	cmd := exec.Command("go", "build", "-o", greeterPath, "./cmd/greeter")
	cmd.Dir = projectRoot
	if output, err := cmd.CombinedOutput(); err != nil {
		panic("Failed to build greeter: " + err.Error() + "\n" + string(output))
	}

	// Run tests
	code := m.Run()

	// Cleanup
	os.Remove(greeterPath)

	// Print summary banner
	test.PrintCategorySummary("INTEGRATION TESTS",
		int(atomic.LoadInt32(&testCount)),
		int(atomic.LoadInt32(&passedCount)))

	os.Exit(code)
}

// findProjectRoot finds the project root directory.
func findProjectRoot() string {
	// Start from current directory and walk up
	dir, _ := os.Getwd()
	for {
		// Check if this directory has cmd/greeter
		if _, err := os.Stat(filepath.Join(dir, "cmd", "greeter")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	// Fallback: assume we're in test/integration
	abs, _ := filepath.Abs(filepath.Join("..", ".."))
	return abs
}

// runGreeter executes the greeter binary with the given args.
func runGreeter(args ...string) (stdout, stderr string, exitCode int) {
	cmd := exec.Command(greeterPath, args...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()

	stdout = stdoutBuf.String()
	stderr = stderrBuf.String()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	return
}

// ============================================================================
// Valid Input Tests
// ============================================================================

func TestGreeter_ValidName_Success(t *testing.T) {
	registerTest(t)
	stdout, stderr, exitCode := runGreeter("Alice")

	assert.Equal(t, 0, exitCode, "exit code should be 0")
	assert.Equal(t, "Hello, Alice!\n", stdout, "stdout should contain greeting")
	assert.Empty(t, stderr, "stderr should be empty")
}

func TestGreeter_NameWithSpaces_Success(t *testing.T) {
	registerTest(t)
	stdout, stderr, exitCode := runGreeter("Bob Smith")

	assert.Equal(t, 0, exitCode, "exit code should be 0")
	assert.Equal(t, "Hello, Bob Smith!\n", stdout, "stdout should contain full name")
	assert.Empty(t, stderr, "stderr should be empty")
}

func TestGreeter_UnicodeCharacters_Success(t *testing.T) {
	registerTest(t)
	stdout, stderr, exitCode := runGreeter("José García")

	assert.Equal(t, 0, exitCode, "exit code should be 0")
	assert.Equal(t, "Hello, José García!\n", stdout, "stdout should contain unicode name")
	assert.Empty(t, stderr, "stderr should be empty")
}

func TestGreeter_SingleCharacter_Success(t *testing.T) {
	registerTest(t)
	stdout, stderr, exitCode := runGreeter("X")

	assert.Equal(t, 0, exitCode, "exit code should be 0")
	assert.Equal(t, "Hello, X!\n", stdout, "stdout should contain single char greeting")
	assert.Empty(t, stderr, "stderr should be empty")
}

func TestGreeter_MaxLengthName_Success(t *testing.T) {
	registerTest(t)
	// MaxNameLength is 100 characters
	maxName := strings.Repeat("a", 100)
	stdout, stderr, exitCode := runGreeter(maxName)

	assert.Equal(t, 0, exitCode, "exit code should be 0")
	assert.Contains(t, stdout, "Hello, "+maxName+"!", "stdout should contain max length greeting")
	assert.Empty(t, stderr, "stderr should be empty")
}

// ============================================================================
// Invalid Input Tests
// ============================================================================

func TestGreeter_NoArguments_ShowsUsage(t *testing.T) {
	registerTest(t)
	stdout, stderr, exitCode := runGreeter()

	assert.Equal(t, 1, exitCode, "exit code should be 1")
	assert.Empty(t, stdout, "stdout should be empty")
	assert.Contains(t, stderr, "Usage:", "stderr should contain usage")
}

func TestGreeter_TooManyArguments_ShowsUsage(t *testing.T) {
	registerTest(t)
	stdout, stderr, exitCode := runGreeter("Alice", "Bob")

	assert.Equal(t, 1, exitCode, "exit code should be 1")
	assert.Empty(t, stdout, "stdout should be empty")
	assert.Contains(t, stderr, "Usage:", "stderr should contain usage")
}

func TestGreeter_EmptyName_ValidationError(t *testing.T) {
	registerTest(t)
	stdout, stderr, exitCode := runGreeter("")

	assert.Equal(t, 1, exitCode, "exit code should be 1")
	assert.Empty(t, stdout, "stdout should be empty")
	assert.Contains(t, stderr, "Error:", "stderr should contain error")
	assert.Contains(t, stderr, "valid name", "stderr should mention valid name")
}

func TestGreeter_NameTooLong_ValidationError(t *testing.T) {
	registerTest(t)
	// MaxNameLength is 100, so 101 should fail
	longName := strings.Repeat("x", 101)
	stdout, stderr, exitCode := runGreeter(longName)

	assert.Equal(t, 1, exitCode, "exit code should be 1")
	assert.Empty(t, stdout, "stdout should be empty")
	assert.Contains(t, stderr, "Error:", "stderr should contain error")
}

// ============================================================================
// Edge Case Tests
// ============================================================================

func TestGreeter_WhitespaceOnlyName_ValidationError(t *testing.T) {
	registerTest(t)
	// Name with only whitespace should still be valid (preserved as-is)
	// Based on the Ada design: "whitespace is preserved exactly as provided"
	stdout, stderr, exitCode := runGreeter("   ")

	assert.Equal(t, 0, exitCode, "exit code should be 0 (whitespace preserved)")
	assert.Contains(t, stdout, "Hello,    !", "stdout should contain whitespace greeting")
	assert.Empty(t, stderr, "stderr should be empty")
}

func TestGreeter_SpecialCharacters_Success(t *testing.T) {
	registerTest(t)
	stdout, stderr, exitCode := runGreeter("O'Connor")

	assert.Equal(t, 0, exitCode, "exit code should be 0")
	assert.Equal(t, "Hello, O'Connor!\n", stdout, "stdout should contain special chars")
	assert.Empty(t, stderr, "stderr should be empty")
}

// ============================================================================
// Table-Driven Tests
// ============================================================================

func TestGreeter_ValidNames_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple name", "Alice", "Hello, Alice!\n"},
		{"name with space", "John Doe", "Hello, John Doe!\n"},
		{"unicode name", "北京", "Hello, 北京!\n"},
		{"emoji name", "🎉", "Hello, 🎉!\n"},
		{"hyphenated name", "Mary-Jane", "Hello, Mary-Jane!\n"},
		{"name with numbers", "User123", "Hello, User123!\n"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			registerTest(t)
			stdout, stderr, exitCode := runGreeter(tc.input)

			require.Equal(t, 0, exitCode, "exit code should be 0")
			assert.Equal(t, tc.expected, stdout)
			assert.Empty(t, stderr)
		})
	}
}

func TestGreeter_InvalidInputs_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectExitCode int
		expectInStderr string
	}{
		{"no args", []string{}, 1, "Usage:"},
		{"too many args", []string{"a", "b"}, 1, "Usage:"},
		{"empty string", []string{""}, 1, "Error:"},
		{"name too long", []string{strings.Repeat("x", 101)}, 1, "Error:"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			registerTest(t)
			stdout, stderr, exitCode := runGreeter(tc.args...)

			assert.Equal(t, tc.expectExitCode, exitCode)
			assert.Empty(t, stdout)
			assert.Contains(t, stderr, tc.expectInStderr)
		})
	}
}
