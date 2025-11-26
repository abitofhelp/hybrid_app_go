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
//   - Implements ports defined by Application layer (WriterPort interface)
//   - Depends on Application + Domain layers
//   - Converts exceptions/errors to Result types
//   - Handles all technical/platform-specific details
//   - Enables STATIC DISPATCH when used with generic use cases
//
// Static Dispatch Pattern:
//   - ConsoleWriter implements WriterPort interface
//   - GreetUseCase[*ConsoleWriter] is instantiated with concrete type
//   - Compiler knows exact type → devirtualizes method calls
//   - Equivalent to Ada's generic instantiation
//
// Design Pattern: Dependency Injection via io.Writer
//   - ConsoleWriter.w accepts any io.Writer for flexibility and testability
//   - NewConsoleWriter() is a convenience that uses os.Stdout
//   - Tests can inject bytes.Buffer to capture output
//   - Production can inject file writers, network writers, etc.
//
// Mapping to Ada:
//   - Ada: Infrastructure.Adapter.Console_Writer package with Write function
//   - Go: ConsoleWriter struct with Write method implementing WriterPort
//   - Both: Match the signature required by the output port
//
// Usage:
//
//	import "github.com/abitofhelp/hybrid_app_go/infrastructure/adapter"
//
//	// Production: write to console
//	writer := adapter.NewConsoleWriter()
//	result := writer.Write(ctx, "Hello, World!")
//
//	// Testing: capture output
//	var buf bytes.Buffer
//	writer := adapter.NewWriter(&buf)
//	result := writer.Write(ctx, "Hello, World!")
//	captured := buf.String()
//
//	// Static dispatch with generic use case
//	uc := usecase.NewGreetUseCase[*adapter.ConsoleWriter](writer)
package adapter

import (
	"context"
	"fmt"
	"io"
	"os"

	apperr "github.com/abitofhelp/hybrid_app_go/application/error"
	"github.com/abitofhelp/hybrid_app_go/application/model"
	domerr "github.com/abitofhelp/hybrid_app_go/domain/error"
)

// ConsoleWriter is an infrastructure adapter that writes to an io.Writer.
//
// This struct implements the WriterPort interface, enabling static dispatch
// when used as a type parameter in generic use cases.
//
// Static Dispatch:
//   - When used as GreetUseCase[*ConsoleWriter], compiler knows concrete type
//   - Method calls are devirtualized → zero runtime overhead
//   - Equivalent to Ada's generic package instantiation
//
// Design Pattern: Adapter
//   - Adapts io.Writer to WriterPort interface
//   - Converts I/O errors and panics to Result types
//   - Handles context cancellation
//
// Implements: outward.WriterPort
type ConsoleWriter struct {
	w io.Writer
}

// NewWriter creates a ConsoleWriter that writes to the provided io.Writer.
//
// This is the core adapter factory that demonstrates production-ready patterns:
//   - Accepts any io.Writer for flexibility (files, network, buffers)
//   - Enables testability by injecting test doubles (bytes.Buffer)
//   - Returns struct pointer for use with generic type parameters
//
// Static Dispatch Usage:
//   - writer := NewWriter(&buf)
//   - uc := usecase.NewGreetUseCase[*adapter.ConsoleWriter](writer)
//   - The uc.writer.Write() call is statically dispatched
//
// Dependency Inversion:
//   - Application defines the WriterPort interface it NEEDS
//   - Infrastructure provides this implementation
//   - Bootstrap wires them together via generic instantiation
//   - Application never depends on Infrastructure
//
// Example - Production:
//
//	writer := NewWriter(os.Stdout)
//	result := writer.Write(ctx, "Hello!")
//
// Example - Testing:
//
//	var buf bytes.Buffer
//	writer := NewWriter(&buf)
//	result := writer.Write(ctx, "Hello!")
//	assert.Equal(t, "Hello!\n", buf.String())
//
// Example - File Output:
//
//	file, _ := os.Create("output.txt")
//	defer file.Close()
//	writer := NewWriter(file)
//	result := writer.Write(ctx, "Hello!")
func NewWriter(w io.Writer) *ConsoleWriter {
	return &ConsoleWriter{w: w}
}

// Write writes the message to the underlying io.Writer.
//
// This method implements the WriterPort interface, enabling static dispatch
// when ConsoleWriter is used as a type parameter in generic use cases.
//
// Production-Ready Patterns:
//   - Recovers from panics and converts to InfrastructureError
//   - Checks context cancellation before I/O
//   - Maps all io.Writer errors to InfrastructureError
//   - Always returns Result (never panics across boundary)
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
//
// Contract:
//   - ctx parameter carries cancellation and deadline signals
//   - message can be any string
//   - Returns Ok(Unit) on success
//   - Returns Err(InfrastructureError) on I/O failure, panic, or cancellation
//   - Never panics (panics are caught and converted to Err)
func (cw *ConsoleWriter) Write(ctx context.Context, message string) (result domerr.Result[model.Unit]) {
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
	_, err := fmt.Fprintln(cw.w, message)
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

// NewConsoleWriter creates a ConsoleWriter that writes to standard output.
//
// This is a convenience function that wraps NewWriter with os.Stdout.
// Use this for production CLI applications.
//
// Static Dispatch Usage:
//   - writer := adapter.NewConsoleWriter()
//   - uc := usecase.NewGreetUseCase[*adapter.ConsoleWriter](writer)
//
// For testing, use NewWriter with a bytes.Buffer instead to capture output.
//
// Usage:
//
//	writer := adapter.NewConsoleWriter()
//	result := writer.Write(ctx, "Hello, World!")
func NewConsoleWriter() *ConsoleWriter {
	return NewWriter(os.Stdout)
}

// NewStderrWriter creates a ConsoleWriter that writes to standard error.
//
// Use this for error messages or diagnostic output that should go to stderr.
//
// Usage:
//
//	errWriter := adapter.NewStderrWriter()
//	result := errWriter.Write(ctx, "Error: something went wrong")
func NewStderrWriter() *ConsoleWriter {
	return NewWriter(os.Stderr)
}
