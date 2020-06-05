package ast

import (
	"github.com/jensneuse/graphql-go-tools/internal/pkg/unsafebytes"
	"github.com/jensneuse/graphql-go-tools/pkg/lexer/position"
)

type OperationType int

const (
	OperationTypeUnknown OperationType = iota
	OperationTypeQuery
	OperationTypeMutation
	OperationTypeSubscription
)

type OperationDefinition struct {
	OperationType          OperationType      // one of query, mutation, subscription
	OperationTypeLiteral   position.Position  // position of the operation type literal, if present
	Name                   ByteSliceReference // optional, user defined name of the operation
	HasVariableDefinitions bool
	VariableDefinitions    VariableDefinitionList // optional, e.g. ($devicePicSize: Int)
	HasDirectives          bool
	Directives             DirectiveList // optional, e.g. @foo
	SelectionSet           int           // e.g. {field}
	HasSelections          bool
}

func (d *Document) OperationDefinitionNameBytes(ref int) ByteSlice {
	return d.Input.ByteSlice(d.OperationDefinitions[ref].Name)
}

func (d *Document) OperationDefinitionNameString(ref int) string {
	return unsafebytes.BytesToString(d.Input.ByteSlice(d.OperationDefinitions[ref].Name))
}

func (d *Document) AddOperationDefinitionToRootNodes(definition OperationDefinition) Node {
	d.OperationDefinitions = append(d.OperationDefinitions, definition)
	node := Node{Kind: NodeKindOperationDefinition, Ref: len(d.OperationDefinitions) - 1}
	d.RootNodes = append(d.RootNodes, node)
	return node
}

func (d *Document) AddVariableDefinitionToOperationDefinition(operationDefinitionRef, variableValueRef, typeRef int) {
	if !d.OperationDefinitions[operationDefinitionRef].HasVariableDefinitions {
		d.OperationDefinitions[operationDefinitionRef].HasVariableDefinitions = true
		d.OperationDefinitions[operationDefinitionRef].VariableDefinitions.Refs = d.Refs[d.NextRefIndex()][:0]
	}
	variableDefinition := VariableDefinition{
		VariableValue: Value{
			Kind: ValueKindVariable,
			Ref:  variableValueRef,
		},
		Type: typeRef,
	}
	d.VariableDefinitions = append(d.VariableDefinitions, variableDefinition)
	ref := len(d.VariableDefinitions) - 1
	d.OperationDefinitions[operationDefinitionRef].VariableDefinitions.Refs =
		append(d.OperationDefinitions[operationDefinitionRef].VariableDefinitions.Refs, ref)
}

const (
	alphabet     = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
	alphaNumeric = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`
)

func (d *Document) generateUniqueShortIdentifier(exists func(b []byte) bool) []byte {
	out := make([]byte, 0, 1)
	input := alphabet
	for i := 1; true; i++ {
		for j := 0; j < len(input); j++ {
			out = append(out, input[j])
			if exists(out) {
				out = out[:i-1]
				continue
			}
			return out
		}
		out = append(out, input[i-1])
		if i == 1 {
			input = alphaNumeric
		}
	}
	return nil
}

func (d *Document) GenerateUnusedVariableDefinitionName(operationDefinition int) []byte {
	return d.generateUniqueShortIdentifier(func(b []byte) bool {
		_, exists := d.VariableDefinitionByNameAndOperation(operationDefinition, b)
		return exists
	})
}
