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
//   - Uses interfaces with generics for STATIC DISPATCH (compile-time resolution)
//
// Static Dispatch Pattern:
//  1. Application defines WriterPort interface (the contract)
//  2. Infrastructure implements a struct that satisfies WriterPort
//  3. Use case is generic over WriterPort: GreetUseCase[W WriterPort]
//  4. Bootstrap instantiates with concrete type: NewGreetUseCase[*ConsoleWriter](writer)
//  5. Compiler knows exact type → static dispatch, no vtable lookup
//
// Mapping to Ada:
//   - Ada: generic with function Write(...) return Result;
//   - Go: interface WriterPort + generic type parameter
//   - Both achieve: static dispatch, compile-time resolution, zero runtime overhead
//
// Usage:
//
//	import "github.com/abitofhelp/hybrid_app_go/application/port/outward"
//
//	type GreetUseCase[W outward.WriterPort] struct {
//	    writer W  // Concrete type known at compile time
//	}
//
//	func (uc *GreetUseCase[W]) Execute(ctx context.Context, cmd GreetCommand) domerr.Result[Unit] {
//	    result := uc.writer.Write(ctx, "Hello, World!")  // Static dispatch
//	    return result
//	}
package outward

import (
	"context"

	"github.com/abitofhelp/hybrid_app_go/application/model"
	domerr "github.com/abitofhelp/hybrid_app_go/domain/error"
)

// WriterPort is an output port contract for writing operations.
//
// This interface defines the contract between Application and Infrastructure layers.
// Any infrastructure adapter that wants to provide write output must:
//  1. Implement this interface with a concrete struct
//  2. Be injected into use cases via generic type parameter
//
// Static Dispatch:
//   - When used as generic type parameter: GreetUseCase[W WriterPort]
//   - The concrete type W is known at compile time
//   - Method calls are statically dispatched (devirtualized by compiler)
//   - Zero runtime overhead compared to dynamic interface dispatch
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
type WriterPort interface {
	Write(ctx context.Context, message string) domerr.Result[model.Unit]
}
