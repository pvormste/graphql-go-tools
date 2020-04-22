package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jensneuse/graphql-go-tools/pkg/introspection"
	"github.com/jensneuse/graphql-go-tools/pkg/starwars"
)

func TestSchemaFieldsGenerator_Generate(t *testing.T) {
	generator := schemaFieldsGenerator{}

	t.Run("It should return an error when schema is nil", func(t *testing.T) {
		result, err := generator.Generate("")

		assert.Nil(t, result.Types)
		assert.Equal(t, ErrEmptySchema, err)
	})

	t.Run("It should return types for valid schema", func(t *testing.T) {
		expectedTypes := []SchemaType{
			{Name: "Query", Type: "OBJECT",
				Fields: []SchemaField{
					{Name: "hero", Type: "Character", TypeRef: "Character"},
					{Name: "droid", Type: "Droid", TypeRef: "Droid"},
					{Name: "search", Type: "SearchResult", TypeRef: "UNION"},
				}},
			{Name: "Mutation", Type: "OBJECT",
				Fields: []SchemaField{
					{Name: "createReview", Type: "Review", TypeRef: "Review"},
				}},
			{Name: "Subscription", Type: "OBJECT",
				Fields: []SchemaField{
					{Name: "remainingJedis", Type: "Int!", TypeRef: "SCALAR"},
				}},
			{Name: "Review", Type: "OBJECT",
				Fields: []SchemaField{
					{Name: "id", Type: "ID!", TypeRef: "SCALAR"},
					{Name: "stars", Type: "Int!", TypeRef: "SCALAR"},
					{Name: "commentary", Type: "String", TypeRef: "SCALAR"},
				}},
			{Name: "Character", Type: "INTERFACE",
				Fields: []SchemaField{
					{Name: "name", Type: "String!", TypeRef: "SCALAR"},
					{Name: "friends", Type: "[Character]", TypeRef: "Character"},
				}},
			{Name: "Human", Type: "OBJECT",
				Fields: []SchemaField{
					{Name: "name", Type: "String!", TypeRef: "SCALAR"},
					{Name: "height", Type: "String!", TypeRef: "SCALAR"},
					{Name: "friends", Type: "[Character]", TypeRef: "Character"},
				}},
			{Name: "Droid", Type: "OBJECT",
				Fields: []SchemaField{
					{Name: "name", Type: "String!", TypeRef: "SCALAR"},
					{Name: "primaryFunction", Type: "String!", TypeRef: "SCALAR"},
					{Name: "friends", Type: "[Character]", TypeRef: "Character"},
				}},
			{Name: "Startship", Type: "OBJECT",
				Fields: []SchemaField{
					{Name: "name", Type: "String!", TypeRef: "SCALAR"},
					{Name: "length", Type: "Float!", TypeRef: "SCALAR"},
				}},
		}
		starwars.SetRelativePathToStarWarsPackage("../starwars")
		schemaBytes := starwars.Schema(t)

		result, err := generator.Generate(string(schemaBytes))
		assert.NoError(t, err)
		assert.Nil(t, result.Errors)
		assert.Equal(t, expectedTypes, result.Types)
	})
}

func Test_getType(t *testing.T) {
	tRef := introspection.TypeRef{
		Kind: introspection.NONNULL,
		OfType: &introspection.TypeRef{
			Kind: introspection.LIST,
			OfType: &introspection.TypeRef{
				Kind: introspection.OBJECT,
				Name: strPtr("Droid"),
			},
		},
	}

	fullType, underlyingType := fieldType(tRef)

	assert.Equal(t, "[Droid]!", fullType)
	assert.Equal(t, "Droid", underlyingType)
}

func strPtr(in string) *string {
	return &in
}
