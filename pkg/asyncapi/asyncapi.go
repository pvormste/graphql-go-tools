package asyncapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"unicode"

	"github.com/asyncapi/parser-go/pkg/parser"
	"github.com/jensneuse/graphql-go-tools/internal/pkg/unsafebytes"
	"github.com/jensneuse/graphql-go-tools/pkg/ast"
	"github.com/jensneuse/graphql-go-tools/pkg/lexer/literal"
	"github.com/jensneuse/graphql-go-tools/pkg/operationreport"
)

type importer struct {
	doc       *ast.Document
	fieldRefs []int
	schema    *Schema
}

type Property struct {
	Type        string `json:"type"`
	Format      string `json:"format"`
	Minimum     string `json:"minimum"`
	Maximum     string `json:"maximum"`
	Description string `json:"description"`
}

type SchemaObject struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
}

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

type ParameterObject struct {
	Description string       `json:"description"`
	Schema      SchemaObject `json:"schema"`
}

type Channel struct {
	Description string                      `json:"description"`
	Publish     *OperationObject            `json:"publish"`
	Parameters  map[string]*ParameterObject `json:"parameters"`
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

func asyncApiTypeToGQLType(asyncApiType string) ([]byte, error) {
	switch asyncApiType {
	case "string":
		return literal.STRING, nil
	case "integer":
		return literal.INT, nil
	case "number":
		return literal.FLOAT, nil
	case "boolean":
		return literal.BOOLEAN, nil
	default:
		return nil, fmt.Errorf("unknown type: %s", asyncApiType)
	}
}

func (i *importer) extractSchemaObject(payload map[string]interface{}) (*SchemaObject, error) {
	schemaObject := &SchemaObject{
		Type:       "object",
		Properties: make(map[string]Property),
	}

	properties, ok := payload["properties"]
	if !ok {
		return nil, fmt.Errorf("missing keyword: properties")
	}

	for key := range properties.(map[string]interface{}) {
		// TODO: Fix this
		value := properties.(map[string]interface{})[key].(map[string]interface{})
		schemaObject.Properties[key] = Property{
			Type:        value["type"].(string),
			Description: value["description"].(string),
		}
	}
	return schemaObject, nil
}

func (i *importer) processMessage(msg *MessageObject) error {
	var objectFieldRefs []int
	switch payload := msg.Payload.(type) {
	case map[string]interface{}:
		// TODO: Check type
		payloadType, ok := payload["type"]
		if !ok {
			return fmt.Errorf("missing keyword: type")
		}
		// TODO: Payload type can be any type. Currently, we only handle schema objects.
		if payloadType == "object" {
			schemaObject, err := i.extractSchemaObject(payload)
			if err != nil {
				return err
			}

			// Sort the property names to produce a deterministic result. Good for tests.
			var keys []string
			for key := range schemaObject.Properties {
				keys = append(keys, key)
			}
			sort.Strings(keys)

			for _, key := range keys {
				property := schemaObject.Properties[key]
				gtype, err := asyncApiTypeToGQLType(property.Type)
				if err != nil {
					return err
				}
				// TODO: How do you determine type kind?
				typeRef := i.getOrCreateType(string(gtype), ast.TypeKindNamed)
				objectFieldRefs = append(objectFieldRefs, i.doc.ImportFieldDefinition(key, property.Description, typeRef, nil, nil))
			}
		}
	}
	if _, ok := i.doc.NodeByNameStr(msg.Name); !ok {
		i.doc.ImportObjectTypeDefinition(msg.Name, msg.Description, objectFieldRefs, nil)
	}
	return nil
}

func (i *importer) processParameters(params map[string]*ParameterObject) ([]int, error) {
	var inputValueDefRefs []int
	for param, properties := range params {
		gtype, err := asyncApiTypeToGQLType(properties.Schema.Type)
		if err != nil {
			return nil, err
		}
		// TODO: How do you determine type kind?
		argTypeRef := i.getOrCreateType(string(gtype), ast.TypeKindNamed)
		inputValueDefRef := i.doc.ImportInputValueDefinition(lowercaseFirstLetter(param), "", argTypeRef, ast.DefaultValue{})
		inputValueDefRefs = append(inputValueDefRefs, inputValueDefRef)
	}
	return inputValueDefRefs, nil
}

func (i *importer) processPublishField(channel Channel) error {
	if channel.Publish.Message == nil {
		return nil
	}

	pb := channel.Publish
	if pb.Message.Payload != nil {
		err := i.processMessage(pb.Message)
		if err != nil {
			return err
		}
	}

	typeRef := i.doc.AddNonNullNamedType([]byte(pb.Message.Name))
	inputValueDefRefs, err := i.processParameters(channel.Parameters)
	if err != nil {
		return err
	}

	fieldRef := i.doc.ImportFieldDefinition(pb.OperationID, pb.Summary, typeRef, inputValueDefRefs, nil)
	i.fieldRefs = append(i.fieldRefs, fieldRef)
	return nil
}

func (i *importer) importAsyncAPIDocument() error {
	// TODO: Extract server configuration from the AsyncAPI document
	for _, channel := range i.schema.Channels {
		if channel.Publish != nil {
			err := i.processPublishField(channel)
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
