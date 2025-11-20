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
// Usage:
//
//	import "github.com/abitofhelp/hybrid_app_go/infrastructure/adapter"
//
//	writer := adapter.NewConsoleWriter()
//	result := writer("Hello, World!")
package adapter

import (
	"fmt"

	apperr "github.com/abitofhelp/hybrid_app_go/application/error"
	"github.com/abitofhelp/hybrid_app_go/application/model"
	"github.com/abitofhelp/hybrid_app_go/application/port/outward"
	domerr "github.com/abitofhelp/hybrid_app_go/domain/error"
)

// NewConsoleWriter creates a WriterFunc that writes to standard output.
//
// This adapter implements the outward.WriterFunc port interface defined
// by the application layer. It conforms to the interface without the
// application knowing about this implementation.
//
// Dependency Inversion:
//   - Application defines the WriterFunc interface it NEEDS
//   - Infrastructure provides this implementation
//   - Bootstrap wires them together
//   - Application never depends on Infrastructure
//
// Error Handling:
//   - Captures any panics and converts to Err
//   - Maps I/O errors to InfrastructureError
//   - Always returns Result (never panics across boundary)
//
// Returns:
//   - outward.WriterFunc that can be injected into use cases
func NewConsoleWriter() outward.WriterFunc {
	return func(message string) domerr.Result[model.Unit] {
		// Perform the I/O operation
		_, err := fmt.Println(message)
		if err != nil {
			// Convert I/O error to domain error
			return domerr.Err[model.Unit](apperr.NewInfrastructureError(
				fmt.Sprintf("Console write failed: %v", err)))
		}

		// Success case
		return domerr.Ok(model.UnitValue)
	}
}

// Write is an alternative function-style adapter (not method-based).
//
// This function directly implements the outward.WriterFunc signature and
// can be used in places where function references are preferred.
//
// Usage:
//
//	useCase := NewGreetUseCase(adapter.Write)
func Write(message string) domerr.Result[model.Unit] {
	_, err := fmt.Println(message)
	if err != nil {
		return domerr.Err[model.Unit](apperr.NewInfrastructureError(
			fmt.Sprintf("Console write failed: %v", err)))
	}
	return domerr.Ok(model.UnitValue)
}
