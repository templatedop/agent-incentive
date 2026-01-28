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

type required_unlessValidator struct {
	pass          *codegen.Pass
	field         *ast.Field
	otherField    string
	expectedValue string
	structName    string
	ruleName      string
	parentPath    string
}

var _ validator.Validator = (*required_unlessValidator)(nil)

const required_unlessKey = "%s-required_unless"

func (r *required_unlessValidator) Validate() string {
	typ := r.pass.TypesInfo.TypeOf(r.field.Type)
	zero := getZeroValueForType(typ)
	fieldName := r.FieldName()

	// Generate: if otherField != expectedValue && thisField == zeroValue { fail }
	return fmt.Sprintf("t.%s != %s && t.%s == %s",
		r.otherField, r.expectedValue, fieldName, zero)
}

func (r *required_unlessValidator) FieldName() string {
	return r.field.Names[0].Name
}

func (r *required_unlessValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(r.structName, r.parentPath, r.FieldName())
}

func (r *required_unlessValidator) Err() string {
	key := fmt.Sprintf(required_unlessKey, r.structName+r.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the field is required unless another field has a specific value.
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] is required unless [@OTHER] equals [@VALUE]", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sRequiredUnlessValidation", r.structName, r.FieldName())
	currentErrVarName := r.ErrVariable()

	// Escape quotes in the value for error message
	escapedValue := strings.ReplaceAll(r.expectedValue, `"`, `\"`)

	replacer := strings.NewReplacer(
		"[@ERRVARIABLE]", currentErrVarName,
		"[@LEGACYERRVAR]", legacyErrVarName,
		"[@FIELD]", r.FieldName(),
		"[@PATH]", r.FieldPath().String(),
		"[@OTHER]", r.otherField,
		"[@VALUE]", escapedValue,
		"[@TYPE]", r.ruleName,
	)

	if currentErrVarName != legacyErrVarName {
		return replacer.Replace(deprecationNoticeTemplate + errTemplate)
	}

	return replacer.Replace(errTemplate)
}

func (r *required_unlessValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]RequiredUnlessValidation", "[@PATH]", r.FieldPath().CleanedPath())
}

func (r *required_unlessValidator) Imports() []string {
	return []string{}
}

// ValidateRequiredUnless creates a new required_unlessValidator.
// Format: required_unless=OtherField Value
func ValidateRequiredUnless(input registry.ValidatorInput) validator.Validator {
	expr, ok := input.Expressions[markers.GoValidMarkerRequired_unless]
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

	return &required_unlessValidator{
		pass:          input.Pass,
		field:         input.Field,
		otherField:    otherField,
		expectedValue: expectedValue,
		structName:    input.StructName,
		ruleName:      input.RuleName,
		parentPath:    input.ParentPath,
	}
}
