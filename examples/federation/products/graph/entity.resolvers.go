package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/jensneuse/graphql-go-tools/examples/federation/products/graph/generated"
	"github.com/jensneuse/graphql-go-tools/examples/federation/products/graph/model"
)

func (r *entityResolver) FindCarByID(ctx context.Context, id string) (*model.Car, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *entityResolver) FindFurnitureByUpc(ctx context.Context, upc string) (*model.Furniture, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *entityResolver) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	fmt.Println(id)
	return &model.User{
		ID: id,
	}, nil
}

func (r *entityResolver) FindVanByID(ctx context.Context, id string) (*model.Van, error) {
	panic(fmt.Errorf("not implemented"))
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
