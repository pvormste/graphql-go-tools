//go:generate go run -mod=mod github.com/99designs/gqlgen
package products

import (
	"net/http"

	"github.com/pvormste/graphql-go-tools/pkg/graphql/federationtesting/products/graph"
)

func Handler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", graph.GraphQLEndpointHandler(graph.EndpointOptions{EnableDebug: true}))
	mux.HandleFunc("/websocket_connections", graph.WebsocketConnectionsHandler)

	return mux
}
