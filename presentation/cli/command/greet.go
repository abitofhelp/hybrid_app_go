// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
// Package: command
// Description: CLI command for greet use case

// Package command provides CLI command handlers for the presentation layer.
// Command handlers are responsible for UI concerns: parsing arguments,
// displaying output, and mapping results to exit codes.
//
// Architecture Notes:
//   - Part of the PRESENTATION layer (driving/primary adapters)
//   - Handles user interface concerns (CLI args, output formatting)
//   - Calls APPLICATION layer use cases (through input ports)
//   - Does NOT depend on Infrastructure or Domain directly
//   - Does NOT contain business logic (delegates to use case)
//
// Dependency Flow (all pointing INWARD):
//   - GreetCommand -> application.GreetUseCase (calls use case)
//   - GreetCommand -> application.Error (re-exported from domain)
//   - GreetCommand -> application.Command (DTOs)
//
// Critical Architectural Rule:
//   - Presentation MUST NOT import domain/* packages
//   - Presentation MUST use application/error re-exports
//   - This prevents tight coupling between UI and business logic
//
// Usage:
//
//	import "github.com/abitofhelp/hybrid_app_go/presentation/cli/command"
//
//	cmd := command.NewGreetCommand(greetUseCase)
//	exitCode := cmd.Run(args)
package command

import (
	"context"
	"fmt"
	"os"

	"github.com/abitofhelp/hybrid_app_go/application/command"
	apperr "github.com/abitofhelp/hybrid_app_go/application/error"
	"github.com/abitofhelp/hybrid_app_go/application/model"
)

// GreetUseCaseFunc is the input port contract for the greet use case.
//
// This type defines the interface between Presentation and Application layers.
// The use case is injected via this function type, enabling dependency inversion.
//
// Pattern: Input Port (Driving Adapter calls Application)
//   - Presentation defines what it needs (this function signature)
//   - Application provides implementation (use case Execute method)
//   - Bootstrap wires them together
//
// Context Usage:
//   - ctx carries cancellation signals from CLI (e.g., Ctrl+C handling)
//   - For CLI apps, context.Background() is typically used
//   - Enables future support for timeouts and graceful shutdown
type GreetUseCaseFunc func(ctx context.Context, cmd command.GreetCommand) apperr.Result[model.Unit]

// GreetCommand is a CLI command handler for the greet use case.
//
// This command demonstrates presentation-layer concerns:
//  1. Parse command-line arguments
//  2. Create application DTOs
//  3. Call use case
//  4. Display results to user
//  5. Map results to exit codes
//
// Design Pattern: Command Handler
//   - Single responsibility (one CLI command)
//   - Coordinates UI concerns
//   - Depends on abstractions (use case func), not implementations
//   - Returns exit code for shell
type GreetCommand struct {
	useCase GreetUseCaseFunc
}

// NewGreetCommand creates a new GreetCommand with injected use case.
//
// Dependency Injection Pattern:
//   - Use case function is injected via constructor
//   - Command doesn't know the implementation
//   - Application provides the implementation
//   - Bootstrap wires them together
func NewGreetCommand(useCase GreetUseCaseFunc) *GreetCommand {
	return &GreetCommand{useCase: useCase}
}

// Run executes the CLI command logic.
//
// Responsibilities:
//  1. Parse command-line arguments
//  2. Extract the name parameter
//  3. Create GreetCommand DTO
//  4. Call the use case with context and DTO
//  5. Handle the result and display appropriate messages
//  6. Return exit code (0 = success, non-zero = error)
//
// CLI Usage: greeter <name>
// Example: ./greeter Alice
//
// This is where presentation concerns live:
//   - CLI argument parsing
//   - Context creation (for cancellation support)
//   - User-facing error messages
//   - Exit code mapping
//
// Contract:
//   - Pre: args can be any slice (validation happens inside)
//   - Post: Returns 0 if greeting succeeded
//   - Post: Returns 1 if validation or infrastructure error occurred
//   - Post: Displays error message to stderr on failure
func (c *GreetCommand) Run(args []string) int {
	// Check if user provided exactly one argument (the name)
	if len(args) != 2 { // args[0] is program name, args[1] is the name
		programName := args[0]
		fmt.Fprintf(os.Stderr, "Usage: %s <name>\n", programName)
		fmt.Fprintf(os.Stderr, "Example: %s Alice\n", programName)
		return 1 // Exit code 1 indicates error
	}

	// Extract the name from command-line arguments
	name := args[1]

	// Create DTO for crossing presentation -> application boundary
	cmd := command.NewGreetCommand(name)

	// Create context for the request
	// For CLI apps, we use Background context. Future enhancement could
	// add signal handling for graceful shutdown on Ctrl+C.
	ctx := context.Background()

	// Call the use case (injected via constructor)
	// This is the key architectural boundary:
	// Presentation -> Application (through input port)
	result := c.useCase(ctx, cmd)

	// Handle the result from the use case
	if result.IsOk() {
		// Success! Greeting was displayed via console port
		// Use case already wrote to console, just exit cleanly
		return 0 // Exit code 0 indicates success
	}

	// Use case failed - display error to user
	domErr := result.ErrorInfo()

	// Display user-friendly error message
	fmt.Fprintf(os.Stderr, "Error: %s\n", domErr.Message)

	// Add detailed error handling based on ErrorKind
	// Note: We use apperr types here but the error comes through domain layer
	switch domErr.Kind {
	case apperr.ValidationError:
		fmt.Fprintln(os.Stderr, "Please provide a valid name.")

	case apperr.InfrastructureError:
		fmt.Fprintln(os.Stderr, "A system error occurred.")
	}

	return 1 // Exit code 1 indicates error
}
