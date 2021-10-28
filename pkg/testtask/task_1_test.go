package testtask

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jensneuse/graphql-go-tools/pkg/astparser"
)

func TestGatherStringFieldsStats(t *testing.T) {
	definition, report := astparser.ParseGraphqlDocumentString(StarWarsSchema)
	require.False(t, report.HasErrors())

	documentStats := GatherStringFieldsStats(&definition, &report)
	sort.Strings(documentStats.stringFieldNames)

	expectedStats := &StringFieldStats{
		stringFieldNames: []string{
			"commentary",
			"height",
			"name",
			"primaryFunction",
		},
		stringFieldCount: 7,
	}

	assert.Equal(t, expectedStats, documentStats)
}
