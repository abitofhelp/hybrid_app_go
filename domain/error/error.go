// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
// Package: error
// Description: Domain error types and error handling primitives

// Package error provides domain error types used throughout the application for
// consistent error reporting. This package defines the core error primitives
// that enable railway-oriented programming via Result monads.
//
// Architecture Notes:
//   - Part of the DOMAIN layer (innermost, zero external dependencies except mo)
//   - Error types are concrete (not generic) for consistency
//   - Used with mo.Result[T] monad for functional error handling
//
// Usage:
//
//	import "github.com/abitofhelp/hybrid_app_go/domain/error"
//
//	result := someOperation()
//	if result.IsError() {
//	    err := result.Error()
//	    switch err.Kind {
//	    case error.ValidationError:
//	        // Handle validation error
//	    case error.InfrastructureError:
//	        // Handle infrastructure error
//	    }
//	}
package error

import "fmt"

// ErrorKind represents categories of errors that can occur in the application.
// This enables pattern matching and different handling strategies per category.
type ErrorKind int

const (
	// ValidationError indicates domain validation failures (invalid input)
	ValidationError ErrorKind = iota

	// InfrastructureError indicates infrastructure failures (I/O, network, DB)
	InfrastructureError
)

// String returns a human-readable representation of the ErrorKind.
func (k ErrorKind) String() string {
	switch k {
	case ValidationError:
		return "ValidationError"
	case InfrastructureError:
		return "InfrastructureError"
	default:
		return "UnknownError"
	}
}

// ErrorType is the concrete error type used throughout the application.
// It combines an error category (Kind) with a descriptive message.
//
// Contract:
//   - Message should be non-empty when creating errors
//   - Kind should be a valid ErrorKind value
type ErrorType struct {
	Kind    ErrorKind
	Message string
}

// Error implements the error interface for ErrorType.
// This allows ErrorType to be used as a standard Go error when needed.
func (e ErrorType) Error() string {
	return fmt.Sprintf("%s: %s", e.Kind, e.Message)
}

// NewValidationError creates a new validation error with the given message.
func NewValidationError(message string) ErrorType {
	return ErrorType{
		Kind:    ValidationError,
		Message: message,
	}
}

// NewInfrastructureError creates a new infrastructure error with the given message.
func NewInfrastructureError(message string) ErrorType {
	return ErrorType{
		Kind:    InfrastructureError,
		Message: message,
	}
}
