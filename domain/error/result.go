// ===========================================================================
// result.go
// ===========================================================================
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
// SPDX-License-Identifier: BSD-3-Clause
//
// Purpose:
//   Result[T] monad for railway-oriented programming. This is the core
//   functional error handling primitive for the entire application.
//
// Architecture Notes:
//   - Generic over success type T
//   - Uses Domain.Error.ErrorType for all errors
//   - Pure domain implementation (ZERO external dependencies)
//   - Enables functional composition and error propagation
//
// Design Pattern:
//   Railway-Oriented Programming:
//   - Ok track: Successful computation continues
//   - Error track: Error propagates (short-circuit)
//   - Forces explicit error handling at compile time
// ===========================================================================

// Package error provides domain error types and Result monad for error handling.
package error

// Result represents either a successful value of type T or an error.
// This is the core functional error handling type.
//
// States:
//   - Ok: Contains a value of type T
//   - Err: Contains an ErrorType
//
// Usage:
//
//	result := Ok[string]("success")
//	if result.IsOk() {
//	    value := result.Value()
//	}
type Result[T any] struct {
	value T
	err   ErrorType
	isOk  bool
}

// ============================================================================
// Constructors
// ============================================================================

// Ok creates a Result containing a successful value.
//
// Example:
//
//	result := Ok[int](42)
func Ok[T any](value T) Result[T] {
	return Result[T]{
		value: value,
		isOk:  true,
	}
}

// Err creates a Result containing an error.
//
// Example:
//
//	result := Err[int](NewValidationError("invalid input"))
func Err[T any](err ErrorType) Result[T] {
	return Result[T]{
		err:  err,
		isOk: false,
	}
}

// ============================================================================
// Query functions
// ============================================================================

// IsOk returns true if the Result contains a successful value.
func (r Result[T]) IsOk() bool {
	return r.isOk
}

// IsError returns true if the Result contains an error.
func (r Result[T]) IsError() bool {
	return !r.isOk
}

// ============================================================================
// Value extraction (UNSAFE - require precondition checks)
// ============================================================================
//
// IMPORTANT: Value() and ErrorInfo() are "unsafe" accessors that panic if
// called in the wrong state. This mirrors Rust's unwrap() behavior.
//
// Design Rationale:
//   - Ada's discriminated records make invalid access impossible at compile time
//   - Go lacks this, so we use runtime panics as the idiomatic alternative
//   - These methods exist for ergonomics after precondition checks
//
// Safe Usage Patterns:
//
//	// Pattern 1: Check first (recommended)
//	if result.IsOk() {
//	    value := result.Value()  // Safe - precondition verified
//	}
//
//	// Pattern 2: Use safe alternatives
//	value := result.UnwrapOr(defaultValue)
//	value := result.UnwrapOrElse(func() T { return computeDefault() })
//
//	// Pattern 3: Document assumption with Expect
//	value := result.Expect("config validated at startup")
//
// When to Use Each:
//   - Value()/ErrorInfo(): After IsOk()/IsError() check, or in tests
//   - UnwrapOr(): When you have a sensible default value
//   - UnwrapOrElse(): When default is expensive to compute
//   - Expect(): When you can document WHY the precondition holds
//
// ============================================================================

// Value returns the success value.
//
// PRECONDITION: Result must be Ok. Caller must verify with IsOk() first.
//
// Panics if the Result is an error. This is intentional - it indicates
// a programmer error (violated precondition), not a runtime failure.
//
// For safe alternatives, see: UnwrapOr, UnwrapOrElse, Expect.
func (r Result[T]) Value() T {
	if !r.isOk {
		panic("called Value() on error Result - precondition violated: must check IsOk() first")
	}
	return r.value
}

// ErrorInfo returns the error value.
//
// PRECONDITION: Result must be Error. Caller must verify with IsError() first.
//
// Panics if the Result is Ok. This is intentional - it indicates
// a programmer error (violated precondition), not a runtime failure.
func (r Result[T]) ErrorInfo() ErrorType {
	if r.isOk {
		panic("called ErrorInfo() on ok Result - precondition violated: must check IsError() first")
	}
	return r.err
}

// ============================================================================
// Unwrap operations (extract value or use default)
// ============================================================================

// UnwrapOr returns the value if Ok, otherwise returns the default value.
//
// Example:
//
//	value := result.UnwrapOr("default")
func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.isOk {
		return r.value
	}
	return defaultValue
}

// UnwrapOrElse returns the value if Ok, otherwise computes default lazily via f.
// Use when default is expensive to compute.
//
// Example:
//
//	value := result.UnwrapOrElse(func() string { return computeDefault() })
func (r Result[T]) UnwrapOrElse(f func() T) T {
	if r.isOk {
		return r.value
	}
	return f()
}

// Expect returns the value if Ok, otherwise panics with the given message.
// Use only when you can document why Result should be Ok at the call site.
//
// Example:
//
//	value := result.Expect("expected valid configuration")
func (r Result[T]) Expect(message string) T {
	if r.isOk {
		return r.value
	}
	panic(message)
}

// ============================================================================
// Functional operations (transform and chain)
// ============================================================================

// Map transforms the success value if Ok, propagates error if Error.
//
// Example:
//
//	doubled := intResult.Map(func(x int) int { return x * 2 })
func (r Result[T]) Map(f func(T) T) Result[T] {
	if r.isOk {
		return Ok(f(r.value))
	}
	return r
}

// MapTo transforms the success value to a different type U if Ok, propagates error if Error.
//
// Example:
//
//	strResult := intResult.MapTo(func(x int) string { return fmt.Sprintf("%d", x) })
func MapTo[T any, U any](r Result[T], f func(T) U) Result[U] {
	if r.isOk {
		return Ok(f(r.value))
	}
	return Err[U](r.err)
}

// AndThen chains fallible operations (monadic bind).
// If Self is Error, propagates error without calling f.
// If Self is Ok, calls f with value (f might return Error).
//
// Example:
//
//	result := parseFile(path).AndThen(validate)
func (r Result[T]) AndThen(f func(T) Result[T]) Result[T] {
	if r.isOk {
		return f(r.value)
	}
	return r
}

// AndThenTo chains fallible operations that return a different type U.
//
// Example:
//
//	userResult := idResult.AndThenTo(func(id int) Result[User] { return findUser(id) })
func AndThenTo[T any, U any](r Result[T], f func(T) Result[U]) Result[U] {
	if r.isOk {
		return f(r.value)
	}
	return Err[U](r.err)
}

// MapError transforms the error value if Error, propagates Ok if Ok.
// Use to add context to errors as they propagate up call stack.
//
// Example:
//
//	result := operation().MapError(func(e ErrorType) ErrorType {
//	    return ErrorType{Kind: e.Kind, Message: "context: " + e.Message}
//	})
func (r Result[T]) MapError(f func(ErrorType) ErrorType) Result[T] {
	if !r.isOk {
		return Err[T](f(r.err))
	}
	return r
}

// ============================================================================
// Fallback and recovery
// ============================================================================

// Fallback tries Primary, if Error then uses Alternative.
// Both are eagerly evaluated.
//
// Example:
//
//	result := primary.Fallback(alternative)
func (r Result[T]) Fallback(alternative Result[T]) Result[T] {
	if r.isOk {
		return r
	}
	return alternative
}

// FallbackWith tries Self, if Error then computes alternative lazily via f.
// Use when alternative is expensive to compute.
//
// Example:
//
//	result := primary.FallbackWith(func() Result[T] { return computeAlternative() })
func (r Result[T]) FallbackWith(f func() Result[T]) Result[T] {
	if r.isOk {
		return r
	}
	return f()
}

// Recover turns error into value via handle function.
// Always returns T (never fails).
//
// Example:
//
//	value := result.Recover(func(e ErrorType) string { return "default" })
func (r Result[T]) Recover(handle func(ErrorType) T) T {
	if r.isOk {
		return r.value
	}
	return handle(r.err)
}

// RecoverWith turns error into another Result via handle function.
// Handle might succeed or return different error.
//
// Example:
//
//	result := original.RecoverWith(func(e ErrorType) Result[T] { return retry() })
func (r Result[T]) RecoverWith(handle func(ErrorType) Result[T]) Result[T] {
	if r.isOk {
		return r
	}
	return handle(r.err)
}

// ============================================================================
// Side effects (for logging/debugging)
// ============================================================================

// Tap executes side effects without changing the Result.
// Returns the same Result for chaining.
//
// Example:
//
//	result := operation().Tap(
//	    func(v T) { log.Info("success", v) },
//	    func(e ErrorType) { log.Error("failed", e) },
//	)
func (r Result[T]) Tap(onOk func(T), onError func(ErrorType)) Result[T] {
	if r.isOk {
		onOk(r.value)
	} else {
		onError(r.err)
	}
	return r
}
