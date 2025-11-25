// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.

// Package valueobject_test provides unit tests for Person value object
// using the Ada-style test framework for consistent cross-language reporting.
package valueobject_test

import (
	"strings"
	"testing"

	domerr "github.com/abitofhelp/hybrid_app_go/domain/error"
	"github.com/abitofhelp/hybrid_app_go/domain/test"
	"github.com/abitofhelp/hybrid_app_go/domain/valueobject"
)

// TestDomainValueObjectPerson tests the Person value object.
// Uses Ada-style [PASS]/[FAIL] output for uniform cross-language reporting.
func TestDomainValueObjectPerson(t *testing.T) {
	tf := test.New("Domain.ValueObject.Person")

	// ========================================================================
	// Test: CreatePerson with valid name
	// ========================================================================

	r1 := valueobject.CreatePerson("Alice")
	tf.RunTest("CreatePerson valid - IsOk returns true", r1.IsOk())
	if r1.IsOk() {
		person := r1.Value()
		tf.RunTest("CreatePerson valid - GetName returns correct name",
			person.GetName() == "Alice")
		tf.RunTest("CreatePerson valid - IsValid returns true",
			person.IsValid())
	}

	// ========================================================================
	// Test: CreatePerson with empty name (validation error)
	// ========================================================================

	r2 := valueobject.CreatePerson("")
	tf.RunTest("CreatePerson empty - IsError returns true", r2.IsError())
	if r2.IsError() {
		info := r2.ErrorInfo()
		tf.RunTest("CreatePerson empty - error kind is ValidationError",
			info.Kind == domerr.ValidationError)
		tf.RunTest("CreatePerson empty - error message mentions 'empty'",
			strings.Contains(info.Message, "empty"))
	}

	// ========================================================================
	// Test: CreatePerson with name too long (validation error)
	// ========================================================================

	longName := strings.Repeat("a", valueobject.MaxNameLength+1)
	r3 := valueobject.CreatePerson(longName)
	tf.RunTest("CreatePerson too long - IsError returns true", r3.IsError())
	if r3.IsError() {
		info := r3.ErrorInfo()
		tf.RunTest("CreatePerson too long - error kind is ValidationError",
			info.Kind == domerr.ValidationError)
		tf.RunTest("CreatePerson too long - error message mentions 'exceeds'",
			strings.Contains(info.Message, "exceeds"))
	}

	// ========================================================================
	// Test: GreetingMessage format
	// ========================================================================

	r4 := valueobject.CreatePerson("Bob")
	if r4.IsOk() {
		person := r4.Value()
		greeting := person.GreetingMessage()
		tf.RunTest("GreetingMessage - starts with 'Hello, '",
			strings.HasPrefix(greeting, "Hello, "))
		tf.RunTest("GreetingMessage - contains name",
			strings.Contains(greeting, "Bob"))
		tf.RunTest("GreetingMessage - ends with '!'",
			strings.HasSuffix(greeting, "!"))
		tf.RunTest("GreetingMessage - exact format",
			greeting == "Hello, Bob!")
	}

	// ========================================================================
	// Test: Name with spaces
	// ========================================================================

	r5 := valueobject.CreatePerson("Bob Smith")
	tf.RunTest("Name with spaces - IsOk returns true", r5.IsOk())
	if r5.IsOk() {
		person := r5.Value()
		tf.RunTest("Name with spaces - preserves spaces",
			person.GetName() == "Bob Smith")
		tf.RunTest("Name with spaces - greeting correct",
			person.GreetingMessage() == "Hello, Bob Smith!")
	}

	// ========================================================================
	// Test: Name with unicode characters
	// ========================================================================

	r6 := valueobject.CreatePerson("José García")
	tf.RunTest("Unicode name - IsOk returns true", r6.IsOk())
	if r6.IsOk() {
		person := r6.Value()
		tf.RunTest("Unicode name - preserves unicode",
			person.GetName() == "José García")
	}

	// ========================================================================
	// Test: Single character name (boundary)
	// ========================================================================

	r7 := valueobject.CreatePerson("X")
	tf.RunTest("Single char name - IsOk returns true", r7.IsOk())
	if r7.IsOk() {
		person := r7.Value()
		tf.RunTest("Single char name - GetName correct",
			person.GetName() == "X")
	}

	// ========================================================================
	// Test: Max length name (boundary)
	// ========================================================================

	maxName := strings.Repeat("a", valueobject.MaxNameLength)
	r8 := valueobject.CreatePerson(maxName)
	tf.RunTest("Max length name - IsOk returns true", r8.IsOk())
	if r8.IsOk() {
		person := r8.Value()
		tf.RunTest("Max length name - GetName correct",
			person.GetName() == maxName)
		tf.RunTest("Max length name - length is MaxNameLength",
			len(person.GetName()) == valueobject.MaxNameLength)
	}

	// Print summary and fail test if any failed
	tf.Summary(t)
}
