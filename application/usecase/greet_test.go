// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.

package usecase

import (
	"bytes"
	"context"
	"testing"

	"github.com/abitofhelp/hybrid_app_go/application/command"
	"github.com/abitofhelp/hybrid_app_go/application/model"
	"github.com/abitofhelp/hybrid_app_go/application/port/outward"
	domerr "github.com/abitofhelp/hybrid_app_go/domain/error"
	"github.com/abitofhelp/hybrid_app_go/domain/test"
)

// TestApplicationUseCaseGreet tests the GreetUseCase application service.
//
// This demonstrates canonical application-layer unit testing:
//   - Stubbing output ports with in-memory implementations
//   - Testing successful orchestration (valid input -> output)
//   - Testing error propagation (invalid input -> no side effects)
//   - Verifying domain validation is respected
func TestApplicationUseCaseGreet(t *testing.T) {
	tf := test.New("Application.UseCase.Greet")

	// ========================================================================
	// Test Fixtures
	// ========================================================================

	// writerCallCount tracks how many times the writer was called
	var writerCallCount int

	// capturedOutput captures what was written
	var capturedOutput bytes.Buffer

	// stubWriter is an in-memory WriterFunc implementation for testing.
	// It captures output and tracks call count without actual I/O.
	stubWriter := func(ctx context.Context, message string) domerr.Result[model.Unit] {
		writerCallCount++
		capturedOutput.WriteString(message)
		capturedOutput.WriteString("\n")
		return domerr.Ok(model.UnitValue)
	}

	// resetFixtures resets the test state between tests
	resetFixtures := func() {
		writerCallCount = 0
		capturedOutput.Reset()
	}

	// ========================================================================
	// Valid Name Tests
	// ========================================================================

	// Test: Valid name produces successful Result
	resetFixtures()
	uc := NewGreetUseCase(stubWriter)
	cmd := command.NewGreetCommand("Alice")
	result := uc.Execute(context.Background(), cmd)

	tf.RunTest("Valid name - IsOk returns true", result.IsOk())
	tf.RunTest("Valid name - writer was called once", writerCallCount == 1)
	tf.RunTest("Valid name - output contains greeting",
		capturedOutput.String() == "Hello, Alice!\n")

	// Test: Name with spaces works correctly
	resetFixtures()
	uc = NewGreetUseCase(stubWriter)
	cmd = command.NewGreetCommand("Bob Smith")
	result = uc.Execute(context.Background(), cmd)

	tf.RunTest("Name with spaces - IsOk returns true", result.IsOk())
	tf.RunTest("Name with spaces - output correct",
		capturedOutput.String() == "Hello, Bob Smith!\n")

	// Test: Single character name works
	resetFixtures()
	uc = NewGreetUseCase(stubWriter)
	cmd = command.NewGreetCommand("X")
	result = uc.Execute(context.Background(), cmd)

	tf.RunTest("Single char name - IsOk returns true", result.IsOk())
	tf.RunTest("Single char name - output correct",
		capturedOutput.String() == "Hello, X!\n")

	// ========================================================================
	// Invalid Name Tests (ValidationError propagation)
	// ========================================================================

	// Test: Empty name bubbles up ValidationError WITHOUT calling writer
	resetFixtures()
	uc = NewGreetUseCase(stubWriter)
	cmd = command.NewGreetCommand("")
	result = uc.Execute(context.Background(), cmd)

	tf.RunTest("Empty name - IsError returns true", result.IsError())
	tf.RunTest("Empty name - error kind is ValidationError",
		result.ErrorInfo().Kind == domerr.ValidationError)
	tf.RunTest("Empty name - writer was NOT called", writerCallCount == 0)
	tf.RunTest("Empty name - no output captured", capturedOutput.Len() == 0)

	// Test: Name too long bubbles up ValidationError WITHOUT calling writer
	resetFixtures()
	uc = NewGreetUseCase(stubWriter)
	longName := string(make([]byte, 101)) // 101 chars exceeds MaxNameLength (100)
	for i := range longName {
		longName = longName[:i] + "a" + longName[i+1:]
	}
	cmd = command.NewGreetCommand(longName)
	result = uc.Execute(context.Background(), cmd)

	tf.RunTest("Long name - IsError returns true", result.IsError())
	tf.RunTest("Long name - error kind is ValidationError",
		result.ErrorInfo().Kind == domerr.ValidationError)
	tf.RunTest("Long name - writer was NOT called", writerCallCount == 0)

	// ========================================================================
	// Writer Error Propagation Tests
	// ========================================================================

	// Test: Writer failure propagates as InfrastructureError
	resetFixtures()
	failingWriter := func(ctx context.Context, message string) domerr.Result[model.Unit] {
		writerCallCount++
		return domerr.Err[model.Unit](domerr.NewInfrastructureError("simulated write failure"))
	}

	uc = NewGreetUseCase(failingWriter)
	cmd = command.NewGreetCommand("Alice")
	result = uc.Execute(context.Background(), cmd)

	tf.RunTest("Writer failure - IsError returns true", result.IsError())
	tf.RunTest("Writer failure - error kind is InfrastructureError",
		result.ErrorInfo().Kind == domerr.InfrastructureError)
	tf.RunTest("Writer failure - writer was called", writerCallCount == 1)

	// ========================================================================
	// Context Cancellation Tests
	// ========================================================================

	// Test: Cancelled context handled by writer (simulated)
	resetFixtures()
	contextAwareWriter := func(ctx context.Context, message string) domerr.Result[model.Unit] {
		writerCallCount++
		select {
		case <-ctx.Done():
			return domerr.Err[model.Unit](domerr.NewInfrastructureError("context cancelled"))
		default:
			capturedOutput.WriteString(message)
			return domerr.Ok(model.UnitValue)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	uc = NewGreetUseCase(contextAwareWriter)
	cmd = command.NewGreetCommand("Alice")
	result = uc.Execute(ctx, cmd)

	tf.RunTest("Cancelled context - IsError returns true", result.IsError())
	tf.RunTest("Cancelled context - error kind is InfrastructureError",
		result.ErrorInfo().Kind == domerr.InfrastructureError)

	// Print summary
	tf.Summary(t)
}

// stubWriterFunc creates a stub WriterFunc that captures output for testing.
// This is exported so other test packages can use it as a reference pattern.
func stubWriterFunc(output *bytes.Buffer, callCount *int) outward.WriterFunc {
	return func(ctx context.Context, message string) domerr.Result[model.Unit] {
		*callCount++
		output.WriteString(message)
		output.WriteString("\n")
		return domerr.Ok(model.UnitValue)
	}
}
