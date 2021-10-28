package testtask

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jensneuse/graphql-go-tools/pkg/astprinter"
)

func TestBuildAst(t *testing.T) {
	doc := BuildAst()

	resultingSchema, err := astprinter.PrintStringIndent(doc, nil, "  ")
	require.NoError(t, err)
	require.Equal(t, SchemaExample, resultingSchema)
}
