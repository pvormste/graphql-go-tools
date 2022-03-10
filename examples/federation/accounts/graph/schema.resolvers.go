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

func (r *userResolver) Metadata(ctx context.Context, obj *model.User) ([]*model.UserMetadata, error) {
	if obj == nil {
		return nil, nil
	}

	for _, metadataEntity := range userMetadataDB {
		if metadataEntity.ID != obj.ID {
			continue
		}

		userMetadataResult := make([]*model.UserMetadata, len(metadataEntity.Metadata))
		for i, m := range metadataEntity.Metadata {
			userMetadataResult[i] = &model.UserMetadata{
				Name:        m.Name,
				Address:     m.Address,
				Description: m.Description,
			}
		}

		return userMetadataResult, nil
	}

	return nil, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
