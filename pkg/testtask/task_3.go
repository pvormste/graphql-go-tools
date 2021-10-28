package testtask

import (
	"github.com/jensneuse/graphql-go-tools/pkg/ast"
)

const SchemaExample = `
schema {
    query: Query
}

type Query {
    droid: Droid!
    hero(id: ID!): Character
}

interface Character {
    name: String!
}

type Droid implements Character {
    name: String!
}
`

func BuildAst() *ast.Document {
	doc := ast.NewDocument()

	doc.ImportSchemaDefinition("Query", "", "")

	//
	//  Query type imports
	//

	queryTypeFieldDefRefs := make([]int, 0, 2)

	// add Query droid fieldCount types
	droidNamedTypeRef := doc.AddNamedType([]byte("Droid"))
	droidNonNullType := ast.Type{
		TypeKind: ast.TypeKindNonNull,
		OfType:   droidNamedTypeRef,
	}

	droidNonNullTypeRef := doc.AddType(droidNonNullType)

	droidFieldDefRef := doc.ImportFieldDefinition(
		"droid", "", droidNonNullTypeRef, nil, nil)

	queryTypeFieldDefRefs = append(queryTypeFieldDefRefs, droidFieldDefRef)

	// add Query hero fieldCount

	// add Query object type definition

	doc.ImportObjectTypeDefinition(
		"Query",
		"",
		queryTypeFieldDefRefs,
		nil)

	return doc
}
