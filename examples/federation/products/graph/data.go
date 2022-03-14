package graph

import (
	"github.com/jensneuse/graphql-go-tools/examples/federation/products/graph/model"
)

var productsDB = []model.Product{
	&model.Furniture{
		Upc:   "1",
		Sku:   "TABLE1",
		Name:  stringPtr("Table"),
		Price: stringPtr("899"),
		Brand: model.Ikea{
			Asile: intPtr(10),
		},
		Metadata: []model.MetadataOrError{
			model.KeyValue{
				Key:   "Condition",
				Value: "excellent",
			},
		},
		Details: nil,
		InStock: 10000,
	},
	&model.Furniture{
		Upc:   "2",
		Sku:   "COUCH1",
		Name:  stringPtr("Couch"),
		Price: stringPtr("1299"),
		Brand: model.Amazon{
			Referrer: stringPtr("https://canopy.co"),
		},
		Metadata: []model.MetadataOrError{
			model.KeyValue{
				Key:   "Condition",
				Value: "used",
			},
		},
		Details: nil,
		InStock: 750,
	},
	&model.Furniture{
		Upc:   "3",
		Sku:   "CHAIR1",
		Name:  stringPtr("Chair"),
		Price: stringPtr("54"),
		Brand: model.Ikea{
			Asile: intPtr(10),
		},
		Metadata: []model.MetadataOrError{
			model.KeyValue{
				Key:   "Condition",
				Value: "like new",
			},
		},
		Details: nil,
		InStock: 2200,
	},
}

var vehiclesDB = []model.Vehicle{
	&model.Car{
		ID:          "1",
		Description: stringPtr("Humble Toyota"),
		Price:       stringPtr("9990"),
	},
	&model.Car{
		ID:          "2",
		Description: stringPtr("Awesome Tesla"),
		Price:       stringPtr("12990"),
	},
	&model.Van{
		ID:          "3",
		Description: stringPtr("Just a van..."),
		Price:       stringPtr("15990"),
	},
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
