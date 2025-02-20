package graph

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/debug"

	"github.com/pvormste/graphql-go-tools/pkg/graphql/federationtesting/accounts/graph/generated"
)

type EndpointOptions struct {
	EnableDebug bool
}

var TestOptions = EndpointOptions{
	EnableDebug: false,
}

func GraphQLEndpointHandler(opts EndpointOptions) http.Handler {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{}}))
	if opts.EnableDebug {
		srv.Use(&debug.Tracer{})
	}

	return srv
}
