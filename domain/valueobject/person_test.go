// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
// Package: valueobject
// Description: Unit tests for Person value object (stdlib testing only - ZERO external dependencies)

package valueobject

import (
	"strings"
	"testing"

	domerr "github.com/abitofhelp/hybrid_app_go/domain/error"
)

func TestCreatePerson_ValidNames(t *testing.T) {
	tests := map[string]struct {
		name         string
		wantName     string
		wantGreeting string
	}{
		"simple name": {
			name:         "Alice",
			wantName:     "Alice",
			wantGreeting: "Hello, Alice!",
		},
		"name with spaces": {
			name:         "Bob Smith",
			wantName:     "Bob Smith",
			wantGreeting: "Hello, Bob Smith!",
		},
		"name with unicode": {
			name:         "José García",
			wantName:     "José García",
			wantGreeting: "Hello, José García!",
		},
		"single character": {
			name:         "X",
			wantName:     "X",
			wantGreeting: "Hello, X!",
		},
		"max length name": {
			name:         strings.Repeat("a", MaxNameLength),
			wantName:     strings.Repeat("a", MaxNameLength),
			wantGreeting: "Hello, " + strings.Repeat("a", MaxNameLength) + "!",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := CreatePerson(tt.name)

			if !result.IsOk() {
				t.Fatalf("expected Ok result for valid name, got Error: %v", result.ErrorInfo())
			}
			person := result.Value()

			if got := person.GetName(); got != tt.wantName {
				t.Errorf("GetName() = %v, want %v", got, tt.wantName)
			}

			if got := person.GreetingMessage(); got != tt.wantGreeting {
				t.Errorf("GreetingMessage() = %v, want %v", got, tt.wantGreeting)
			}

			if !person.IsValid() {
				t.Error("IsValid() = false, want true")
			}
		})
	}
}

func TestCreatePerson_EmptyName(t *testing.T) {
	result := CreatePerson("")

	if !result.IsError() {
		t.Fatal("expected Err result for empty name")
	}

	domErr := result.ErrorInfo()

	if domErr.Kind != domerr.ValidationError {
		t.Errorf("ErrorKind = %v, want %v", domErr.Kind, domerr.ValidationError)
	}

	if !strings.Contains(domErr.Message, "cannot be empty") {
		t.Errorf("Error message %q does not contain 'cannot be empty'", domErr.Message)
	}
}

func TestCreatePerson_NameTooLong(t *testing.T) {
	longName := strings.Repeat("a", MaxNameLength+1)
	result := CreatePerson(longName)

	if !result.IsError() {
		t.Fatal("expected Err result for name exceeding max length")
	}

	domErr := result.ErrorInfo()

	if domErr.Kind != domerr.ValidationError {
		t.Errorf("ErrorKind = %v, want %v", domErr.Kind, domerr.ValidationError)
	}

	if !strings.Contains(domErr.Message, "exceeds maximum length") {
		t.Errorf("Error message %q does not contain 'exceeds maximum length'", domErr.Message)
	}
}

func TestPerson_GetName(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"simple", "Alice"},
		{"with spaces", "Bob Smith"},
		{"with unicode", "José García"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreatePerson(tt.input)
			if !result.IsOk() {
				t.Fatalf("expected Ok result, got Error: %v", result.ErrorInfo())
			}

			person := result.Value()
			got := person.GetName()

			if got != tt.input {
				t.Errorf("GetName() = %v, want %v", got, tt.input)
			}

			if got == "" {
				t.Error("GetName() returned empty string")
			}

			if len(got) > MaxNameLength {
				t.Errorf("GetName() length = %d, want <= %d", len(got), MaxNameLength)
			}
		})
	}
}

func TestPerson_GreetingMessage(t *testing.T) {
	tests := map[string]struct {
		name            string
		wantContains    []string
		wantStartsWith  string
		wantEndsWith    string
		wantMinLength   int
	}{
		"Alice": {
			name:            "Alice",
			wantContains:    []string{"Hello,", "Alice"},
			wantStartsWith:  "Hello, ",
			wantEndsWith:    "!",
			wantMinLength:   9,
		},
		"Bob Smith": {
			name:            "Bob Smith",
			wantContains:    []string{"Hello,", "Bob Smith"},
			wantStartsWith:  "Hello, ",
			wantEndsWith:    "!",
			wantMinLength:   9,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := CreatePerson(tt.name)
			if !result.IsOk() {
				t.Fatalf("expected Ok result, got Error: %v", result.ErrorInfo())
			}

			person := result.Value()
			greeting := person.GreetingMessage()

			if len(greeting) <= tt.wantMinLength {
				t.Errorf("greeting length = %d, want > %d", len(greeting), tt.wantMinLength)
			}

			if !strings.HasPrefix(greeting, tt.wantStartsWith) {
				t.Errorf("greeting %q does not start with %q", greeting, tt.wantStartsWith)
			}

			if !strings.HasSuffix(greeting, tt.wantEndsWith) {
				t.Errorf("greeting %q does not end with %q", greeting, tt.wantEndsWith)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(greeting, want) {
					t.Errorf("greeting %q does not contain %q", greeting, want)
				}
			}
		})
	}
}

func TestPerson_IsValid(t *testing.T) {
	t.Run("valid person", func(t *testing.T) {
		result := CreatePerson("Alice")
		if !result.IsOk() {
			t.Fatalf("expected Ok result, got Error: %v", result.ErrorInfo())
		}

		person := result.Value()
		if !person.IsValid() {
			t.Error("IsValid() = false, want true")
		}
	})

	t.Run("type invariant maintained", func(t *testing.T) {
		// Test that the type invariant is always maintained
		// for any Person created via CreatePerson
		names := []string{"A", "Alice", "Bob Smith", "José García"}

		for _, name := range names {
			result := CreatePerson(name)
			if result.IsOk() {
				person := result.Value()
				if !person.IsValid() {
					t.Errorf("type invariant violated for name: %s", name)
				}
			}
		}
	})
}
