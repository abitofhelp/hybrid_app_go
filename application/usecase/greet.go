// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
// Package: usecase
// Description: Greet use case orchestration

// Package usecase provides application use cases - orchestration logic that
// coordinates domain objects to fulfill specific business operations.
//
// Architecture Notes:
//   - Part of the APPLICATION layer
//   - Use case = application business logic orchestration
//   - Coordinates domain objects
//   - Depends on output ports defined in application layer
//   - Never imports infrastructure layer
//   - Uses GENERICS for STATIC DISPATCH (compile-time resolution)
//
// Static Dispatch Pattern:
//   - GreetUseCase[W WriterPort] is generic over the writer type
//   - At instantiation, concrete type is known: NewGreetUseCase[*ConsoleWriter](writer)
//   - Compiler devirtualizes method calls → zero runtime overhead
//   - Equivalent to Ada's generic package instantiation
//
// Dependency Flow (all pointing INWARD toward Domain):
//   - GreetUseCase[W] -> domain.Person (coordinates)
//   - GreetUseCase[W] -> application.port.outbound.WriterPort (interface constraint)
//   - infrastructure.ConsoleWriter -> WriterPort (implements)
//   - Bootstrap instantiates GreetUseCase[*ConsoleWriter]
//
// Mapping to Ada:
//   - Ada: generic with function Writer(...) return Result; package Application.Usecase.Greet
//   - Go: type GreetUseCase[W WriterPort] struct { writer W }
//   - Both achieve: static dispatch, compile-time resolution
//
// Usage:
//
//	import "github.com/abitofhelp/hybrid_app_go/application/usecase"
//
//	// Bootstrap instantiates with concrete type
//	consoleWriter := &adapter.ConsoleWriter{...}
//	uc := usecase.NewGreetUseCase[*adapter.ConsoleWriter](consoleWriter)
//
//	// Use case Execute is statically dispatched
//	result := uc.Execute(ctx, greetCommand)
package usecase

import (
	"context"

	"github.com/abitofhelp/hybrid_app_go/application/command"
	"github.com/abitofhelp/hybrid_app_go/application/model"
	"github.com/abitofhelp/hybrid_app_go/application/port/outbound"
	domerr "github.com/abitofhelp/hybrid_app_go/domain/error"
	"github.com/abitofhelp/hybrid_app_go/domain/valueobject"
)

// GreetUseCase orchestrates the greeting workflow.
//
// This use case demonstrates application-layer orchestration with static dispatch:
//  1. Generic over WriterPort: GreetUseCase[W WriterPort]
//  2. Receives command DTO from presentation layer
//  3. Validates input using domain layer (Person)
//  4. Generates greeting message (domain logic)
//  5. Writes output via infrastructure port (statically dispatched)
//  6. Returns Result to presentation layer
//
// Static Dispatch:
//   - Type parameter W is constrained to WriterPort interface
//   - At instantiation, concrete type replaces W
//   - Compiler knows exact type → method calls are devirtualized
//   - Zero runtime overhead (no vtable lookup)
//
// Design Pattern: Generic Use Case (matching Ada's generic package)
//   - Single responsibility (one business operation)
//   - Coordinates domain objects
//   - Generic over port abstraction (static dispatch)
//   - Returns Result for functional error handling
//
// Implements: inbound.GreetPort interface
type GreetUseCase[W outbound.WriterPort] struct {
	writer W
}

// NewGreetUseCase creates a new GreetUseCase with injected dependencies.
//
// Static Dependency Injection Pattern:
//   - Type parameter W specifies the concrete writer type
//   - Writer instance is injected via constructor
//   - Use case doesn't know the implementation details
//   - But compiler knows the concrete type for static dispatch
//   - Bootstrap wires them together: NewGreetUseCase[*ConsoleWriter](writer)
//
// Mapping to Ada:
//   - Ada: package Greet_UC is new Application.Usecase.Greet(Writer => Console_Writer.Write);
//   - Go: uc := NewGreetUseCase[*adapter.ConsoleWriter](consoleWriter)
func NewGreetUseCase[W outbound.WriterPort](writer W) *GreetUseCase[W] {
	return &GreetUseCase[W]{writer: writer}
}

// Execute runs the greeting use case.
//
// Orchestration workflow:
//  1. Extract name from GreetCommand DTO
//  2. Validate and create Person from name (domain validation)
//  3. Generate greeting message (application-level formatting)
//  4. Write greeting to console via output port (STATIC DISPATCH)
//  5. Propagate any errors via railway-oriented programming
//
// Railway-Oriented Programming:
//   - Uses AndThenTo for functional composition across Result types
//   - Person validation failure short-circuits to error track
//   - Writer failure propagates error to caller
//   - Success returns Ok(Unit)
//
// Static Dispatch:
//   - uc.writer.Write() is statically dispatched because W is concrete at instantiation
//   - Compiler knows exact implementation → no vtable lookup
//   - Equivalent to Ada's generic instantiation with compile-time resolution
//
// Parameters:
//   - ctx: Context for cancellation and deadlines (passed to infrastructure)
//   - cmd: GreetCommand DTO crossing presentation -> application boundary
//
// Error scenarios:
//   - ValidationError: Invalid person name (empty, too long)
//   - InfrastructureError: Console write failure or context cancellation
//
// Contract:
//   - Pre: ctx is non-nil (use context.Background() if no cancellation needed)
//   - Pre: cmd can be any GreetCommand (validation happens inside)
//   - Post: Returns Ok(Unit) if greeting succeeded
//   - Post: Returns Err(ValidationError) if name validation failed
//   - Post: Returns Err(InfrastructureError) if write failed or ctx cancelled
func (uc *GreetUseCase[W]) Execute(ctx context.Context, cmd command.GreetCommand) domerr.Result[model.Unit] {
	// Step 1: Validate and create Person from name (domain validation)
	personResult := valueobject.CreatePerson(cmd.GetName())

	// Step 2-4: Chain operations using railway-oriented programming
	// AndThenTo enables cross-type chaining: Result[Person] → Result[Unit]
	// If personResult is Error, error propagates without calling the lambda
	// If personResult is Ok, lambda executes and may return Ok or Error
	return domerr.AndThenTo(personResult, func(person valueobject.Person) domerr.Result[model.Unit] {
		// Application-level greeting format (orchestration, not domain logic)
		message := formatGreeting(person.GetName())

		// Write to console via output port (STATIC DISPATCH)
		return uc.writer.Write(ctx, message)
	})
}

// formatGreeting creates the greeting message.
// This is application-level formatting logic, not domain logic.
// The format "Hello, <name>!" is an application decision.
func formatGreeting(name string) string {
	return "Hello, " + name + "!"
}
