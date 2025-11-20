// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.
// Package: main
// Description: Main entry point for greeter CLI application

// Package main is the entry point for the greeter application.
// This is intentionally minimal - all logic lives in the Bootstrap layer.
//
// Architecture Notes:
//   - Minimal entry point (3-line implementation)
//   - Delegates to Bootstrap.Run for all logic
//   - Only responsible for process exit code
//   - No business logic here
//
// Usage:
//
//	./greeter Alice
//	Output: Hello, Alice!
package main

import (
	"os"

	"github.com/abitofhelp/hybrid_app_go/bootstrap/cli"
)

func main() {
	// Delegate to Bootstrap layer for all logic
	exitCode := cli.Run(os.Args)

	// Set process exit code
	os.Exit(exitCode)
}
