// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.

//go:build integration

package integration

import (
	"os"
	"testing"

	"github.com/abitofhelp/hybrid_app_go/domain/test"
)

func TestMain(m *testing.M) {
	test.Reset()
	code := m.Run()

	// Print grand total and final banner
	test.PrintCategorySummary("INTEGRATION TESTS",
		test.GrandTotalTests(),
		test.GrandTotalPassed())

	os.Exit(code)
}
