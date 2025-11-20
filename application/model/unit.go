// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
// Package: model
// Description: Unit type for Result with no meaningful value

// Package model provides application-level model types and DTOs.
//
// Architecture Notes:
//   - Part of the APPLICATION layer
//   - Application depends on Domain only
//   - Provides types that cross layer boundaries
//
// Usage:
//
//	import "github.com/abitofhelp/hybrid_app_go/application/model"
//
//	func Write(message string) mo.Result[model.Unit] {
//	    // Perform write operation
//	    return mo.Ok(model.UnitValue)
//	}
package model

// Unit represents "no meaningful value" in Result returns.
//
// This type is used for operations that return Result but have no meaningful
// success value. It represents "void" or "no value" in the Result monad,
// similar to () in Rust, void in C, or Unit in Scala.
//
// Design Notes:
//   - Used for operations with side effects (like console writes)
//   - Allows consistent mo.Result[Unit] return type instead of separate error handling
//   - Distinguishes success from failure even when there's no data to return
//
// Usage:
//
//	func Write(message string) mo.Result[Unit] {
//	    err := performWrite(message)
//	    if err != nil {
//	        return mo.Err[Unit](err)
//	    }
//	    return mo.Ok(UnitValue)
//	}
type Unit struct{}

// UnitValue is a singleton instance for convenience.
// Use this when you need to return a Unit value.
var UnitValue = Unit{}
