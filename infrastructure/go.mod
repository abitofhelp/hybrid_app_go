// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.

module github.com/abitofhelp/hybrid_app_go/infrastructure

go 1.23

// Infrastructure layer - Driven/Secondary adapters
// Depends on application + domain layers

require github.com/abitofhelp/hybrid_app_go/application v0.0.0

require github.com/abitofhelp/hybrid_app_go/domain v0.0.0

replace (
	github.com/abitofhelp/hybrid_app_go/application => ../application
	github.com/abitofhelp/hybrid_app_go/domain => ../domain
)
