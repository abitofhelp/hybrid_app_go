// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.

// Package test provides a reusable test framework for professional test suites.
//
// This is a Go port of the Ada Test_Framework package, providing:
//   - Test result tracking across multiple test modules
//   - Standardized [PASS]/[FAIL] output formatting
//   - Professional color-coded category summary banners
//
// The framework lives in domain/test to ensure consistent test appearance
// across all languages in our hybrid architecture projects.
//
// Usage Pattern (similar to Ada):
//
//	func TestDomainErrorResult(t *testing.T) {
//	    tf := test.New("Domain.Error.Result")
//
//	    // Run individual tests
//	    tf.RunTest("Ok construction - Is_Ok returns true", result.IsOk())
//	    tf.RunTest("Ok value extraction", value == 42)
//
//	    // Print summary and fail if any tests failed
//	    tf.Summary(t)
//	}
//
// For test runners that aggregate multiple test modules:
//
//	func TestMain(m *testing.M) {
//	    test.Reset()
//	    code := m.Run()
//	    test.PrintCategorySummary("UNIT TESTS",
//	        test.GrandTotalTests(),
//	        test.GrandTotalPassed())
//	    os.Exit(code)
//	}
package test

import (
	"fmt"
	"sync"
	"testing"
)

// ANSI color codes for professional output.
const (
	ColorGreen = "\033[1;92m" // Bright green
	ColorRed   = "\033[1;91m" // Bright red
	ColorReset = "\033[0m"    // Reset to default
)

// Global test counters (thread-safe for parallel tests).
var (
	mu          sync.Mutex
	totalTests  int
	totalPassed int
)

// Framework tracks test results for a single test module.
type Framework struct {
	name   string
	total  int
	passed int
}

// New creates a new test framework instance for a test module.
func New(moduleName string) *Framework {
	fmt.Println("========================================")
	fmt.Printf("Testing: %s\n", moduleName)
	fmt.Println("========================================")
	fmt.Println()

	return &Framework{
		name:   moduleName,
		total:  0,
		passed: 0,
	}
}

// RunTest executes a single test and records the result.
// Prints [PASS] (green) or [FAIL] (red) with the test name.
func (f *Framework) RunTest(name string, passed bool) {
	f.total++
	if passed {
		f.passed++
		fmt.Printf("%s[PASS]%s %s\n", ColorGreen, ColorReset, name)
	} else {
		fmt.Printf("%s[FAIL]%s %s\n", ColorRed, ColorReset, name)
	}
}

// RunTestWithError executes a test that may return an error.
// The test passes if err is nil, fails otherwise.
func (f *Framework) RunTestWithError(name string, err error) {
	f.total++
	if err == nil {
		f.passed++
		fmt.Printf("%s[PASS]%s %s\n", ColorGreen, ColorReset, name)
	} else {
		fmt.Printf("%s[FAIL]%s %s: %v\n", ColorRed, ColorReset, name, err)
	}
}

// Total returns the total number of tests run in this module.
func (f *Framework) Total() int {
	return f.total
}

// Passed returns the number of tests that passed in this module.
func (f *Framework) Passed() int {
	return f.passed
}

// Failed returns the number of tests that failed in this module.
func (f *Framework) Failed() int {
	return f.total - f.passed
}

// Summary prints the test summary for this module and registers results.
// If any tests failed, it calls t.Fail() to mark the test as failed.
func (f *Framework) Summary(t *testing.T) {
	fmt.Println()
	fmt.Println("========================================")
	fmt.Printf("Test Summary: %s\n", f.name)
	fmt.Println("========================================")
	fmt.Printf("Total tests: %d\n", f.total)
	fmt.Printf("Passed:      %d\n", f.passed)
	fmt.Printf("Failed:      %d\n", f.total-f.passed)
	fmt.Println()

	// Register results with global counters
	RegisterResults(f.total, f.passed)

	// Fail the Go test if any tests failed
	if f.passed != f.total {
		t.Fail()
	}
}

// SummaryNoFail prints the test summary without failing the Go test.
// Use this for informational output when you want to aggregate results.
func (f *Framework) SummaryNoFail() {
	fmt.Println()
	fmt.Println("========================================")
	fmt.Printf("Test Summary: %s\n", f.name)
	fmt.Println("========================================")
	fmt.Printf("Total tests: %d\n", f.total)
	fmt.Printf("Passed:      %d\n", f.passed)
	fmt.Printf("Failed:      %d\n", f.total-f.passed)
	fmt.Println()

	// Register results with global counters
	RegisterResults(f.total, f.passed)
}

// RegisterResults adds test results to the global counters.
// Thread-safe for parallel test execution.
func RegisterResults(total, passed int) {
	mu.Lock()
	defer mu.Unlock()
	totalTests += total
	totalPassed += passed
}

// GrandTotalTests returns the cumulative total tests across all modules.
func GrandTotalTests() int {
	mu.Lock()
	defer mu.Unlock()
	return totalTests
}

// GrandTotalPassed returns the cumulative passed tests across all modules.
func GrandTotalPassed() int {
	mu.Lock()
	defer mu.Unlock()
	return totalPassed
}

// Reset clears the global test counters.
// Call this at the start of a test runner.
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	totalTests = 0
	totalPassed = 0
}

// PrintCategorySummary prints a professional color-coded summary banner.
// Returns 0 for success (all tests passed), 1 for failure (any tests failed).
//
// Success output (bright green):
//
//	########################################
//	###                                  ###
//	###    UNIT TESTS: SUCCESS
//	###    All  42 tests passed!
//	###                                  ###
//	########################################
//
// Failure output (bright red):
//
//	########################################
//	###                                  ###
//	###    UNIT TESTS: FAILURE
//	###    3 of 42 tests failed
//	###                                  ###
//	########################################
func PrintCategorySummary(categoryName string, total, passed int) int {
	fmt.Println()

	if passed == total {
		// Success: Bright green box
		fmt.Println(ColorGreen + "########################################")
		fmt.Println("###                                  ###")
		fmt.Printf("###    %s: SUCCESS\n", categoryName)
		fmt.Printf("###    All %d tests passed!\n", total)
		fmt.Println("###                                  ###")
		fmt.Println("########################################" + ColorReset)
		fmt.Println()
		return 0 // Success exit code
	}

	// Failure: Bright red box
	fmt.Println(ColorRed + "########################################")
	fmt.Println("###                                  ###")
	fmt.Printf("###    %s: FAILURE\n", categoryName)
	fmt.Printf("###    %d of %d tests failed\n", total-passed, total)
	fmt.Println("###                                  ###")
	fmt.Println("########################################" + ColorReset)
	fmt.Println()
	return 1 // Failure exit code
}

// AllPassed returns true if all registered tests passed.
func AllPassed() bool {
	mu.Lock()
	defer mu.Unlock()
	return totalTests > 0 && totalPassed == totalTests
}
