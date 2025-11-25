// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.

package adapter

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	domerr "github.com/abitofhelp/hybrid_app_go/domain/error"
	"github.com/abitofhelp/hybrid_app_go/domain/test"
)

// TestInfrastructureAdapterConsoleWriter tests the NewWriter adapter factory.
//
// This demonstrates canonical infrastructure-layer unit testing:
//   - Testing adapters in isolation using bytes.Buffer
//   - Testing successful write operations
//   - Testing error handling (I/O errors mapped to InfrastructureError)
//   - Testing context cancellation handling
//   - Testing panic recovery at boundaries
func TestInfrastructureAdapterConsoleWriter(t *testing.T) {
	tf := test.New("Infrastructure.Adapter.ConsoleWriter")

	// ========================================================================
	// Successful Write Tests
	// ========================================================================

	// Test: NewWriter with bytes.Buffer captures output
	var buf bytes.Buffer
	writer := NewWriter(&buf)
	result := writer(context.Background(), "Hello, World!")

	tf.RunTest("Successful write - IsOk returns true", result.IsOk())
	tf.RunTest("Successful write - buffer contains message",
		buf.String() == "Hello, World!\n")

	// Test: Multiple writes accumulate
	buf.Reset()
	writer = NewWriter(&buf)
	_ = writer(context.Background(), "First")
	_ = writer(context.Background(), "Second")

	tf.RunTest("Multiple writes - buffer contains both",
		buf.String() == "First\nSecond\n")

	// Test: Empty message works
	buf.Reset()
	writer = NewWriter(&buf)
	result = writer(context.Background(), "")

	tf.RunTest("Empty message - IsOk returns true", result.IsOk())
	tf.RunTest("Empty message - buffer contains newline only",
		buf.String() == "\n")

	// Test: Unicode message works
	buf.Reset()
	writer = NewWriter(&buf)
	result = writer(context.Background(), "Hello, 世界! 🌍")

	tf.RunTest("Unicode message - IsOk returns true", result.IsOk())
	tf.RunTest("Unicode message - buffer contains unicode",
		buf.String() == "Hello, 世界! 🌍\n")

	// ========================================================================
	// Context Cancellation Tests
	// ========================================================================

	// Test: Cancelled context returns InfrastructureError
	buf.Reset()
	writer = NewWriter(&buf)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	result = writer(ctx, "Should not write")

	tf.RunTest("Cancelled context - IsError returns true", result.IsError())
	tf.RunTest("Cancelled context - error kind is InfrastructureError",
		result.ErrorInfo().Kind == domerr.InfrastructureError)
	tf.RunTest("Cancelled context - error message mentions cancelled",
		containsSubstring(result.ErrorInfo().Message, "cancelled"))
	tf.RunTest("Cancelled context - buffer is empty", buf.Len() == 0)

	// ========================================================================
	// I/O Error Handling Tests
	// ========================================================================

	// Test: Writer that returns error gets mapped to InfrastructureError
	failWriter := &failingWriter{err: errors.New("disk full")}
	writer = NewWriter(failWriter)
	result = writer(context.Background(), "Test")

	tf.RunTest("I/O error - IsError returns true", result.IsError())
	tf.RunTest("I/O error - error kind is InfrastructureError",
		result.ErrorInfo().Kind == domerr.InfrastructureError)
	tf.RunTest("I/O error - error message contains original error",
		containsSubstring(result.ErrorInfo().Message, "disk full"))

	// ========================================================================
	// Panic Recovery Tests
	// ========================================================================

	// Test: Panicking writer gets recovered and mapped to InfrastructureError
	panicWriter := &panickingWriter{}
	writer = NewWriter(panicWriter)
	result = writer(context.Background(), "Test")

	tf.RunTest("Panic recovery - IsError returns true", result.IsError())
	tf.RunTest("Panic recovery - error kind is InfrastructureError",
		result.ErrorInfo().Kind == domerr.InfrastructureError)
	tf.RunTest("Panic recovery - error message mentions panic",
		containsSubstring(result.ErrorInfo().Message, "panic"))

	// ========================================================================
	// Convenience Function Tests
	// ========================================================================

	// Test: NewConsoleWriter returns a valid WriterFunc
	consoleWriter := NewConsoleWriter()
	tf.RunTest("NewConsoleWriter - returns non-nil", consoleWriter != nil)

	// Test: NewStderrWriter returns a valid WriterFunc
	stderrWriter := NewStderrWriter()
	tf.RunTest("NewStderrWriter - returns non-nil", stderrWriter != nil)

	// Print summary
	tf.Summary(t)
}

// ============================================================================
// Test Helpers
// ============================================================================

// failingWriter is an io.Writer that always returns an error.
type failingWriter struct {
	err error
}

func (w *failingWriter) Write(p []byte) (n int, err error) {
	return 0, w.err
}

// panickingWriter is an io.Writer that panics on Write.
type panickingWriter struct{}

func (w *panickingWriter) Write(p []byte) (n int, err error) {
	panic("simulated panic in writer")
}

// containsSubstring checks if s contains substr (case-sensitive).
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Ensure test helper types implement io.Writer
var _ io.Writer = (*failingWriter)(nil)
var _ io.Writer = (*panickingWriter)(nil)
