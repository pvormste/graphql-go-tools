package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/jensneuse/graphql-go-tools/examples/federation/products/graph/generated"
	"github.com/jensneuse/graphql-go-tools/examples/federation/products/graph/model"
)

func (r *queryResolver) Product(ctx context.Context, upc string) (model.Product, error) {
	for i := 0; i < len(productsDB); i++ {
		switch p := productsDB[i].(type) {
		case *model.Furniture:
			if p.Upc == upc {
				return p, nil
			}
		}
	}

	return nil, nil
}

func (r *queryResolver) Vehicle(ctx context.Context, id string) (model.Vehicle, error) {
	for i := 0; i < len(vehiclesDB); i++ {
		switch v := vehiclesDB[i].(type) {
		case model.Car:
			if v.ID == id {
				return v, nil
			}
		case model.Van:
			if v.ID == id {
				return v, nil
			}
		}
	}

	return nil, nil
}

func (r *queryResolver) TopProducts(ctx context.Context, first *int) ([]model.Product, error) {
	const defaultTop = 5
	productsCount := len(productsDB)
	if first == nil && productsCount <= defaultTop {
		return productsDB, nil
	} else if first == nil && productsCount > defaultTop {
		return productsDB[0:defaultTop], nil
	} else if productsCount < *first {
		return productsDB, nil
	}

	return productsDB[0:*first], nil
}

func (r *queryResolver) TopCars(ctx context.Context, first *int) ([]*model.Car, error) {
	const defaultTop = 5
	cars := make([]*model.Car, 0)
	for _, vehicle := range vehiclesDB {
		switch v := vehicle.(type) {
		case *model.Car:
			cars = append(cars, v)
		}
	}

	carsCount := len(cars)
	if first == nil && carsCount <= defaultTop {
		return cars, nil
	} else if first == nil && carsCount > defaultTop {
		return cars[0:defaultTop], nil
	} else if carsCount < *first {
		return cars, nil
	}

	return cars[0:*first], nil
}

func (r *subscriptionResolver) UpdatedPrice(ctx context.Context) (<-chan model.Product, error) {
	updatedPrice := make(chan model.Product)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(updateInterval):
				rand.Seed(time.Now().UnixNano())
				product := productsDB[0]
				price := currentPrice

				if randomnessEnabled {
					product = productsDB[rand.Intn(len(productsDB)-1)]
					price = rand.Intn(maxPrice-minPrice+1) + minPrice
				} else {
					currentPrice += 1
				}

				switch p := product.(type) {
				case *model.Furniture:
					p.Price = stringPtr(strconv.Itoa(price))
					updatedPrice <- product
				}
			}
		}
	}()
	return updatedPrice, nil
}

func (r *subscriptionResolver) UpdateProductPrice(ctx context.Context, upc string) (<-chan model.Product, error) {
	updatedPrice := make(chan model.Product)
	var product model.Product

	for _, productEntity := range productsDB {
		switch p := productEntity.(type) {
		case *model.Furniture:
			if p.Upc == upc {
				product = p
			}
		}
	}

	if product == nil {
		return nil, fmt.Errorf("unknown product upc: %s", upc)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
				rand.Seed(time.Now().UnixNano())
				min := 50
				max := 2000

				switch p := product.(type) {
				case *model.Furniture:
					randPrice := rand.Intn(max-min+1) + min
					p.Price = stringPtr(strconv.Itoa(randPrice))
					updatedPrice <- p
				}

			}
		}
	}()

	return updatedPrice, nil
}

func (r *subscriptionResolver) Stock(ctx context.Context) (<-chan []model.Product, error) {
	stock := make(chan []model.Product)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(2 * time.Second):
				rand.Seed(time.Now().UnixNano())
				randIndex := rand.Intn(len(productsDB))

				switch product := productsDB[randIndex].(type) {
				case *model.Furniture:
					if product.InStock > 0 {
						product.InStock--
					}
				}

				stock <- productsDB
			}
		}
	}()

	return stock, nil
}

func (r *userResolver) Vehicle(ctx context.Context, obj *model.User) (model.Vehicle, error) {
	if obj == nil {
		return nil, nil
	}

	for _, vehicle := range vehiclesDB {
		switch v := vehicle.(type) {
		case *model.Car:
			if v.ID == obj.ID {
				return v, nil
			}
		case *model.Van:
			if v.ID == obj.ID {
				return v, nil
			}
		}
	}

	return nil, nil
}

func (r *userResolver) Thing(ctx context.Context, obj *model.User) (model.Thing, error) {
	if obj == nil {
		return nil, nil
	}

	for _, vehicle := range vehiclesDB {
		switch v := vehicle.(type) {
		case *model.Car:
			if v.ID == obj.ID {
				return v, nil
			}
		}
	}

	return nil, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
