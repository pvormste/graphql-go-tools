package testtask

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jensneuse/graphql-go-tools/pkg/astparser"
)

func TestGatherDocumentStats(t *testing.T) {
	definition, report := astparser.ParseGraphqlDocumentString(StarWarsSchema)
	require.False(t, report.HasErrors())

	documentStats := GatherDocumentStats(&definition, &report)
	sort.Strings(documentStats.uniqFieldNames)
	sort.Strings(documentStats.objectTypesNames)
	sort.Strings(documentStats.enumValues)

	expectedStats := &DocumentStats{
		uniqFieldNames: []string{
			"commentary",
			"createReview",
			"droid",
			"friends",
			"height",
			"hero",
			"id",
			"length",
			"name",
			"primaryFunction",
			"remainingJedis",
			"search",
			"stars",
		},
		objectTypesNames: []string{
			"Droid",
			"Human",
			"Mutation",
			"Query",
			"Review",
			"Starship",
			"Subscription",
		},
		enumValues: []string{
			"EMPIRE",
			"JEDI",
			"NEWHOPE",
		},
	}

	assert.Equal(t, expectedStats, documentStats)
}
