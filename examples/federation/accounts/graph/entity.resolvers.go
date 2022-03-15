package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/jensneuse/graphql-go-tools/examples/federation/accounts/graph/generated"
	"github.com/jensneuse/graphql-go-tools/examples/federation/accounts/graph/model"
)

func (r *entityResolver) FindPasswordAccountByEmail(ctx context.Context, email string) (*model.PasswordAccount, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *entityResolver) FindSMSAccountByNumber(ctx context.Context, number *string) (*model.SMSAccount, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *entityResolver) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	return &model.User{
		ID: id,
	}, nil
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
