package astnormalization

import (
	"github.com/jensneuse/graphql-go-tools/pkg/ast"
	"github.com/jensneuse/graphql-go-tools/pkg/astvisitor"
)

func removeScalarDuplicates(walker *astvisitor.Walker) {
	visitor := removeScalarDuplicatesVisitor{
		Walker: walker,
	}
	walker.RegisterScalarTypeDefinitionVisitor(&visitor)
	walker.RegisterEnterDocumentVisitor(&visitor)
	walker.RegisterLeaveDocumentVisitor(&visitor)
	walker.RegisterEnterScalarTypeDefinitionVisitor(&visitor)
}

type removeScalarDuplicatesVisitor struct {
	*astvisitor.Walker
	operation    *ast.Document
	definition   *ast.Document
	typesSeen    map[string]bool
	refsToRemove []int
	lastRef      int
}

func (d *removeScalarDuplicatesVisitor) EnterDocument(operation, definition *ast.Document) {
	d.operation, d.definition = operation, definition
	d.typesSeen = make(map[string]bool)
	d.refsToRemove = make([]int, 0)
	d.lastRef = -1
}

func (d *removeScalarDuplicatesVisitor) LeaveDocument(operation, definition *ast.Document) {
	if len(d.refsToRemove) < 1 {
		return
	}
	d.operation, d.definition = operation, definition
	newRootNodes := make([]ast.Node, 0)
MainLoop:
	for _, node := range d.operation.RootNodes {
		if node.Kind == ast.NodeKindScalarTypeDefinition {
			for i, ref := range d.refsToRemove {
				if node.Ref == ref {
					lastIndex := len(d.refsToRemove) - 1
					d.refsToRemove[i] = d.refsToRemove[lastIndex]
					d.refsToRemove = d.refsToRemove[:lastIndex]
					continue MainLoop
				}
			}
		}
		newRootNodes = append(newRootNodes, node)
	}
	d.operation.RootNodes = newRootNodes
}

func (d *removeScalarDuplicatesVisitor) EnterScalarTypeDefinition(ref int) {
	if ref <= d.lastRef {
		return
	}
	name := d.operation.ScalarTypeDefinitionNameString(ref)
	if ok := d.typesSeen[name]; ok {
		d.refsToRemove = append(d.refsToRemove, ref)
	} else {
		d.typesSeen[name] = true
	}
	d.lastRef = ref
}

func (d removeScalarDuplicatesVisitor) LeaveScalarTypeDefinition(ref int) {
}
