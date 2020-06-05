package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocument_GenerateUnusedVariableDefinitionName(t *testing.T) {
	doc := &Document{}
	runCase := func(expected string, existing ...string) func(t *testing.T) {
		return func(t *testing.T) {
			actual := string(doc.generateUniqueShortIdentifier(func(b []byte) bool {
				for i := range existing {
					if string(b) == existing[i] {
						return true
					}
				}
				return false
			}))
			assert.Equal(t, expected, actual)
		}
	}

	t.Run("empty -> a", runCase("a"))
	t.Run("existing -> d", runCase("d", "a", "b", "c"))
	all := make([]string, len(alphabet))
	for i := range alphabet {
		all[i] = alphabet[i : i+1]
	}
	t.Run("all -> aa", runCase("aa", all...))
}
