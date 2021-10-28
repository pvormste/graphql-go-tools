package testtask

import (
	"github.com/jensneuse/graphql-go-tools/pkg/ast"
	"github.com/jensneuse/graphql-go-tools/pkg/astvisitor"
	"github.com/jensneuse/graphql-go-tools/pkg/operationreport"
)

type DocumentStats struct {
	uniqFieldNames   []string
	objectTypesNames []string
	stringFieldCount int
	enumValues       []string
}

func GatherDocumentStats(doc *ast.Document, report *operationreport.Report) *DocumentStats {
	walker := astvisitor.NewWalker(48)
	visitor := &DocumentStatsVisitor{
		Walker: &walker,
	}

	walker.RegisterEnterDocumentVisitor(visitor)

	// register additional walk methods here
	walker.RegisterEnterEnumValueDefinitionVisitor(visitor)

	// run walker
	walker.Walk(doc, nil, report)

	// obtain results

	return &DocumentStats{
		enumValues: visitor.enumValues,
	}
}

type DocumentStatsVisitor struct {
	*astvisitor.Walker
	definition *ast.Document
	enumValues []string
}

func (v *DocumentStatsVisitor) EnterEnumValueDefinition(ref int) {
	v.enumValues = append(v.enumValues, v.definition.EnumValueDefinitionNameString(ref))
}

func (v *DocumentStatsVisitor) EnterDocument(operation, _ *ast.Document) {
	v.definition = operation
	v.enumValues = make([]string, 0, 3)
}
