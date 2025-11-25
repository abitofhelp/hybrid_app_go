module github.com/abitofhelp/hybrid_app_go/test/integration

go 1.23

require (
	github.com/abitofhelp/hybrid_app_go/application v0.0.0
	github.com/abitofhelp/hybrid_app_go/domain v0.0.0
	github.com/abitofhelp/hybrid_app_go/infrastructure v0.0.0
)

replace (
	github.com/abitofhelp/hybrid_app_go/application => ../../application
	github.com/abitofhelp/hybrid_app_go/domain => ../../domain
	github.com/abitofhelp/hybrid_app_go/infrastructure => ../../infrastructure
)
