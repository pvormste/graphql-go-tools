module github.com/TykTechnologies/graphql-go-tools

go 1.15

require (
	github.com/OneOfOne/xxhash v1.2.8
	github.com/TykTechnologies/graphql-go-tools/examples/chat v0.0.0-20211001181535-d99d0e668bd2
	github.com/TykTechnologies/graphql-go-tools/examples/federation v0.0.0-20211001181535-d99d0e668bd2
	github.com/buger/jsonparser v1.1.1
	github.com/cespare/xxhash v1.1.0
	github.com/dave/jennifer v1.4.0
	github.com/davecgh/go-spew v1.1.1
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/evanphx/json-patch/v5 v5.1.0
	github.com/go-test/deep v1.0.4
	github.com/gobuffalo/packr v1.30.1
	github.com/gobwas/ws v1.0.4
	github.com/golang/mock v1.4.1
	github.com/google/uuid v1.1.1
	github.com/hashicorp/golang-lru v0.5.4
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334
	github.com/jensneuse/abstractlogger v0.0.4
	github.com/jensneuse/byte-template v0.0.0-20200214152254-4f3cf06e5c68
	github.com/jensneuse/diffview v1.0.0
	github.com/jensneuse/pipeline v0.0.0-20200117120358-9fb4de085cd6
	github.com/mitchellh/go-homedir v1.1.0
	github.com/nats-io/nats-server/v2 v2.6.1 // indirect
	github.com/nats-io/nats.go v1.12.3
	github.com/sebdah/goldie v0.0.0-20180424091453-8784dd1ab561
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/gjson v1.8.1
	github.com/tidwall/sjson v1.0.4
	go.uber.org/atomic v1.9.0
	go.uber.org/zap v1.18.1
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	nhooyr.io/websocket v1.8.7
)

// replace github.com/TykTechnologies/graphql-go-tools/examples/federation => ./examples/federation

// replace github.com/TykTechnologies/graphql-go-tools/examples/chat => ./examples/chat
