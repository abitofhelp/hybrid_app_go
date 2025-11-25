module github.com/abitofhelp/hybrid_app_go/test/e2e

go 1.23

require (
	github.com/abitofhelp/hybrid_app_go/domain v0.0.0
)

replace (
	github.com/abitofhelp/hybrid_app_go/domain => ../../domain
)
