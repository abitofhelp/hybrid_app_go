// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.

package error_test

import (
	"testing"

	domerr "github.com/abitofhelp/hybrid_app_go/domain/error"
	"github.com/abitofhelp/hybrid_app_go/domain/test"
)

// TestDomainErrorResult tests the domain Result[T] monad functionality.
// This is a Go port of test_domain_error_result.adb from the Ada project.
func TestDomainErrorResult(t *testing.T) {
	tf := test.New("Domain.Error.Result")

	// ========================================================================
	// Test: Ok construction and Is_Ok query
	// ========================================================================

	r1 := domerr.Ok(42)
	tf.RunTest("Ok construction - IsOk returns true", r1.IsOk())
	tf.RunTest("Ok construction - IsError returns false", !r1.IsError())

	// ========================================================================
	// Test: Ok value extraction
	// ========================================================================

	r2 := domerr.Ok(123)
	if r2.IsOk() {
		val := r2.Value()
		tf.RunTest("Ok value extraction - correct value", val == 123)
	} else {
		tf.RunTest("Ok value extraction - Result should be Ok", false)
	}

	// ========================================================================
	// Test: Error construction and Is_Error query
	// ========================================================================

	r3 := domerr.Err[int](domerr.ErrorType{
		Kind:    domerr.ValidationError,
		Message: "Test validation error",
	})
	tf.RunTest("Error construction - IsError returns true", r3.IsError())
	tf.RunTest("Error construction - IsOk returns false", !r3.IsOk())

	// ========================================================================
	// Test: Error info extraction
	// ========================================================================

	r4 := domerr.Err[int](domerr.ErrorType{
		Kind:    domerr.InfrastructureError,
		Message: "Test infra error",
	})
	if r4.IsError() {
		info := r4.ErrorInfo()
		tf.RunTest("Error info - correct kind", info.Kind == domerr.InfrastructureError)
		tf.RunTest("Error info - correct message", info.Message == "Test infra error")
	} else {
		tf.RunTest("Error info extraction - Result should be Error", false)
	}

	// ========================================================================
	// Test: Result with boolean type
	// ========================================================================

	r5 := domerr.Ok(true)
	tf.RunTest("Boolean Result - IsOk returns true", r5.IsOk())
	if r5.IsOk() {
		tf.RunTest("Boolean Result - correct value", r5.Value() == true)
	}

	// ========================================================================
	// Test: Error with empty message
	// ========================================================================

	r6 := domerr.Err[int](domerr.ErrorType{
		Kind:    domerr.ValidationError,
		Message: "",
	})
	tf.RunTest("Error with empty message - IsError", r6.IsError())
	if r6.IsError() {
		info := r6.ErrorInfo()
		tf.RunTest("Error with empty message - message is empty", info.Message == "")
	}

	// ========================================================================
	// Test: Multiple Ok values don't interfere
	// ========================================================================

	r7 := domerr.Ok(100)
	r8 := domerr.Ok(200)
	tf.RunTest("Multiple Ok values - R1 has correct value",
		r7.IsOk() && r7.Value() == 100)
	tf.RunTest("Multiple Ok values - R2 has correct value",
		r8.IsOk() && r8.Value() == 200)

	// ========================================================================
	// Test: Multiple Error values don't interfere
	// ========================================================================

	r9 := domerr.Err[int](domerr.ErrorType{
		Kind:    domerr.ValidationError,
		Message: "Error 1",
	})
	r10 := domerr.Err[int](domerr.ErrorType{
		Kind:    domerr.InfrastructureError,
		Message: "Error 2",
	})
	if r9.IsError() && r10.IsError() {
		info1 := r9.ErrorInfo()
		info2 := r10.ErrorInfo()
		tf.RunTest("Multiple errors - R1 has correct kind",
			info1.Kind == domerr.ValidationError)
		tf.RunTest("Multiple errors - R1 has correct message",
			info1.Message == "Error 1")
		tf.RunTest("Multiple errors - R2 has correct kind",
			info2.Kind == domerr.InfrastructureError)
		tf.RunTest("Multiple errors - R2 has correct message",
			info2.Message == "Error 2")
	} else {
		tf.RunTest("Multiple errors test failed", false)
	}

	// ========================================================================
	// Test: UnwrapOr with Ok value
	// ========================================================================

	r11 := domerr.Ok(42)
	tf.RunTest("UnwrapOr with Ok - returns value", r11.UnwrapOr(0) == 42)

	// ========================================================================
	// Test: UnwrapOr with Error value
	// ========================================================================

	r12 := domerr.Err[int](domerr.ErrorType{
		Kind:    domerr.ValidationError,
		Message: "error",
	})
	tf.RunTest("UnwrapOr with Error - returns default", r12.UnwrapOr(99) == 99)

	// Print summary and fail test if any failed
	tf.Summary(t)
}
