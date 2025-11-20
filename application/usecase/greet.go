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
//
// Dependency Flow (all pointing INWARD toward Domain):
//   - GreetUseCase -> domain.Person
//   - GreetUseCase -> application.port.outward.WriterFunc (interface)
//   - infrastructure.ConsoleWriter -> WriterFunc (implements)
//   - Bootstrap wires them together
//
// Usage:
//
//	import "github.com/abitofhelp/hybrid_app_go/application/usecase"
//
//	uc := usecase.NewGreetUseCase(consoleWriter)
//	result := uc.Execute(greetCommand)
package usecase

import (
	"github.com/abitofhelp/hybrid_app_go/application/command"
	"github.com/abitofhelp/hybrid_app_go/application/model"
	"github.com/abitofhelp/hybrid_app_go/application/port/outward"
	domerr "github.com/abitofhelp/hybrid_app_go/domain/error"
	"github.com/abitofhelp/hybrid_app_go/domain/valueobject"
)

// GreetUseCase orchestrates the greeting workflow.
//
// This use case demonstrates application-layer orchestration:
//  1. Receives command DTO from presentation layer
//  2. Validates input using domain layer (Person)
//  3. Generates greeting message (domain logic)
//  4. Writes output via infrastructure port
//  5. Returns Result to presentation layer
//
// Design Pattern: Use Case
//   - Single responsibility (one business operation)
//   - Coordinates domain objects
//   - Depends on abstractions (ports), not implementations
//   - Returns Result for functional error handling
type GreetUseCase struct {
	writer outward.WriterFunc
}

// NewGreetUseCase creates a new GreetUseCase with injected dependencies.
//
// Dependency Injection Pattern:
//   - Writer function is injected via constructor
//   - Use case doesn't know the implementation
//   - Infrastructure provides the implementation
//   - Bootstrap wires them together
func NewGreetUseCase(writer outward.WriterFunc) *GreetUseCase {
	return &GreetUseCase{writer: writer}
}

// Execute runs the greeting use case.
//
// Orchestration workflow:
//  1. Extract name from GreetCommand DTO
//  2. Validate and create Person from name
//  3. Generate greeting message from Person
//  4. Write greeting to console via output port
//  5. Propagate any errors up to caller
//
// Input: GreetCommand DTO crossing presentation -> application boundary
//
// Error scenarios:
//   - ValidationError: Invalid person name (empty, too long)
//   - InfrastructureError: Console write failure (rare, but possible)
//
// Contract:
//   - Pre: cmd can be any GreetCommand (validation happens inside)
//   - Post: Returns Ok(Unit) if greeting succeeded
//   - Post: Returns Err(ValidationError) if name validation failed
//   - Post: Returns Err(InfrastructureError) if write failed
func (uc *GreetUseCase) Execute(cmd command.GreetCommand) domerr.Result[model.Unit] {
	// Step 1: Extract name from DTO
	name := cmd.GetName()

	// Step 2: Validate and create Person from name (domain validation)
	personResult := valueobject.CreatePerson(name)

	// Check if person creation failed (railway-oriented programming)
	if personResult.IsError() {
		// Propagate validation error to caller
		domErr := personResult.ErrorInfo()
		return domerr.Err[model.Unit](domErr)
	}

	// Extract validated Person
	person := personResult.Value()

	// Step 3: Generate greeting message from Person (pure domain logic)
	message := person.GreetingMessage()

	// Step 4: Write to console via output port (injected dependency)
	writeResult := uc.writer(message)

	// Step 5: Propagate result (success or failure) to caller
	return writeResult
}
