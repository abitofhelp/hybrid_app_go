// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
// Package: valueobject
// Description: Person value object for the greeter domain

// Package valueobject provides domain value objects - immutable objects
// defined by their attributes rather than identity.
//
// Architecture Notes:
//   - Part of the DOMAIN layer (innermost, pure business logic)
//   - Value objects are immutable after creation
//   - Smart constructors enforce validation
//   - Returns Result[T] for validation (no panics)
//   - Pure domain logic - ZERO external module dependencies
//
// Usage:
//
//	import "github.com/abitofhelp/hybrid_app_go/domain/valueobject"
//
//	result := valueobject.CreatePerson("Alice")
//	if result.IsOk() {
//	    person := result.Value()
//	    message := person.GreetingMessage()
//	}
package valueobject

import (
	"fmt"

	domerr "github.com/abitofhelp/hybrid_app_go/domain/error"
)

const (
	// MaxNameLength is the maximum allowed length for a person's name.
	// This is a reasonable limit for person names in most applications.
	MaxNameLength = 100
)

// Person represents a person's name as an immutable value object.
//
// Design Pattern: Value Object
//   - Immutable after creation
//   - Validation enforced via Create smart constructor
//   - Defined by attributes (name) not identity
//   - No setters - create new instance for changes
//
// Contract:
//   - Name is never empty (enforced by Create)
//   - Name never exceeds MaxNameLength (enforced by Create)
//   - Use Create() to instantiate, not struct literal
type Person struct {
	name string
}

// CreatePerson creates a new Person value object with validation.
//
// This is the RECOMMENDED way to create a Person. Direct struct instantiation
// bypasses validation and should be avoided.
//
// Validation rules:
//  1. Name must not be empty
//  2. Name must not exceed MaxNameLength
//  3. Name may contain spaces and Unicode characters
//  4. Whitespace is preserved exactly as provided
//
// Returns:
//   - domerr.Result[Person] - Ok if valid, Err if validation fails
//
// Contract (expressed in comments since Go lacks Pre/Post):
//   - Pre: name parameter can be any string
//   - Post: If name is empty or exceeds MaxNameLength, returns Err
//   - Post: If valid, returns Ok with Person where GetName() returns exact input
func CreatePerson(name string) domerr.Result[Person] {
	// Validation 1: Check for empty string
	if len(name) == 0 {
		return domerr.Err[Person](domerr.NewValidationError("Person name cannot be empty"))
	}

	// Validation 2: Check maximum length
	if len(name) > MaxNameLength {
		return domerr.Err[Person](domerr.NewValidationError(
			fmt.Sprintf("Person name exceeds maximum length of %d characters", MaxNameLength)))
	}

	// All validations passed - create the value object
	return domerr.Ok(Person{name: name})
}

// GetName returns the string representation of the person's name.
//
// Contract:
//   - Post: Result is never empty (enforced by Create validation)
//   - Post: Result length <= MaxNameLength (enforced by Create validation)
func (p Person) GetName() string {
	return p.name
}

// GreetingMessage generates a greeting message for this person.
//
// Pure domain logic - no side effects.
//
// Contract:
//   - Post: Result always starts with "Hello, " and ends with "!"
//   - Post: Result length is always > 9 (len("Hello, !") == 8)
func (p Person) GreetingMessage() string {
	return fmt.Sprintf("Hello, %s!", p.name)
}

// IsValid checks if the person satisfies the type invariant.
//
// Type Invariant: A Person is valid if and only if its name is non-empty.
// This invariant must always hold for any Person instance.
//
// This method is primarily used for testing and debugging to verify that
// the invariant is maintained.
func (p Person) IsValid() bool {
	return len(p.name) > 0
}
