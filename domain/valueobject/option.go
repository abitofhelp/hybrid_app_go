// ===========================================================================
// option.go
// ===========================================================================
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
// SPDX-License-Identifier: BSD-3-Clause
//
// Purpose:
//   Option[T] monad for handling nullable/optional values without nil pointers.
//   Provides type-safe handling of values that may or may not be present.
//
// Architecture Notes:
//   - Generic over value type T
//   - Pure domain implementation (ZERO external module dependencies)
//   - Enables functional composition of optional values
//   - Alternative to nil pointers
//
// Design Pattern:
//   - Some(value): Contains a value
//   - None: No value present
//   - Forces explicit handling of missing values at compile time
// ===========================================================================

// Package valueobject provides domain value object types.
package valueobject

// Option represents a value that may or may not be present.
// This is the core functional type for handling optional values.
//
// States:
//   - Some: Contains a value of type T
//   - None: No value present
//
// Usage:
//
//	opt := Some[string]("value")
//	if opt.IsSome() {
//	    value := opt.Value()
//	}
type Option[T any] struct {
	value  T
	isSome bool
}

// ============================================================================
// Constructors
// ============================================================================

// Some creates an Option containing a value.
//
// Example:
//
//	opt := Some[int](42)
func Some[T any](value T) Option[T] {
	return Option[T]{
		value:  value,
		isSome: true,
	}
}

// None creates an empty Option (no value).
//
// Example:
//
//	opt := None[int]()
func None[T any]() Option[T] {
	return Option[T]{
		isSome: false,
	}
}

// ============================================================================
// Query functions
// ============================================================================

// IsSome returns true if Option contains a value.
func (o Option[T]) IsSome() bool {
	return o.isSome
}

// IsNone returns true if Option is empty.
func (o Option[T]) IsNone() bool {
	return !o.isSome
}

// ============================================================================
// Value extraction
// ============================================================================

// Value returns the contained value.
// Panics if Option is None. Check IsSome() first.
func (o Option[T]) Value() T {
	if !o.isSome {
		panic("called Value() on None Option")
	}
	return o.value
}

// ============================================================================
// Unwrap operations (extract value or use default)
// ============================================================================

// UnwrapOr returns the value if Some, otherwise returns the default value.
//
// Example:
//
//	value := opt.UnwrapOr("default")
func (o Option[T]) UnwrapOr(defaultValue T) T {
	if o.isSome {
		return o.value
	}
	return defaultValue
}

// UnwrapOrElse returns the value if Some, otherwise computes default lazily via f.
// Use when default is expensive to compute.
//
// Example:
//
//	value := opt.UnwrapOrElse(func() string { return computeDefault() })
func (o Option[T]) UnwrapOrElse(f func() T) T {
	if o.isSome {
		return o.value
	}
	return f()
}

// ============================================================================
// Functional operations (transform and chain)
// ============================================================================

// Map transforms the contained value if Some, propagates None if None.
//
// Example:
//
//	doubled := intOpt.Map(func(x int) int { return x * 2 })
func (o Option[T]) Map(f func(T) T) Option[T] {
	if o.isSome {
		return Some(f(o.value))
	}
	return o
}

// MapTo transforms the contained value to a different type U if Some, propagates None if None.
//
// Example:
//
//	strOpt := intOpt.MapTo(func(x int) string { return fmt.Sprintf("%d", x) })
func MapTo[T any, U any](o Option[T], f func(T) U) Option[U] {
	if o.isSome {
		return Some(f(o.value))
	}
	return None[U]()
}

// AndThen chains optional operations (monadic bind).
// If Self is None, propagates None without calling f.
// If Self is Some, calls f with value (f might return None).
//
// Example:
//
//	emailOpt := getUserOpt(id).AndThen(func(user User) Option[string] { return user.Email })
func (o Option[T]) AndThen(f func(T) Option[T]) Option[T] {
	if o.isSome {
		return f(o.value)
	}
	return o
}

// AndThenTo chains optional operations that return a different type U.
//
// Example:
//
//	userOpt := idOpt.AndThenTo(func(id int) Option[User] { return findUser(id) })
func AndThenTo[T any, U any](o Option[T], f func(T) Option[U]) Option[U] {
	if o.isSome {
		return f(o.value)
	}
	return None[U]()
}

// Filter keeps value only if predicate holds, otherwise None.
//
// Example:
//
//	evenOpt := intOpt.Filter(func(x int) bool { return x%2 == 0 })
func (o Option[T]) Filter(pred func(T) bool) Option[T] {
	if o.isSome && pred(o.value) {
		return o
	}
	return None[T]()
}

// ============================================================================
// Fallback
// ============================================================================

// OrElse tries Primary, if None then uses Alternative.
// Both are eagerly evaluated.
//
// Example:
//
//	opt := primary.OrElse(alternative)
func (o Option[T]) OrElse(alternative Option[T]) Option[T] {
	if o.isSome {
		return o
	}
	return alternative
}

// OrElseWith tries Self, if None then computes alternative lazily via f.
// Use when alternative is expensive to compute.
//
// Example:
//
//	opt := primary.OrElseWith(func() Option[T] { return computeAlternative() })
func (o Option[T]) OrElseWith(f func() Option[T]) Option[T] {
	if o.isSome {
		return o
	}
	return f()
}
