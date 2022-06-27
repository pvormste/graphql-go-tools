module github.com/pvormste/graphql-go-tools/examples/federation

go 1.16

require (
	github.com/99designs/gqlgen v0.17.10
	github.com/gobwas/ws v1.0.4
	github.com/gorilla/websocket v1.5.0
	github.com/jensneuse/abstractlogger v0.0.4
	github.com/nats-io/nats-server/v2 v2.3.2 // indirect
	github.com/pvormste/graphql-go-tools v1.6.2-0.20220627101903-ba97d05b8fdb
	github.com/vektah/gqlparser/v2 v2.4.5
	go.uber.org/atomic v1.9.0
	go.uber.org/zap v1.18.1
)

replace github.com/pvormste/graphql-go-tools => ../../
