package asyncapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jensneuse/graphql-go-tools/internal/pkg/unsafebytes"
	"github.com/jensneuse/graphql-go-tools/pkg/lexer/literal"
	"unicode"

	"github.com/asyncapi/parser-go/pkg/parser"
	"github.com/jensneuse/graphql-go-tools/pkg/ast"
	"github.com/jensneuse/graphql-go-tools/pkg/operationreport"
)

type SchemaObject struct {
	Type       string                            `json:"type"`
	Properties map[string]map[string]interface{} `json:"properties"`
}

// https://studio.asyncapi.com/#operation-publish-smartylighting.streetlights.1.0.event.{streetlightId}.lighting.measured
type OperationObject struct {
	OperationID string         `json:"operationId"`
	Summary     string         `json:"summary"`
	Description string         `json:"description"`
	Message     *MessageObject `json:"message"`
}

type MessageObject struct {
	Name        string      `json:"name"`
	Summary     string      `json:"summary"`
	Description string      `json:"description"`
	Payload     interface{} `json:"payload"`
}

type Channel struct {
	Ref         string           `json:"$ref"`
	Description string           `json:"description"`
	Publish     *OperationObject `json:"publish"`
}

type Schema struct {
	AsyncAPI string             `json:"asyncapi"`
	Channels map[string]Channel `json:"channels"`
}

func lowercaseFirstLetter(s string) string {
	copyStr := []rune(s)
	copyStr[0] = unicode.ToLower(copyStr[0])
	return string(copyStr)
}

func validateAndParseAsyncAPIDocument(input []byte) (*Schema, error) {
	r := bytes.NewBuffer(input)
	p, err := parser.New()
	if err != nil {
		return nil, err
	}

	w := bytes.NewBuffer(nil)
	err = p(r, w)
	if err != nil {
		return nil, err
	}

	var jsonData Schema
	err = json.NewDecoder(w).Decode(&jsonData)
	if err != nil {
		return nil, err
	}
	return &jsonData, nil
}

func (i *importer) findTypeByNameAndKind(name string, kind ast.TypeKind) int {
	for ref, ty := range i.doc.Types {
		if i.doc.Input.ByteSliceString(ty.Name) == name && ty.TypeKind == kind {
			return ref
		}
	}
	return ast.InvalidRef
}

func (i *importer) getOrCreateType(name string, kind ast.TypeKind) int {
	ref := i.findTypeByNameAndKind(name, kind)
	if ref != ast.InvalidRef {
		return ref
	}

	switch kind {
	case ast.TypeKindNamed:
		return i.doc.AddNamedType(unsafebytes.StringToBytes(name))
	default:
		// TODO: Handle other type kinds
		return ast.InvalidRef
	}
}

func asyncApiTypesToGQLTypes(t string) ([]byte, error) {
	switch t {
	case "string":
		return literal.STRING, nil
	case "integer":
		return literal.INT, nil
	default:
		return nil, fmt.Errorf("unknown type: %s", t)
	}
}

func (i *importer) processPayload(msg *MessageObject) error {
	var inputValueDefinitionRefs []int
	switch payload := msg.Payload.(type) {
	case map[string]interface{}:
		// TODO: Check type
		payloadType, ok := payload["type"]
		if !ok {
			return fmt.Errorf("missing keyword: type")
		}
		if payloadType == "object" {
			properties, ok := payload["properties"]
			if !ok {
				return fmt.Errorf("missing keyword: properties")
			}
			for key, rawValue := range properties.(map[string]interface{}) {
				value := rawValue.(map[string]interface{})
				gtype, err := asyncApiTypesToGQLTypes(value["type"].(string))
				if err != nil {
					return err
				}
				// TODO: How do you determine type kind?
				typeRef := i.getOrCreateType(string(gtype), ast.TypeKindNamed)
				inputValueDefinitionRef := i.doc.ImportInputValueDefinition(key, value["description"].(string), typeRef, ast.DefaultValue{})
				inputValueDefinitionRefs = append(inputValueDefinitionRefs, inputValueDefinitionRef)
			}
		}
	}
	i.doc.ImportInputObjectTypeDefinition(msg.Name, msg.Description, inputValueDefinitionRefs)
	return nil
}

func (i *importer) processPublishField(pb *OperationObject) error {
	var inputValueRefs []int
	if pb.Message != nil {
		if pb.Message.Payload != nil {
			err := i.processPayload(pb.Message)
			if err != nil {
				return err
			}
		}
		typeRef := i.doc.AddNamedType([]byte(pb.Message.Name))
		inputValueRef := i.doc.ImportInputValueDefinition(lowercaseFirstLetter(pb.Message.Name), "", typeRef, ast.DefaultValue{})
		inputValueRefs = append(inputValueRefs, inputValueRef)
		fieldRef := i.doc.ImportFieldDefinition(pb.OperationID, pb.Summary, typeRef, inputValueRefs, nil)
		i.fieldRefs = append(i.fieldRefs, fieldRef)
	}
	return nil
}

type importer struct {
	doc       *ast.Document
	fieldRefs []int
	schema    *Schema
}

func (i *importer) importAsyncAPIDocument() error {
	for _, obj := range i.schema.Channels {
		if obj.Publish != nil {
			err := i.processPublishField(obj.Publish)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ImportAsyncAPIDocumentByte(input []byte) (ast.Document, operationreport.Report) {
	doc := ast.NewDocument()
	report := operationreport.Report{}
	schema, err := validateAndParseAsyncAPIDocument(input)
	if err != nil {
		report.AddInternalError(err)
		return *doc, report
	}

	i := importer{
		doc:    doc,
		schema: schema,
	}

	doc.ImportSchemaDefinition("", "", "Subscription")

	err = i.importAsyncAPIDocument()
	if err != nil {
		report.AddInternalError(err)
		return *doc, report
	}

	doc.ImportObjectTypeDefinition("Subscription", "", i.fieldRefs, []int{})
	return *doc, report
}
