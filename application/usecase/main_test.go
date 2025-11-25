// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.

package usecase

import (
	"os"
	"testing"

	"github.com/abitofhelp/hybrid_app_go/domain/test"
)

// TestMain is the test runner for the usecase package.
// It aggregates test results and prints a professional summary banner.
func TestMain(m *testing.M) {
	// Reset global counters for fresh run
	test.Reset()

	// Run all tests
	code := m.Run()

	// Print category summary banner
	test.PrintCategorySummary("UNIT TESTS",
		test.GrandTotalTests(),
		test.GrandTotalPassed())

	os.Exit(code)
}
