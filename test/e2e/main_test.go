// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.

//go:build e2e

package e2e

import (
	"os"
	"testing"

	"github.com/abitofhelp/hybrid_app_go/domain/test"
)

func TestMain(m *testing.M) {
	test.Reset()
	code := m.Run()

	// Print grand total and final banner
	test.PrintCategorySummary("E2E TESTS",
		test.GrandTotalTests(),
		test.GrandTotalPassed())

	os.Exit(code)
}
