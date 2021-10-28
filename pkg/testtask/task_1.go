package testtask

import (
	"github.com/jensneuse/graphql-go-tools/pkg/ast"
	"github.com/jensneuse/graphql-go-tools/pkg/astvisitor"
	"github.com/jensneuse/graphql-go-tools/pkg/operationreport"
)

const StarWarsSchema = `
union SearchResult = Human | Droid | Starship

schema {
    query: Query
    mutation: Mutation
    subscription: Subscription
}

type Query {
    hero: Character
    droid(id: ID!): Droid
    search(name: String!): SearchResult
}

type Mutation {
    createReview(episode: Episode!, review: ReviewInput!): Review
}

type Subscription {
    remainingJedis: Int!
}

input ReviewInput {
    stars: Int!
    commentary: String
}

type Review {
    id: ID!
    stars: Int!
    commentary: String
}

enum Episode {
    NEWHOPE
    EMPIRE
    JEDI
}

interface Character {
    name: String!
    friends: [Character]
}

type Human implements Character {
    name: String!
    height: String!
    friends: [Character]
}

type Droid implements Character {
    name: String!
    primaryFunction: String!
    friends: [Character]
}

type Starship {
    name: String!
    length: Float!
}`

type StringFieldStats struct {
	stringFieldNames []string
	stringFieldCount int
}

func GatherStringFieldsStats(doc *ast.Document, report *operationreport.Report) *StringFieldStats {
	walker := astvisitor.NewWalker(48)
	visitor := &StringFieldsStats{
		Walker: &walker,
	}

	walker.RegisterEnterDocumentVisitor(visitor)
	walker.RegisterEnterFieldDefinitionVisitor(visitor)

	// run walker
	walker.Walk(doc, nil, report)

	// obtain results

	var fieldNames []string
	for s, _ := range visitor.fieldNames {
		fieldNames = append(fieldNames, s)
	}

	return &StringFieldStats{
		stringFieldNames: fieldNames,
		stringFieldCount: visitor.fieldCount,
	}
}

type StringFieldsStats struct {
	*astvisitor.Walker
	definition *ast.Document
	fieldNames map[string]struct{}
	fieldCount int
}

func (v *StringFieldsStats) EnterFieldDefinition(ref int) {
	fieldTypeRef := v.definition.FieldDefinitionType(ref)
	fieldType := v.definition.Types[fieldTypeRef]

	switch fieldType.TypeKind {
	case ast.TypeKindNamed:
		if v.definition.TypeNameString(fieldTypeRef) == "Int" {
			v.fieldCount++
		}
	case ast.TypeKindNonNull:
		if v.definition.TypeNameString(fieldType.OfType) == "String" {
			v.fieldCount++
			v.fieldNames[v.definition.FieldDefinitionNameString(ref)] = struct{}{}
		}
	}
}

func (v *StringFieldsStats) EnterDocument(operation, _ *ast.Document) {
	v.definition = operation
	v.fieldNames = make(map[string]struct{})
}
