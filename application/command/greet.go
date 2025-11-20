// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
// Package: command
// Description: DTOs for application use cases

// Package command provides Data Transfer Objects (DTOs) that cross
// architectural boundaries. DTOs are simple data structures with no business
// logic, used to transfer data between layers.
//
// Architecture Notes:
//   - Part of the APPLICATION layer (use cases and contracts)
//   - DTOs should be simple, serializable data structures
//   - No business logic in DTOs
//   - DTOs are different from domain entities
//   - Crosses boundary: Presentation -> Application
//
// Usage:
//
//	import "github.com/abitofhelp/hybrid_app_go/application/command"
//
//	cmd := command.NewGreetCommand("Alice")
//	result := greetUseCase.Execute(cmd)
package command

// GreetCommand is a Data Transfer Object for the greet use case.
//
// This DTO crosses the presentation -> application boundary. It may carry
// invalid data; the domain layer is responsible for validating the name
// and returning appropriate Result errors.
//
// Design Notes:
//   - Simple data structure (no methods except accessors)
//   - No validation logic (validation is in domain layer)
//   - Separates external API from internal domain model
type GreetCommand struct {
	Name string
}

// NewGreetCommand creates a new GreetCommand DTO from a name string.
//
// This function does not perform validation; it simply packages the raw
// input. Validation is performed in domain.Person.CreatePerson via Result.
func NewGreetCommand(name string) GreetCommand {
	return GreetCommand{Name: name}
}

// GetName extracts the name as a string.
func (c GreetCommand) GetName() string {
	return c.Name
}
