// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.

package valueobject_test

import (
	"os"
	"testing"

	"github.com/abitofhelp/hybrid_app_go/domain/test"
)

func TestMain(m *testing.M) {
	test.Reset()
	code := m.Run()

	// Print grand total and final banner
	test.PrintCategorySummary("UNIT TESTS",
		test.GrandTotalTests(),
		test.GrandTotalPassed())

	os.Exit(code)
}
