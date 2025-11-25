// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
// Package: adapter
// Description: Console output adapter

// Package adapter provides concrete implementations of application ports
// (adapter pattern). Adapters implement interfaces defined by the application
// layer, converting between the application's needs and external systems.
//
// Architecture Notes:
//   - Part of the INFRASTRUCTURE layer (driven/secondary adapters)
//   - Implements ports defined by Application layer
//   - Depends on Application + Domain layers
//   - Converts exceptions/errors to Result types
//   - Handles all technical/platform-specific details
//
// Design Pattern: Dependency Injection via io.Writer
//   - NewWriter accepts any io.Writer for flexibility and testability
//   - NewConsoleWriter is a convenience that uses os.Stdout
//   - Tests can inject bytes.Buffer to capture output
//   - Production can inject file writers, network writers, etc.
//
// Usage:
//
//	import "github.com/abitofhelp/hybrid_app_go/infrastructure/adapter"
//
//	// Production: write to console
//	writer := adapter.NewConsoleWriter()
//	result := writer(ctx, "Hello, World!")
//
//	// Testing: capture output
//	var buf bytes.Buffer
//	writer := adapter.NewWriter(&buf)
//	result := writer(ctx, "Hello, World!")
//	captured := buf.String()
package adapter

import (
	"context"
	"fmt"
	"io"
	"os"

	apperr "github.com/abitofhelp/hybrid_app_go/application/error"
	"github.com/abitofhelp/hybrid_app_go/application/model"
	"github.com/abitofhelp/hybrid_app_go/application/port/outward"
	domerr "github.com/abitofhelp/hybrid_app_go/domain/error"
)

// NewWriter creates a WriterFunc that writes to the provided io.Writer.
//
// This is the core adapter factory that demonstrates production-ready patterns:
//   - Accepts any io.Writer for flexibility (files, network, buffers)
//   - Enables testability by injecting test doubles (bytes.Buffer)
//   - Properly maps all I/O errors to domain InfrastructureError
//
// Dependency Inversion:
//   - Application defines the WriterFunc interface it NEEDS
//   - Infrastructure provides this implementation
//   - Bootstrap wires them together
//   - Application never depends on Infrastructure
//
// Context Handling:
//   - Checks ctx.Done() before performing I/O
//   - Returns InfrastructureError if context is cancelled
//   - Enables graceful shutdown and timeout support
//
// Error Handling:
//   - Recovers from panics and converts to InfrastructureError
//   - Maps all io.Writer errors to InfrastructureError
//   - Includes original error message for debugging
//   - Always returns Result (never panics across boundary)
//
// Example - Production:
//
//	writer := NewWriter(os.Stdout)
//	result := writer(ctx, "Hello!")
//
// Example - Testing:
//
//	var buf bytes.Buffer
//	writer := NewWriter(&buf)
//	result := writer(ctx, "Hello!")
//	assert.Equal(t, "Hello!\n", buf.String())
//
// Example - File Output:
//
//	file, _ := os.Create("output.txt")
//	defer file.Close()
//	writer := NewWriter(file)
//	result := writer(ctx, "Hello!")
func NewWriter(w io.Writer) outward.WriterFunc {
	return func(ctx context.Context, message string) (result domerr.Result[model.Unit]) {
		// Recover from any panics and convert to InfrastructureError
		// This ensures NO panics escape across the infrastructure boundary
		// Pattern: Infrastructure adapters are the "exception boundary" where
		// all panics/exceptions must be caught and converted to Result errors
		defer func() {
			if r := recover(); r != nil {
				result = domerr.Err[model.Unit](apperr.NewInfrastructureError(
					fmt.Sprintf("write panicked: %v", r)))
			}
		}()

		// Check for context cancellation before I/O
		// This is important for long-running operations or network writers
		select {
		case <-ctx.Done():
			return domerr.Err[model.Unit](apperr.NewInfrastructureError(
				fmt.Sprintf("write cancelled: %v", ctx.Err())))
		default:
			// Context is still active, proceed with I/O
		}

		// Perform the I/O operation using the injected writer
		// fmt.Fprintln handles the newline and returns any write errors
		_, err := fmt.Fprintln(w, message)
		if err != nil {
			// Map the I/O error to a domain InfrastructureError
			// This keeps infrastructure concerns (specific error types)
			// from leaking into application/domain layers
			return domerr.Err[model.Unit](apperr.NewInfrastructureError(
				fmt.Sprintf("write failed: %v", err)))
		}

		// Success case - return Unit to indicate completion
		return domerr.Ok(model.UnitValue)
	}
}

// NewConsoleWriter creates a WriterFunc that writes to standard output.
//
// This is a convenience function that wraps NewWriter with os.Stdout.
// Use this for production CLI applications.
//
// For testing, use NewWriter with a bytes.Buffer instead to capture output.
//
// Usage:
//
//	writer := adapter.NewConsoleWriter()
//	result := writer(ctx, "Hello, World!")
func NewConsoleWriter() outward.WriterFunc {
	return NewWriter(os.Stdout)
}

// NewStderrWriter creates a WriterFunc that writes to standard error.
//
// Use this for error messages or diagnostic output that should go to stderr.
//
// Usage:
//
//	errWriter := adapter.NewStderrWriter()
//	result := errWriter(ctx, "Error: something went wrong")
func NewStderrWriter() outward.WriterFunc {
	return NewWriter(os.Stderr)
}
