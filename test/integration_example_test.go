// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.

// Package test contains integration and e2e tests
// This module CAN use external test frameworks like testify
// because it's separate from the /src modules
package test

import (
	"testing"

	"github.com/abitofhelp/hybrid_app_go/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntegration_PersonCreation demonstrates integration testing with testify
// Note: testify IS allowed here because we're in the /test module
func TestIntegration_PersonCreation(t *testing.T) {
	// Arrange
	name := "Alice"

	// Act
	result := valueobject.CreatePerson(name)

	// Assert - Using testify is allowed in /test module
	require.True(t, result.IsOk(), "Expected valid person creation")

	person := result.Value()
	assert.Equal(t, name, person.GetName())
	assert.True(t, person.IsValid())

	greeting := person.GreetingMessage()
	assert.Contains(t, greeting, name)
	assert.Contains(t, greeting, "Hello,")
}
