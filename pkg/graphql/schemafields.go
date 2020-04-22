package graphql

import (
	"fmt"

	"github.com/jensneuse/graphql-go-tools/pkg/introspection"
	"github.com/jensneuse/graphql-go-tools/pkg/operationreport"
)

var DefaultSchemaFieldsGenerator SchemaFieldsGenerator = schemaFieldsGenerator{}

type (
	SchemaType struct {
		Name   string        `json:"name"`
		Type   string        `json:"type"`
		Fields []SchemaField `json:"fields"`
	}

	SchemaField struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		TypeRef string `json:"type_ref"`
	}

	SchemaFieldsGenerator interface {
		Generate(schema string) (SchemaFieldsResult, error)
	}
)

type schemaFieldsGenerator struct{}

func (g schemaFieldsGenerator) Generate(schema string) (SchemaFieldsResult, error) {
	if schema == "" {
		return SchemaFieldsResult{}, ErrEmptySchema
	}

	parsedSchema, err := NewSchemaFromString(schema)
	if err != nil {
		return SchemaFieldsResult{}, err
	}

	var (
		report operationreport.Report
		data   introspection.Data
	)

	generator := introspection.NewGenerator()
	generator.Generate(&parsedSchema.document, &report, &data)

	if report.HasErrors() {
		return schemaFieldsResult(nil, report)
	}

	types := g.extractTypes(&data)
	return schemaFieldsResult(types, report)
}

func (g schemaFieldsGenerator) extractTypes(data *introspection.Data) []SchemaType {
	var types []SchemaType

	q, m, s := "Query", "Mutation", "Subscription"

	// TODO add checks for different names of query, mutation, subscription

	var (
		query        *SchemaType
		mutation     *SchemaType
		subscription *SchemaType
	)

	var objectTypes []SchemaType
	for _, fullType := range data.Schema.Types {
		switch fullType.Kind {
		case introspection.INTERFACE, introspection.OBJECT:
			t := SchemaType{}
			t.Name = fullType.Name
			t.Type = fullType.Kind.String()
			for _, field := range fullType.Fields {
				f := SchemaField{}
				f.Name = field.Name
				f.Type, f.TypeRef = fieldType(field.Type)
				t.Fields = append(t.Fields, f)
			}

			switch t.Name {
			case q:
				query = &t
			case m:
				mutation = &t
			case s:
				subscription = &t
			default:
				objectTypes = append(objectTypes, t)
			}
		}
	}
	// place main objects first
	if query != nil {
		types = append(types, *query)
	}
	if mutation != nil {
		types = append(types, *mutation)
	}
	if subscription != nil {
		types = append(types, *subscription)
	}
	types = append(types, objectTypes...)

	return types
}

func fieldType(t introspection.TypeRef) (fullType, underlyingType string) {
	switch t.Kind {
	case introspection.NONNULL:
		f, u := fieldType(*t.OfType)
		return fmt.Sprintf("%s!", f), u
	case introspection.LIST:
		f, u := fieldType(*t.OfType)
		return fmt.Sprintf("[%s]", f), u
	case introspection.SCALAR, introspection.ENUM, introspection.UNION:
		return *t.Name, t.Kind.String()
	default: // only object as we do filtering
		return *t.Name, *t.Name
	}
}

type SchemaFieldsResult struct {
	Types  []SchemaType
	Errors Errors
}

func schemaFieldsResult(types []SchemaType, report operationreport.Report) (SchemaFieldsResult, error) {
	result := SchemaFieldsResult{
		Types:  types,
		Errors: nil,
	}

	if !report.HasErrors() {
		return result, nil
	}

	result.Errors = operationValidationErrorsFromOperationReport(report)

	var err error
	if len(report.InternalErrors) > 0 {
		err = report.InternalErrors[0]
	}

	return result, err
}
