package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/jensneuse/graphql-go-tools/examples/federation/accounts/graph/generated"
	"github.com/jensneuse/graphql-go-tools/examples/federation/accounts/graph/model"
)

func (r *mutationResolver) Login(ctx context.Context, username string, password string, userID *string) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	user := users["1"]
	return &user, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	user, ok := users[id]
	if !ok {
		return nil, nil
	}
	return &user, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
