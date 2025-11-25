// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.

//go:build integration

// Package integration provides cross-layer integration tests.
//
// Integration tests verify that real components work together correctly.
// They use real implementations (not mocks) to test complete flows.
//
// Run with: go test -v -tags=integration ./test/integration/...
package integration

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/abitofhelp/hybrid_app_go/application/command"
	"github.com/abitofhelp/hybrid_app_go/application/usecase"
	domerr "github.com/abitofhelp/hybrid_app_go/domain/error"
	"github.com/abitofhelp/hybrid_app_go/domain/test"
	"github.com/abitofhelp/hybrid_app_go/infrastructure/adapter"
)

// TestGreetFlowIntegration demonstrates integration testing pattern.
//
// This test uses REAL implementations across all layers:
// - Real Domain (Person value object)
// - Real Application (GreetUseCase)
// - Real Infrastructure (ConsoleWriter - captured stdout)
//
// Per Testing Standards: Integration tests use real components,
// mock only external services (databases, APIs, etc.)
func TestGreetFlowIntegration(t *testing.T) {
	tf := test.New("Integration.GreetFlow")

	// ========================================================================
	// Test: Full flow with valid input
	// ========================================================================

	// Capture stdout to verify console output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Setup: Wire real components together (same as bootstrap)
	// NewConsoleWriter returns a WriterFunc directly (not a struct)
	writerFunc := adapter.NewConsoleWriter()
	greetUseCase := usecase.NewGreetUseCase(writerFunc)

	// Execute: Run use case with valid input
	cmd := command.GreetCommand{Name: "Integration Test"}
	result := greetUseCase.Execute(cmd)

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Assert
	tf.RunTest("Valid input - IsOk returns true", result.IsOk())
	tf.RunTest("Valid input - output contains greeting",
		strings.Contains(buf.String(), "Hello, Integration Test!"))

	// ========================================================================
	// Test: Full flow with invalid (empty) input - error propagation
	// ========================================================================

	result2 := greetUseCase.Execute(command.GreetCommand{Name: ""})

	tf.RunTest("Empty input - IsError returns true", result2.IsError())
	if result2.IsError() {
		err := result2.ErrorInfo()
		tf.RunTest("Empty input - error kind is ValidationError",
			err.Kind == domerr.ValidationError)
		tf.RunTest("Empty input - error propagates from domain",
			strings.Contains(err.Message, "empty"))
	}

	// ========================================================================
	// Test: Component wiring - valid scenarios
	// ========================================================================

	validNames := []string{"Alice", "Bob Smith", "José García"}
	for _, name := range validNames {
		// Suppress stdout during test
		oldStdout = os.Stdout
		_, w, _ = os.Pipe()
		os.Stdout = w

		wf := adapter.NewConsoleWriter()
		uc := usecase.NewGreetUseCase(wf)
		res := uc.Execute(command.GreetCommand{Name: name})

		w.Close()
		os.Stdout = oldStdout

		tf.RunTest("Wiring test - "+name+" produces Ok", res.IsOk())
	}

	// ========================================================================
	// Test: Error propagation through layers
	// ========================================================================

	// Test that domain errors propagate correctly through application layer
	wf2 := adapter.NewConsoleWriter()
	uc := usecase.NewGreetUseCase(wf2)

	// Long name should produce validation error from domain
	longName := strings.Repeat("x", 200)
	res := uc.Execute(command.GreetCommand{Name: longName})
	tf.RunTest("Long name - error propagates from domain", res.IsError())
	if res.IsError() {
		tf.RunTest("Long name - error kind is ValidationError",
			res.ErrorInfo().Kind == domerr.ValidationError)
	}

	// Print summary
	tf.Summary(t)
}
