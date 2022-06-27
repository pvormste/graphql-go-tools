package ast

import "github.com/pvormste/graphql-go-tools/pkg/lexer/position"

type SchemaExtension struct {
	ExtendLiteral position.Position
	SchemaDefinition
}
