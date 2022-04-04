package astnormalization

import (
	"github.com/jensneuse/graphql-go-tools/pkg/ast"
	"testing"
)

func TestFoo(t *testing.T) {
	t.Run("Single duplicate scalar is removed", func(t *testing.T) {
		run(removeScalarDuplicates, "", testDataOne, "scalar ScalarOne")
	})
	t.Run("Single duplicate scalar among other data is removed", func(t *testing.T) {
		run(removeScalarDuplicates, "", testDataTwo, "scalar ScalarOne type TypeOne scalar ScalarTwo")
	})
	t.Run("Several duplicate scalars among other data are removed", func(t *testing.T) {
		run(removeScalarDuplicates, "", testDataThree, "type TypeOne scalar ScalarOne type TypeTwo scalar ScalarTwo scalar ScalarThree type TypeThree scalar ScalarFour")
	})
}

func BenchmarkScalarPointers(b *testing.B) {
	for r := 0; r < b.N; r++ {
		refsToRemove := []int{1, 3, 4, 7, 8}
		rootNodes := []ast.Node{
			{ast.NodeKindObjectTypeDefinition, 0},
			{ast.NodeKindScalarTypeDefinition, 0},
			{ast.NodeKindObjectTypeDefinition, 1},
			{ast.NodeKindScalarTypeDefinition, 1},
			{ast.NodeKindScalarTypeDefinition, 2},
			{ast.NodeKindScalarTypeDefinition, 3},
			{ast.NodeKindScalarTypeDefinition, 4},
			{ast.NodeKindScalarTypeDefinition, 5},
			{ast.NodeKindObjectTypeDefinition, 2},
			{ast.NodeKindScalarTypeDefinition, 6},
			{ast.NodeKindScalarTypeDefinition, 7},
			{ast.NodeKindScalarTypeDefinition, 8},
		}
		nodesToRemove := make([]*ast.Node, 0, len(refsToRemove))
	Bone:
		for i, node := range rootNodes {
			if node.Kind != ast.NodeKindScalarTypeDefinition {
				continue
			}
			for j, ref := range refsToRemove {
				if node.Ref == ref {
					nodesToRemove = append(nodesToRemove, &rootNodes[i])
					lastIndex := len(refsToRemove) - 1
					refsToRemove[j] = refsToRemove[lastIndex]
					refsToRemove = refsToRemove[:lastIndex]
					continue Bone
				}
			}
		}
		// New implementation of DeleteRootNodes
		for i := len(nodesToRemove) - 1; i > -1; i-- {
			for j := len(rootNodes) - 1; j > -1; j-- {
				if rootNodes[j] == *nodesToRemove[i] {
					rootNodes = append(rootNodes[:j], rootNodes[j+1:]...)
					break
				}
			}
		}
	}
}

func BenchmarkScalarDeleteNodes(b *testing.B) {
	for r := 0; r < b.N; r++ {
		refsToRemove := []int{1, 3, 4, 7, 8}
		rootNodes := []ast.Node{
			{ast.NodeKindObjectTypeDefinition, 0},
			{ast.NodeKindScalarTypeDefinition, 0},
			{ast.NodeKindObjectTypeDefinition, 1},
			{ast.NodeKindScalarTypeDefinition, 1},
			{ast.NodeKindScalarTypeDefinition, 2},
			{ast.NodeKindScalarTypeDefinition, 3},
			{ast.NodeKindScalarTypeDefinition, 4},
			{ast.NodeKindScalarTypeDefinition, 5},
			{ast.NodeKindObjectTypeDefinition, 2},
			{ast.NodeKindScalarTypeDefinition, 6},
			{ast.NodeKindScalarTypeDefinition, 7},
			{ast.NodeKindScalarTypeDefinition, 8},
		}
		nodesToRemove := make([]ast.Node, 0)
	Btwo:
		for _, node := range rootNodes {
			if node.Kind != ast.NodeKindScalarTypeDefinition {
				continue
			}
			for i, ref := range refsToRemove {
				if node.Ref == ref {
					nodesToRemove = append(nodesToRemove, node)
					lastIndex := len(refsToRemove) - 1
					refsToRemove[i] = refsToRemove[lastIndex]
					refsToRemove = refsToRemove[:lastIndex]
					continue Btwo
				}
			}
		}
		// ast/DeleteRootNodes
		for _, node := range nodesToRemove {
			for j := range rootNodes {
				if rootNodes[j].Kind == node.Kind && rootNodes[j].Ref == node.Ref {
					rootNodes = append(rootNodes[:j], rootNodes[j+1:]...)
					break
				}
			}
		}
	}
}

func BenchmarkScalarAlterNode(b *testing.B) {
	for r := 0; r < b.N; r++ {
		refsToRemove := []int{1, 3, 4, 7, 8}
		rootNodes := []ast.Node{
			{ast.NodeKindObjectTypeDefinition, 0},
			{ast.NodeKindScalarTypeDefinition, 0},
			{ast.NodeKindObjectTypeDefinition, 1},
			{ast.NodeKindScalarTypeDefinition, 1},
			{ast.NodeKindScalarTypeDefinition, 2},
			{ast.NodeKindScalarTypeDefinition, 3},
			{ast.NodeKindScalarTypeDefinition, 4},
			{ast.NodeKindScalarTypeDefinition, 5},
			{ast.NodeKindObjectTypeDefinition, 2},
			{ast.NodeKindScalarTypeDefinition, 6},
			{ast.NodeKindScalarTypeDefinition, 7},
			{ast.NodeKindScalarTypeDefinition, 8},
		}
	Bthree:
		for j, node := range rootNodes {
			if node.Kind != ast.NodeKindScalarTypeDefinition {
				continue
			}
			for i, ref := range refsToRemove {
				if node.Ref == ref {
					node.Kind = ast.NodeKindUnknown
					node.Ref = -1
					rootNodes[j] = node
					lastIndex := len(refsToRemove) - 1
					refsToRemove[i] = refsToRemove[lastIndex]
					refsToRemove = refsToRemove[:lastIndex]
					continue Bthree
				}
			}
		}
	}
}

func BenchmarkScalarOriginalImplementation(b *testing.B) {
	for r := 0; r < b.N; r++ {
		refsToRemove := []int{1, 3, 4, 7, 8}
		rootNodes := []ast.Node{
			{ast.NodeKindObjectTypeDefinition, 0},
			{ast.NodeKindScalarTypeDefinition, 0},
			{ast.NodeKindObjectTypeDefinition, 1},
			{ast.NodeKindScalarTypeDefinition, 1},
			{ast.NodeKindScalarTypeDefinition, 2},
			{ast.NodeKindScalarTypeDefinition, 3},
			{ast.NodeKindScalarTypeDefinition, 4},
			{ast.NodeKindScalarTypeDefinition, 5},
			{ast.NodeKindObjectTypeDefinition, 2},
			{ast.NodeKindScalarTypeDefinition, 6},
			{ast.NodeKindScalarTypeDefinition, 7},
			{ast.NodeKindScalarTypeDefinition, 8},
		}
		newRootNodes := make([]ast.Node, 0, len(rootNodes))
	Bfour:
		for _, node := range rootNodes {
			if node.Kind == ast.NodeKindScalarTypeDefinition {
				for i, ref := range refsToRemove {
					if node.Ref == ref {
						lastIndex := len(refsToRemove) - 1
						refsToRemove[i] = refsToRemove[lastIndex]
						refsToRemove = refsToRemove[:lastIndex]
						continue Bfour
					}
				}
			}
			newRootNodes = append(newRootNodes, node)
		}
		rootNodes = newRootNodes
	}
}

const testDataOne = `
scalar ScalarOne
scalar ScalarOne
`
const testDataTwo = `
scalar ScalarOne
type TypeOne
scalar ScalarOne
scalar ScalarOne
scalar ScalarTwo
`
const testDataThree = `
type TypeOne
scalar ScalarOne
type TypeTwo
scalar ScalarOne
scalar ScalarTwo
scalar ScalarOne
scalar ScalarTwo
scalar ScalarThree
type TypeThree
scalar ScalarFour
scalar ScalarThree
scalar ScalarOne
`
