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
//   - Uses GENERICS for STATIC DISPATCH (compile-time resolution)
//
// Static Dispatch Pattern:
//   - GreetCommand[UC GreetPort] is generic over the use case type
//   - At instantiation, concrete type is known: NewGreetCommand[*GreetUseCase[*ConsoleWriter]](uc)
//   - Compiler devirtualizes method calls → zero runtime overhead
//   - Equivalent to Ada's generic package instantiation
//
// Dependency Flow (all pointing INWARD):
//   - GreetCommand[UC] -> application.port.inbound.GreetPort (interface constraint)
//   - GreetCommand[UC] -> application.Error (re-exported from domain)
//   - GreetCommand[UC] -> application.Command (DTOs)
//
// Mapping to Ada:
//   - Ada: generic with function Execute_Greet_UseCase(...) return Result; package Presentation.CLI.Command.Greet
//   - Go: type GreetCommand[UC GreetPort] struct { useCase UC }
//   - Both achieve: static dispatch, compile-time resolution
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
//	// Bootstrap instantiates with concrete type
//	uc := usecase.NewGreetUseCase[*adapter.ConsoleWriter](writer)
//	cmd := command.NewGreetCommand[*usecase.GreetUseCase[*adapter.ConsoleWriter]](uc)
//	exitCode := cmd.Run(args)
package command

import (
	"context"
	"fmt"
	"os"

	"github.com/abitofhelp/hybrid_app_go/application/command"
	apperr "github.com/abitofhelp/hybrid_app_go/application/error"
	"github.com/abitofhelp/hybrid_app_go/application/port/inbound"
)

// GreetCommand is a CLI command handler for the greet use case.
//
// This command demonstrates presentation-layer concerns with static dispatch:
//  1. Generic over GreetPort: GreetCommand[UC GreetPort]
//  2. Parse command-line arguments
//  3. Create application DTOs
//  4. Call use case (statically dispatched)
//  5. Display results to user
//  6. Map results to exit codes
//
// Static Dispatch:
//   - Type parameter UC is constrained to GreetPort interface
//   - At instantiation, concrete type replaces UC
//   - Compiler knows exact type → method calls are devirtualized
//   - Zero runtime overhead (no vtable lookup)
//
// Design Pattern: Generic Command Handler (matching Ada's generic package)
//   - Single responsibility (one CLI command)
//   - Coordinates UI concerns
//   - Generic over port abstraction (static dispatch)
//   - Returns exit code for shell
type GreetCommand[UC inbound.GreetPort] struct {
	useCase UC
}

// NewGreetCommand creates a new GreetCommand with injected use case.
//
// Static Dependency Injection Pattern:
//   - Type parameter UC specifies the concrete use case type
//   - Use case instance is injected via constructor
//   - Command doesn't know the implementation details
//   - But compiler knows the concrete type for static dispatch
//   - Bootstrap wires them together
//
// Mapping to Ada:
//   - Ada: package Greet_Command_Instance is new Presentation.CLI.Command.Greet(Execute_Greet_UseCase => Greet_UC.Execute);
//   - Go: cmd := NewGreetCommand[*usecase.GreetUseCase[*adapter.ConsoleWriter]](uc)
func NewGreetCommand[UC inbound.GreetPort](useCase UC) *GreetCommand[UC] {
	return &GreetCommand[UC]{useCase: useCase}
}

// Run executes the CLI command logic.
//
// Responsibilities:
//  1. Parse command-line arguments
//  2. Extract the name parameter
//  3. Create GreetCommand DTO
//  4. Call the use case with context and DTO (STATIC DISPATCH)
//  5. Handle the result and display appropriate messages
//  6. Return exit code (0 = success, non-zero = error)
//
// Static Dispatch:
//   - c.useCase.Execute() is statically dispatched because UC is concrete at instantiation
//   - Compiler knows exact implementation → no vtable lookup
//   - Equivalent to Ada's generic instantiation with compile-time resolution
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
func (c *GreetCommand[UC]) Run(args []string) int {
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

	// Call the use case (STATIC DISPATCH)
	// The useCase.Execute() call is statically dispatched because UC is a
	// concrete type at instantiation time.
	// This is the key architectural boundary:
	// Presentation -> Application (through input port)
	result := c.useCase.Execute(ctx, cmd)

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
