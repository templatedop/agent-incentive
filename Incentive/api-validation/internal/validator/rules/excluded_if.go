// Package rules implements validation rules for fields in structs.
package rules

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/gostaticanalysis/codegen"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/markers"
	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/validator"
	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/validator/registry"
)

type excluded_ifValidator struct {
	pass          *codegen.Pass
	field         *ast.Field
	otherField    string
	expectedValue string
	structName    string
	ruleName      string
	parentPath    string
}

var _ validator.Validator = (*excluded_ifValidator)(nil)

const excluded_ifKey = "%s-excluded_if"

func (e *excluded_ifValidator) Validate() string {
	typ := e.pass.TypesInfo.TypeOf(e.field.Type)
	zero := getZeroValueForType(typ)
	fieldName := e.FieldName()

	// Generate: if otherField == expectedValue && thisField != zeroValue { fail }
	return fmt.Sprintf("t.%s == %s && t.%s != %s",
		e.otherField, e.expectedValue, fieldName, zero)
}

func (e *excluded_ifValidator) FieldName() string {
	return e.field.Names[0].Name
}

func (e *excluded_ifValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(e.structName, e.parentPath, e.FieldName())
}

func (e *excluded_ifValidator) Err() string {
	key := fmt.Sprintf(excluded_ifKey, e.structName+e.FieldPath().CleanedPath())

	if validator.GeneratorMemory[key] {
		return ""
	}

	validator.GeneratorMemory[key] = true

	const deprecationNoticeTemplate = `
		// Deprecated: Use [@ERRVARIABLE]
		//
		// [@LEGACYERRVAR] is deprecated and is kept for compatibility purpose.
		[@LEGACYERRVAR] = [@ERRVARIABLE]
	`

	const errTemplate = `
		// [@ERRVARIABLE] is the error returned when the field must be absent due to another field's value.
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must be absent when [@OTHER] equals [@VALUE]", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sExcludedIfValidation", e.structName, e.FieldName())
	currentErrVarName := e.ErrVariable()

	// Escape quotes in the value for error message
	escapedValue := strings.ReplaceAll(e.expectedValue, `"`, `\"`)

	replacer := strings.NewReplacer(
		"[@ERRVARIABLE]", currentErrVarName,
		"[@LEGACYERRVAR]", legacyErrVarName,
		"[@FIELD]", e.FieldName(),
		"[@PATH]", e.FieldPath().String(),
		"[@OTHER]", e.otherField,
		"[@VALUE]", escapedValue,
		"[@TYPE]", e.ruleName,
	)

	if currentErrVarName != legacyErrVarName {
		return replacer.Replace(deprecationNoticeTemplate + errTemplate)
	}

	return replacer.Replace(errTemplate)
}

func (e *excluded_ifValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]ExcludedIfValidation", "[@PATH]", e.FieldPath().CleanedPath())
}

func (e *excluded_ifValidator) Imports() []string {
	return []string{}
}

// ValidateExcludedIf creates a new excluded_ifValidator.
// Format: excluded_if=OtherField Value
func ValidateExcludedIf(input registry.ValidatorInput) validator.Validator {
	expr, ok := input.Expressions[markers.GoValidMarkerExcluded_if]
	if !ok {
		return nil
	}

	parts := strings.Fields(expr)
	if len(parts) < 2 {
		return nil
	}

	otherField := parts[0]
	expectedValue := strings.Join(parts[1:], " ")

	if !strings.HasPrefix(expectedValue, `"`) && !strings.HasPrefix(expectedValue, "`") {
		expectedValue = fmt.Sprintf(`"%s"`, expectedValue)
	}

	return &excluded_ifValidator{
		pass:          input.Pass,
		field:         input.Field,
		otherField:    otherField,
		expectedValue: expectedValue,
		structName:    input.StructName,
		ruleName:      input.RuleName,
		parentPath:    input.ParentPath,
	}
}
