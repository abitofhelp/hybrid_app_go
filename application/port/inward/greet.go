// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
// Package: inward
// Description: Input port for greet use case

// Package inward defines input (driving/primary) ports - interfaces that
// the application layer EXPOSES and the presentation layer CALLS.
//
// Architecture Notes:
//   - Part of the APPLICATION layer
//   - Application defines the interface it PROVIDES
//   - Presentation layer CALLS through this interface
//   - Enables dependency inversion: Presentation depends on abstraction, not concrete use case
//   - Uses interfaces with generics for STATIC DISPATCH (compile-time resolution)
//
// Static Dispatch Pattern:
//  1. Application defines GreetPort interface (the contract)
//  2. Application implements GreetUseCase[W WriterPort] satisfying GreetPort
//  3. Presentation command is generic over GreetPort: GreetCommand[UC GreetPort]
//  4. Bootstrap instantiates with concrete type
//  5. Compiler knows exact type → static dispatch, no vtable lookup
//
// Mapping to Ada:
//   - Ada: generic with function Execute_Greet_UseCase(...) return Result;
//   - Go: interface GreetPort + generic type parameter
//   - Both achieve: static dispatch, compile-time resolution, zero runtime overhead
//
// Usage:
//
//	import "github.com/abitofhelp/hybrid_app_go/application/port/inward"
//
//	type GreetCommand[UC inward.GreetPort] struct {
//	    useCase UC  // Concrete type known at compile time
//	}
//
//	func (c *GreetCommand[UC]) Run(args []string) int {
//	    result := c.useCase.Execute(ctx, cmd)  // Static dispatch
//	    // ...
//	}
package inward

import (
	"context"

	"github.com/abitofhelp/hybrid_app_go/application/command"
	"github.com/abitofhelp/hybrid_app_go/application/model"
	domerr "github.com/abitofhelp/hybrid_app_go/domain/error"
)

// GreetPort is an input port contract for the greet use case.
//
// This interface defines the contract between Presentation and Application layers.
// Any use case that wants to provide greet functionality must:
//  1. Implement this interface (GreetUseCase does)
//  2. Be injected into presentation commands via generic type parameter
//
// Static Dispatch:
//   - When used as generic type parameter: GreetCommand[UC GreetPort]
//   - The concrete type UC is known at compile time
//   - Method calls are statically dispatched (devirtualized by compiler)
//   - Zero runtime overhead compared to dynamic interface dispatch
//
// Context Usage:
//   - ctx carries cancellation signals and deadlines from caller
//   - For CLI apps, context.Background() is typically used
//   - For HTTP handlers, request context flows through
//
// Contract:
//   - ctx parameter carries cancellation and deadline signals
//   - cmd is a GreetCommand DTO carrying the name to greet
//   - Returns Ok(Unit) on success (greeting was displayed)
//   - Returns Err(ValidationError) if name validation failed
//   - Returns Err(InfrastructureError) if write operation failed
type GreetPort interface {
	Execute(ctx context.Context, cmd command.GreetCommand) domerr.Result[model.Unit]
}
