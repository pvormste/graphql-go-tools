package astnormalization

import (
	"bytes"

	"github.com/tidwall/sjson"

	"github.com/pvormste/graphql-go-tools/internal/pkg/unsafebytes"
	"github.com/pvormste/graphql-go-tools/pkg/ast"
	"github.com/pvormste/graphql-go-tools/pkg/astimport"
	"github.com/pvormste/graphql-go-tools/pkg/astvisitor"
)

func extractVariables(walker *astvisitor.Walker) *variablesExtractionVisitor {
	visitor := &variablesExtractionVisitor{
		Walker: walker,
	}
	walker.RegisterEnterDocumentVisitor(visitor)
	walker.RegisterEnterArgumentVisitor(visitor)
	walker.RegisterEnterOperationVisitor(visitor)
	return visitor
}

type variablesExtractionVisitor struct {
	*astvisitor.Walker
	operation, definition *ast.Document
	importer              astimport.Importer
	operationName         []byte
	skip                  bool
}

func (v *variablesExtractionVisitor) EnterOperationDefinition(ref int) {
	if len(v.operationName) == 0 {
		v.skip = false
		return
	}
	operationName := v.operation.OperationDefinitionNameBytes(ref)
	v.skip = !bytes.Equal(operationName, v.operationName)
}

func (v *variablesExtractionVisitor) EnterArgument(ref int) {
	if v.skip {
		return
	}
	if v.operation.Arguments[ref].Value.Kind == ast.ValueKindVariable {
		return
	}
	if len(v.Ancestors) == 0 || v.Ancestors[0].Kind != ast.NodeKindOperationDefinition {
		return
	}

	for i := range v.Ancestors {
		if v.Ancestors[i].Kind == ast.NodeKindDirective {
			return // skip all directives in any case
		}
	}

	inputValueDefinition, ok := v.Walker.ArgumentInputValueDefinition(ref)
	if !ok {
		return
	}

	containsVariable := v.operation.ValueContainsVariable(v.operation.Arguments[ref].Value)
	if containsVariable {
		v.traverseValue(v.operation.Arguments[ref].Value, ref, inputValueDefinition)
		return
	}

	variableNameBytes := v.operation.GenerateUnusedVariableDefinitionName(v.Ancestors[0].Ref)
	valueBytes, err := v.operation.ValueToJSON(v.operation.Arguments[ref].Value)
	if err != nil {
		return
	}
	v.operation.Input.Variables, err = sjson.SetRawBytes(v.operation.Input.Variables, unsafebytes.BytesToString(variableNameBytes), valueBytes)
	if err != nil {
		v.StopWithInternalErr(err)
		return
	}

	variable := ast.VariableValue{
		Name: v.operation.Input.AppendInputBytes(variableNameBytes),
	}

	v.operation.VariableValues = append(v.operation.VariableValues, variable)

	varRef := len(v.operation.VariableValues) - 1

	v.operation.Arguments[ref].Value.Ref = varRef
	v.operation.Arguments[ref].Value.Kind = ast.ValueKindVariable

	defRef, ok := v.ArgumentInputValueDefinition(ref)
	if !ok {
		return
	}

	defType := v.definition.InputValueDefinitions[defRef].Type

	importedDefType := v.importer.ImportType(defType, v.definition, v.operation)

	v.operation.VariableDefinitions = append(v.operation.VariableDefinitions, ast.VariableDefinition{
		VariableValue: ast.Value{
			Kind: ast.ValueKindVariable,
			Ref:  varRef,
		},
		Type: importedDefType,
	})

	newVariableRef := len(v.operation.VariableDefinitions) - 1

	v.operation.OperationDefinitions[v.Ancestors[0].Ref].VariableDefinitions.Refs =
		append(v.operation.OperationDefinitions[v.Ancestors[0].Ref].VariableDefinitions.Refs, newVariableRef)
	v.operation.OperationDefinitions[v.Ancestors[0].Ref].HasVariableDefinitions = true
}

func (v *variablesExtractionVisitor) EnterDocument(operation, definition *ast.Document) {
	v.operation, v.definition = operation, definition
}

func (v *variablesExtractionVisitor) traverseValue(value ast.Value, argRef, inputValueDefinition int) {
	switch value.Kind {
	case ast.ValueKindList:
		for _, ref := range v.operation.ListValues[value.Ref].Refs {
			listValue := v.operation.Value(ref)
			v.traverseValue(listValue, argRef, inputValueDefinition)
		}
	case ast.ValueKindObject:
		objectValueRefs := make([]int, len(v.operation.ObjectValues[value.Ref].Refs))
		copy(objectValueRefs, v.operation.ObjectValues[value.Ref].Refs)
		for _, ref := range objectValueRefs {
			fieldName := v.operation.Input.ByteSlice(v.operation.ObjectFields[ref].Name)
			fieldValue := v.operation.ObjectFields[ref].Value
			switch fieldValue.Kind {
			case ast.ValueKindVariable:
				continue
			default:

				typeName := v.definition.ResolveTypeNameString(v.definition.InputValueDefinitions[inputValueDefinition].Type)
				typeDefinitionNode, ok := v.definition.Index.FirstNodeByNameStr(typeName)
				if !ok {
					continue
				}
				objectFieldDefinition, ok := v.definition.NodeInputFieldDefinitionByName(typeDefinitionNode, fieldName)
				if !ok {
					continue
				}

				if v.operation.ValueContainsVariable(fieldValue) {
					v.traverseValue(fieldValue, argRef, objectFieldDefinition)
					continue
				}
				v.extractObjectValue(ref, fieldValue, objectFieldDefinition)
			}
		}
	}
}

func (v *variablesExtractionVisitor) extractObjectValue(objectField int, fieldValue ast.Value, inputValueDefinition int) {

	variableNameBytes := v.operation.GenerateUnusedVariableDefinitionName(v.Ancestors[0].Ref)
	valueBytes, err := v.operation.ValueToJSON(fieldValue)
	if err != nil {
		return
	}
	v.operation.Input.Variables, err = sjson.SetRawBytes(v.operation.Input.Variables, unsafebytes.BytesToString(variableNameBytes), valueBytes)
	if err != nil {
		v.StopWithInternalErr(err)
		return
	}

	variable := ast.VariableValue{
		Name: v.operation.Input.AppendInputBytes(variableNameBytes),
	}

	v.operation.VariableValues = append(v.operation.VariableValues, variable)

	varRef := len(v.operation.VariableValues) - 1

	v.operation.ObjectFields[objectField].Value.Kind = ast.ValueKindVariable
	v.operation.ObjectFields[objectField].Value.Ref = varRef

	defType := v.definition.InputValueDefinitions[inputValueDefinition].Type

	importedDefType := v.importer.ImportType(defType, v.definition, v.operation)

	v.operation.VariableDefinitions = append(v.operation.VariableDefinitions, ast.VariableDefinition{
		VariableValue: ast.Value{
			Kind: ast.ValueKindVariable,
			Ref:  varRef,
		},
		Type: importedDefType,
	})

	newVariableRef := len(v.operation.VariableDefinitions) - 1

	v.operation.OperationDefinitions[v.Ancestors[0].Ref].VariableDefinitions.Refs =
		append(v.operation.OperationDefinitions[v.Ancestors[0].Ref].VariableDefinitions.Refs, newVariableRef)
	v.operation.OperationDefinitions[v.Ancestors[0].Ref].HasVariableDefinitions = true
}
