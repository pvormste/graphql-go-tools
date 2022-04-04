package astnormalization

import (
	"github.com/jensneuse/graphql-go-tools/pkg/ast"
	"github.com/jensneuse/graphql-go-tools/pkg/astvisitor"
)

func removeScalarDuplicates(walker *astvisitor.Walker) {
	visitor := removeScalarDuplicatesVisitor{
		Walker: walker,
	}
	walker.RegisterEnterDocumentVisitor(&visitor)
	walker.RegisterLeaveDocumentVisitor(&visitor)
	walker.RegisterEnterScalarTypeDefinitionVisitor(&visitor)
}

type removeScalarDuplicatesVisitor struct {
	*astvisitor.Walker
	operation     *ast.Document
	definition    *ast.Document
	scalarSet     map[string]bool
	refsToRemove  []int
	nodesToRemove []ast.Node
	lastRef       int
}

func (d *removeScalarDuplicatesVisitor) EnterDocument(operation, definition *ast.Document) {
	d.operation, d.definition = operation, definition
	d.scalarSet = make(map[string]bool)
	d.refsToRemove = make([]int, 0)
	d.nodesToRemove = make([]ast.Node, 0)
	d.lastRef = -1
}

// LeaveDocument original implementation
//func (d *removeScalarDuplicatesVisitor) LeaveDocument(operation, definition *ast.Document) {
//	if len(d.refsToRemove) < 1 {
//		return
//	}
//	newRootNodes := make([]ast.Node, 0, len(d.operation.RootNodes))
//MainLoop:
//	for _, node := range d.operation.RootNodes {
//		if node.Kind == ast.NodeKindScalarTypeDefinition {
//			for i, ref := range d.refsToRemove {
//				if node.Ref == ref {
//					lastIndex := len(d.refsToRemove) - 1
//					d.refsToRemove[i] = d.refsToRemove[lastIndex]
//					d.refsToRemove = d.refsToRemove[:lastIndex]
//					continue MainLoop
//				}
//			}
//		}
//		newRootNodes = append(newRootNodes, node)
//	}
//	d.operation.RootNodes = newRootNodes
//}

// LeaveDocument by altering original node
//func (d *removeScalarDuplicatesVisitor) LeaveDocument(operation, definition *ast.Document) {
//	if len(d.refsToRemove) < 1 {
//		return
//	}
//ParentLoop:
//	for j, node := range d.operation.RootNodes {
//		if node.Kind != ast.NodeKindScalarTypeDefinition {
//			continue
//		}
//		for i, ref := range d.refsToRemove {
//			if node.Ref == ref {
//				node.Kind = ast.NodeKindUnknown
//				node.Ref = -1
//				d.operation.RootNodes[j] = node
//				lastIndex := len(d.refsToRemove) - 1
//				if lastIndex < 1 {
//					break ParentLoop
//				}
//				d.refsToRemove[i] = d.refsToRemove[lastIndex]
//				d.refsToRemove = d.refsToRemove[:lastIndex]
//				continue ParentLoop
//			}
//		}
//	}
//}

// LeaveDocument using DeleteRootNodes
//func (d *removeScalarDuplicatesVisitor) LeaveDocument(operation, definition *ast.Document) {
//	if len(d.refsToRemove) < 1 {
//		return
//	}
//	nodesToRemove := make([]ast.Node, 0)
//ParentLoop:
//	for _, node := range d.operation.RootNodes {
//		if node.Kind != ast.NodeKindScalarTypeDefinition {
//			continue
//		}
//		for i, ref := range d.refsToRemove {
//			if node.Ref == ref {
//				nodesToRemove = append(nodesToRemove, node)
//				lastIndex := len(d.refsToRemove) - 1
//				if lastIndex < 1 {
//					break ParentLoop
//				}
//				d.refsToRemove[i] = d.refsToRemove[lastIndex]
//				d.refsToRemove = d.refsToRemove[:lastIndex]
//				continue ParentLoop
//			}
//		}
//	}
//	d.operation.DeleteRootNodes(nodesToRemove)
//}

// LeaveDocument using pointers
//func (d *removeScalarDuplicatesVisitor) LeaveDocument(operation, definition *ast.Document) {
//	if len(d.refsToRemove) < 1 {
//		return
//	}
//	nodesToRemove := make([]*ast.Node, 0, len(d.refsToRemove))
//ParentLoop:
//	for i, node := range d.operation.RootNodes {
//		if node.Kind != ast.NodeKindScalarTypeDefinition {
//			continue
//		}
//		for j, ref := range d.refsToRemove {
//			if node.Ref == ref {
//				nodesToRemove = append(nodesToRemove, &d.operation.RootNodes[i])
//				lastIndex := len(d.refsToRemove) - 1
//				if lastIndex < 1 {
//					break ParentLoop
//				}
//				d.refsToRemove[j] = d.refsToRemove[lastIndex]
//				d.refsToRemove = d.refsToRemove[:lastIndex]
//				continue ParentLoop
//			}
//		}
//	}
//	d.operation.DeleteRootNodesByPointer(nodesToRemove)
//}

//readable
func (d *removeScalarDuplicatesVisitor) LeaveDocument(operation, definition *ast.Document) {
	if len(d.nodesToRemove) < 1 {
		return
	}
	d.operation.DeleteRootNodes(d.nodesToRemove)
}

func (d *removeScalarDuplicatesVisitor) EnterScalarTypeDefinition(ref int) {
	if ref <= d.lastRef {
		return
	}
	name := d.operation.ScalarTypeDefinitionNameString(ref)
	if ok := d.scalarSet[name]; ok {
		d.refsToRemove = append(d.refsToRemove, ref)
		d.nodesToRemove = append(d.nodesToRemove, ast.Node{ast.NodeKindScalarTypeDefinition, ref})
	} else {
		d.scalarSet[name] = true
	}
	d.lastRef = ref
}

func (d removeScalarDuplicatesVisitor) LeaveScalarTypeDefinition(ref int) {
}
