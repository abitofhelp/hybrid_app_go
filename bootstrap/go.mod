// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025 Michael Gardner, A Bit of Help, Inc.

module github.com/abitofhelp/hybrid_app_go/bootstrap

go 1.23

// Bootstrap layer - Composition root
// Depends on ALL layers to wire dependencies together

require (
	github.com/abitofhelp/hybrid_app_go/application v0.0.0
	github.com/abitofhelp/hybrid_app_go/infrastructure v0.0.0
	github.com/abitofhelp/hybrid_app_go/presentation v0.0.0
)

require github.com/abitofhelp/hybrid_app_go/domain v0.0.0 // indirect

replace (
	github.com/abitofhelp/hybrid_app_go/application => ../application
	github.com/abitofhelp/hybrid_app_go/domain => ../domain
	github.com/abitofhelp/hybrid_app_go/infrastructure => ../infrastructure
	github.com/abitofhelp/hybrid_app_go/presentation => ../presentation
)
