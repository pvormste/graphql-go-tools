package introspection_datasource

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/sebdah/goldie"
	"github.com/stretchr/testify/require"

	"github.com/pvormste/graphql-go-tools/pkg/astparser"
	"github.com/pvormste/graphql-go-tools/pkg/asttransform"
	"github.com/pvormste/graphql-go-tools/pkg/introspection"
)

func TestSource_Load(t *testing.T) {
	def, report := astparser.ParseGraphqlDocumentString(testSchema)
	require.False(t, report.HasErrors())
	require.NoError(t, asttransform.MergeDefinitionWithBaseSchema(&def))

	var data introspection.Data
	gen := introspection.NewGenerator()
	gen.Generate(&def, &report, &data)
	require.False(t, report.HasErrors())

	run := func(input string, fixtureName string) func(t *testing.T) {
		t.Helper()
		return func(t *testing.T) {
			buf := &bytes.Buffer{}
			source := &Source{introspectionData: &data}
			require.NoError(t, source.Load(context.Background(), []byte(input), buf))

			actualResponse := &bytes.Buffer{}
			require.NoError(t, json.Indent(actualResponse, buf.Bytes(), "", "  "))
			goldie.Assert(t, fixtureName, actualResponse.Bytes())
		}
	}

	t.Run("schema introspection", run(`{"request_type":1}`, `schema_introspection`))
	t.Run("type introspection", run(`{"request_type":2,"type_name":"Query"}`, `type_introspection`))
	t.Run("type introspection of not existing type", run(`{"request_type":2,"type_name":"NotExisting"}`, `not_existing_type`))

	t.Run("type fields", func(t *testing.T) {
		t.Run("include deprecated", run(`{"request_type":3,"on_type_name":"Query","include_deprecated":true}`, `fields_with_deprecated`))

		t.Run("no deprecated", run(`{"request_type":3,"on_type_name":"Query","include_deprecated":false}`, `fields_without_deprecated`))

		t.Run("of not existing type", run(`{"request_type":3,"on_type_name":"NotExisting","include_deprecated":true}`, `not_existing_type`))
	})

	t.Run("type enum values", func(t *testing.T) {
		t.Run("include deprecated", run(`{"request_type":4,"on_type_name":"Episode","include_deprecated":true}`, `enum_values_with_deprecated`))

		t.Run("no deprecated", run(`{"request_type":4,"on_type_name":"Episode","include_deprecated":false}`, `enum_values_without_deprecated`))

		t.Run("of not existing type", run(`{"request_type":4,"on_type_name":"NotExisting","include_deprecated":true}`, `not_existing_type`))
	})
}

const testSchema = `
schema {
    query: Query
}

type Query {
    me: Droid @deprecated
    droid(id: ID!): Droid
}

enum Episode {
    NEWHOPE
    EMPIRE
    JEDI @deprecated
}

type Droid {
    name: String!
}
`
