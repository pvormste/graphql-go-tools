package graphql

import (
	"fmt"

	"github.com/jensneuse/graphql-go-tools/pkg/operationreport"
)

type RequestFieldsValidator interface {
	Validate(request *Request, schema *Schema, restrictions []Type) (RequestFieldsValidationResult, error)
}

type fieldsValidator struct {
}

func (d fieldsValidator) Validate(request *Request, schema *Schema, allowedFields []Type) (RequestFieldsValidationResult, error) {
	report := operationreport.Report{}
	if len(allowedFields) == 0 {
		return fieldsValidationResult(report, false, "_ALL_", "_ALL_")
	}

	requestedTypes := make(RequestTypes)
	NewExtractor().ExtractFieldsFromRequest(request, schema, &report, requestedTypes)

	allowedTypes := make(RequestTypes)
	for _, field := range allowedFields {
		fields := make(RequestFields)
		for _, fieldName := range field.Fields {
			fields[fieldName] = struct{}{}
		}
		allowedTypes[field.Name] = fields
	}

	for requestedTypeName, requestedFields := range requestedTypes {
		fieldsAllowance, isTypeAllowed := allowedTypes[requestedTypeName]

		if !isTypeAllowed {
			return fieldsValidationResult(report, false, requestedTypeName, "")
		}

		for fieldName, _ := range requestedFields {
			if _, isFieldAllowed := fieldsAllowance[fieldName]; !isFieldAllowed {
				return fieldsValidationResult(report, false, requestedTypeName, fieldName)
			}
		}
	}

	return fieldsValidationResult(report, true, "", "")
}

type RequestFieldsValidationResult struct {
	Valid  bool
	Errors Errors
}

func fieldsValidationResult(report operationreport.Report, valid bool, typeName, fieldName string) (RequestFieldsValidationResult, error) {
	result := RequestFieldsValidationResult{
		Valid:  valid,
		Errors: nil,
	}

	var errors OperationValidationErrors
	if !result.Valid {
		var msgStr string
		if fieldName != "" {
			msgStr = fmt.Sprintf("field: %s is not allowed on type: %s", fieldName, typeName)
		} else {
			msgStr = fmt.Sprintf("type: %s is not allowed", typeName)
		}

		errors = append(errors, OperationValidationError{
			Message: msgStr,
		})
	}
	result.Errors = errors

	if !report.HasErrors() {
		return result, nil
	}

	errors = append(errors, operationValidationErrorsFromOperationReport(report)...)
	result.Errors = errors

	var err error
	if len(report.InternalErrors) > 0 {
		err = report.InternalErrors[0]
	}

	return result, err
}
