// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.

module github.com/abitofhelp/hybrid_app_go/test

go 1.23

// Test module - CAN have external dependencies for integration/e2e tests
// This is separate from /src modules which must have ZERO external dependencies

require (
	github.com/abitofhelp/hybrid_app_go/domain v0.0.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/abitofhelp/hybrid_app_go/domain => ../domain

replace github.com/abitofhelp/hybrid_app_go/application => ../application

replace github.com/abitofhelp/hybrid_app_go/infrastructure => ../infrastructure

replace github.com/abitofhelp/hybrid_app_go/presentation => ../presentation

replace github.com/abitofhelp/hybrid_app_go/bootstrap => ../bootstrap
