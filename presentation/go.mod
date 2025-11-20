// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.

module github.com/abitofhelp/hybrid_app_go/presentation

go 1.23

// Presentation layer - Driving/Primary adapters (CLI)
// Depends ONLY on application layer (NOT domain directly)
// This enforces the architectural rule that Presentation cannot access Domain

require github.com/abitofhelp/hybrid_app_go/application v0.0.0

replace github.com/abitofhelp/hybrid_app_go/application => ../application
