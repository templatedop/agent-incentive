// Package rules implements validation rules for fields in structs.
package rules

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/gostaticanalysis/codegen"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/markers"
	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/validator"
	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/validator/registry"
)

type required_ifValidator struct {
	pass          *codegen.Pass
	field         *ast.Field
	otherField    string
	expectedValue string
	structName    string
	ruleName      string
	parentPath    string
}

var _ validator.Validator = (*required_ifValidator)(nil)

const required_ifKey = "%s-required_if"

func (r *required_ifValidator) Validate() string {
	typ := r.pass.TypesInfo.TypeOf(r.field.Type)
	zero := getZeroValueForType(typ)
	fieldName := r.FieldName()

	// Generate: if otherField == expectedValue && thisField == zeroValue { fail }
	return fmt.Sprintf("t.%s == %s && t.%s == %s",
		r.otherField, r.expectedValue, fieldName, zero)
}

func (r *required_ifValidator) FieldName() string {
	return r.field.Names[0].Name
}

func (r *required_ifValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(r.structName, r.parentPath, r.FieldName())
}

func (r *required_ifValidator) Err() string {
	key := fmt.Sprintf(required_ifKey, r.structName+r.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the field is required due to another field's value.
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] is required when [@OTHER] equals [@VALUE]", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sRequiredIfValidation", r.structName, r.FieldName())
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

func (r *required_ifValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]RequiredIfValidation", "[@PATH]", r.FieldPath().CleanedPath())
}

func (r *required_ifValidator) Imports() []string {
	return []string{}
}

// ValidateRequiredIf creates a new required_ifValidator.
// Format: required_if=OtherField Value
func ValidateRequiredIf(input registry.ValidatorInput) validator.Validator {
	expr, ok := input.Expressions[markers.GoValidMarkerRequired_if]
	if !ok {
		return nil
	}

	// Parse "FieldName Value" format
	parts := strings.Fields(expr)
	if len(parts) < 2 {
		return nil
	}

	otherField := parts[0]
	expectedValue := strings.Join(parts[1:], " ")

	// Wrap value in quotes if it's not already quoted
	if !strings.HasPrefix(expectedValue, `"`) && !strings.HasPrefix(expectedValue, "`") {
		expectedValue = fmt.Sprintf(`"%s"`, expectedValue)
	}

	return &required_ifValidator{
		pass:          input.Pass,
		field:         input.Field,
		otherField:    otherField,
		expectedValue: expectedValue,
		structName:    input.StructName,
		ruleName:      input.RuleName,
		parentPath:    input.ParentPath,
	}
}

// Helper function to get zero value for a type
func getZeroValueForType(typ types.Type) string {
	switch t := typ.Underlying().(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.String:
			return `""`
		case types.Bool:
			return "false"
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
			types.Float32, types.Float64, types.Complex64, types.Complex128:
			return "0"
		}
	case *types.Pointer, *types.Slice, *types.Map, *types.Chan, *types.Interface:
		return "nil"
	case *types.Struct:
		return fmt.Sprintf("%s{}", typ.String())
	}
	return "nil"
}
