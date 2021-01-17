package graphql

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jensneuse/graphql-go-tools/pkg/operationreport"
	"github.com/jensneuse/graphql-go-tools/pkg/starwars"
)

func TestRequest_ValidateForSchema(t *testing.T) {
	t.Run("should return error when schema is nil", func(t *testing.T) {
		request := Request{
			OperationName: "Hello",
			Variables:     nil,
			Query:         `query Hello { hello }`,
		}

		result, err := request.ValidateForSchema(nil)
		assert.Error(t, err)
		assert.Equal(t, ErrNilSchema, err)
		assert.Equal(t, ValidationResult{Valid: false, Errors: nil}, result)
	})

	t.Run("should return gql errors no valid operation is in the the request", func(t *testing.T) {
		request := Request{}

		schema, err := NewSchemaFromString("schema { query: Query } type Query { hello: String }")
		require.NoError(t, err)

		result, err := request.ValidateForSchema(schema)
		assert.NoError(t, err)
		assert.False(t, result.Valid)
		assert.Greater(t, result.Errors.Count(), 0)
	})

	t.Run("should return gql errors when validation fails", func(t *testing.T) {
		request := Request{
			OperationName: "Goodbye",
			Variables:     nil,
			Query:         `query Goodbye { goodbye }`,
		}

		schema, err := NewSchemaFromString("schema { query: Query } type Query { hello: String }")
		require.NoError(t, err)

		result, err := request.ValidateForSchema(schema)
		assert.NoError(t, err)
		assert.False(t, result.Valid)
		assert.Greater(t, result.Errors.Count(), 0)
	})

	t.Run("should successfully validate even when schema definition is missing", func(t *testing.T) {
		request := Request{
			OperationName: "Hello",
			Variables:     nil,
			Query:         `query Hello { hello }`,
		}

		schema, err := NewSchemaFromString("type Query { hello: String }")
		require.NoError(t, err)

		result, err := request.ValidateForSchema(schema)
		assert.NoError(t, err)
		assert.True(t, result.Valid)
		assert.Nil(t, result.Errors)
	})

	t.Run("should return valid result for introspection query after normalization", func(t *testing.T) {
		schema := starwarsSchema(t)
		request := requestForQuery(t, starwars.FileIntrospectionQuery)

		normalizationResult, err := request.Normalize(schema)
		require.NoError(t, err)
		require.True(t, normalizationResult.Successful)
		require.True(t, request.IsNormalized())

		result, err := request.ValidateForSchema(schema)
		assert.NoError(t, err)
		assert.True(t, result.Valid)
		assert.Nil(t, result.Errors)
	})

	t.Run("should return valid result when validation is successful", func(t *testing.T) {
		schema := starwarsSchema(t)
		request := requestForQuery(t, starwars.FileSimpleHeroQuery)

		result, err := request.ValidateForSchema(schema)
		assert.NoError(t, err)
		assert.True(t, result.Valid)
		assert.Nil(t, result.Errors)
	})
}

func TestRequest_ValidateRestrictedFields(t *testing.T) {
	t.Run("should return error when schema is nil", func(t *testing.T) {
		request := Request{}
		result, err := request.ValidateRestrictedFields(nil, nil)
		assert.Error(t, err)
		assert.Equal(t, ErrNilSchema, err)
		assert.False(t, result.Valid)
	})

	t.Run("should deny request when no allowed fields set", func(t *testing.T) {
		schema := starwarsSchema(t)
		request := requestForQuery(t, starwars.FileSimpleHeroQuery)

		result, err := request.ValidateRestrictedFields(schema, nil)
		assert.NoError(t, err)
		assert.False(t, result.Valid)
	})

	t.Run("when allowed fields set", func(t *testing.T) {
		schema := starwarsSchema(t)

		allowedFields := []Type{
			{Name: "Query", Fields: []string{"hero"}},
			{Name: "Character", Fields: []string{"name"}},
		}

		t.Run("should allow request", func(t *testing.T) {
			t.Run("when only allowed fields requested", func(t *testing.T) {
				request := requestForQuery(t, starwars.FileSimpleHeroQuery)
				result, err := request.ValidateRestrictedFields(schema, allowedFields)
				assert.NoError(t, err)
				assert.True(t, result.Valid)
				assert.Empty(t, result.Errors)

				request = requestForQuery(t, starwars.FileHeroWithAliasesQuery)
				result, err = request.ValidateRestrictedFields(schema, allowedFields)
				assert.NoError(t, err)
				assert.True(t, result.Valid)
				assert.Empty(t, result.Errors)
			})
		})

		t.Run("should disallow request", func(t *testing.T) {
			t.Run("when query is not allowed", func(t *testing.T) {
				request := requestForQuery(t, starwars.FileDroidWithArgAndVarQuery)
				result, err := request.ValidateRestrictedFields(schema, allowedFields)
				assert.NoError(t, err)
				assert.False(t, result.Valid)
				assert.Error(t, result.Errors)

				var buf bytes.Buffer
				_, _ = result.Errors.WriteResponse(&buf)
				assert.Equal(t, `{"errors":[{"message":"field: droid is not allowed on type: Query"}]}`, buf.String())
			})

			t.Run("when mutation is not allowed", func(t *testing.T) {
				request := requestForQuery(t, starwars.FileCreateReviewMutation)
				result, err := request.ValidateRestrictedFields(schema, allowedFields)
				assert.NoError(t, err)
				assert.False(t, result.Valid)
				assert.Error(t, result.Errors)

				var buf bytes.Buffer
				_, _ = result.Errors.WriteResponse(&buf)
				assert.Equal(t, `{"errors":[{"message":"type: Mutation is not allowed"}]}`, buf.String())
			})

			t.Run("when requested dissalowed field", func(t *testing.T) {
				request := requestForQuery(t, starwars.FileUnionQuery)
				result, err := request.ValidateRestrictedFields(schema, allowedFields)
				assert.NoError(t, err)
				assert.False(t, result.Valid)
				assert.Error(t, result.Errors)

				var buf bytes.Buffer
				_, _ = result.Errors.WriteResponse(&buf)
				assert.Equal(t, `{"errors":[{"message":"field: friends is not allowed on type: Character"}]}`, buf.String())
			})

			t.Run("when mutation response type has not allowed field", func(t *testing.T) {
				allowedFields := []Type{
					{Name: "Mutation", Fields: []string{"createReview"}},
					{Name: "Review", Fields: []string{"id", "stars"}},
				}

				request := requestForQuery(t, starwars.FileCreateReviewMutation)
				result, err := request.ValidateRestrictedFields(schema, allowedFields)
				assert.NoError(t, err)
				assert.False(t, result.Valid)
				assert.Error(t, result.Errors)

				var buf bytes.Buffer
				_, _ = result.Errors.WriteResponse(&buf)
				assert.Equal(t, `{"errors":[{"message":"field: commentary is not allowed on type: Review"}]}`, buf.String())
			})
		})
	})

}

func Test_operationValidationResultFromReport(t *testing.T) {
	t.Run("should return result for valid when report does not have errors", func(t *testing.T) {
		report := operationreport.Report{}
		result, err := operationValidationResultFromReport(report)

		assert.NoError(t, err)
		assert.Equal(t, ValidationResult{Valid: true, Errors: nil}, result)
	})

	t.Run("should return validation error and internal error when report contain them", func(t *testing.T) {
		internalErr := errors.New("errors occurred")
		externalErr := operationreport.ExternalError{
			Message:   "graphql error",
			Path:      nil,
			Locations: nil,
		}

		report := operationreport.Report{}
		report.AddInternalError(internalErr)
		report.AddExternalError(externalErr)

		result, err := operationValidationResultFromReport(report)

		assert.Error(t, err)
		assert.Equal(t, internalErr, err)
		assert.False(t, result.Valid)
		assert.Len(t, result.Errors.(OperationValidationErrors), 1)
		assert.Equal(t, "graphql error", result.Errors.(OperationValidationErrors)[0].Message)
	})
}
