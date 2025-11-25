// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
// Package: outward
// Description: Output port for writing operations

// Package outward defines output (driven/secondary) ports - interfaces that
// the application layer NEEDS and the infrastructure layer IMPLEMENTS.
//
// Architecture Notes:
//   - Part of the APPLICATION layer
//   - Application defines the interface it NEEDS
//   - Infrastructure layer CONFORMS to this interface
//   - This inverts the dependency: Infrastructure -> Application (not Application -> Infrastructure)
//   - Uses function types for Go's lightweight dependency injection
//
// Port Flow:
//  1. Application defines WriterFunc signature
//  2. Infrastructure implements function matching this signature
//  3. Bootstrap injects infrastructure's implementation into use case
//  4. Use case calls the function without knowing the implementation
//
// Usage:
//
//	import "github.com/abitofhelp/hybrid_app_go/application/port/outward"
//
//	type GreetUseCase struct {
//	    writer outward.WriterFunc
//	}
//
//	func (uc *GreetUseCase) Execute(cmd GreetCommand) domerr.Result[Unit] {
//	    result := uc.writer("Hello, World!")
//	    return result
//	}
package outward

import (
	"context"

	"github.com/abitofhelp/hybrid_app_go/application/model"
	domerr "github.com/abitofhelp/hybrid_app_go/domain/error"
)

// WriterFunc is an output port contract for writing operations.
//
// Any infrastructure adapter that wants to provide write output must:
//  1. Implement a function matching this signature
//  2. Be injected into use cases that need write capabilities
//
// The function accepts a context and message string and returns domerr.Result[Unit]:
//   - Ok(Unit) if write succeeded
//   - Err(error) if write failed (with domain ErrorType)
//
// Context Usage:
//   - ctx carries cancellation signals and deadlines from caller
//   - Implementations SHOULD check ctx.Done() before expensive operations
//   - For CLI apps, context.Background() is typically used
//   - For HTTP handlers, request context flows through
//
// Contract:
//   - ctx parameter carries cancellation and deadline signals
//   - Message parameter can be any string (no length restrictions at this boundary)
//   - Returns Ok(Unit) on success
//   - Returns Err with InfrastructureError on I/O failure or context cancellation
//   - Must not panic (convert panics to Err if needed)
type WriterFunc func(ctx context.Context, message string) domerr.Result[model.Unit]
