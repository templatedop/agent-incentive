// Package rules implements validation rules for fields in structs.
package rules

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/gostaticanalysis/codegen"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/validator"
	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/validator/registry"
)

type longitudeValidator struct {
	pass        *codegen.Pass
	field       *ast.Field
	expressions map[string]string
	structName  string
	ruleName    string
	parentPath  string
}

var _ validator.Validator = (*longitudeValidator)(nil)

const longitudeKey = "%s-longitude"

func (v *longitudeValidator) Validate() string {
	fieldName := v.FieldName()
	validation := fmt.Sprintf("!validationhelper.IsValidLongitude(t.%s)", fieldName)

	// Get the field type for omitempty handling
	typ := v.pass.TypesInfo.TypeOf(v.field.Type)

	// Wrap with omitempty logic if needed
	return wrapWithOmitEmpty(v.pass, typ, fieldName, validation, v.expressions)
}

func (v *longitudeValidator) FieldName() string {
	return v.field.Names[0].Name
}

func (v *longitudeValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(v.structName, v.parentPath, v.FieldName())
}

func (v *longitudeValidator) Err() string {
	key := fmt.Sprintf(longitudeKey, v.structName+v.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the field is not a valid longitude (-180 to 180).
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must be a valid longitude (-180 to 180)", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sLongitudeValidation", v.structName, v.FieldName())
	currentErrVarName := v.ErrVariable()

	replacer := strings.NewReplacer(
		"[@ERRVARIABLE]", currentErrVarName,
		"[@LEGACYERRVAR]", legacyErrVarName,
		"[@FIELD]", v.FieldName(),
		"[@PATH]", v.FieldPath().String(),
		"[@TYPE]", v.ruleName,
	)

	if currentErrVarName != legacyErrVarName {
		return replacer.Replace(deprecationNoticeTemplate + errTemplate)
	}

	return replacer.Replace(errTemplate)
}

func (v *longitudeValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]LongitudeValidation", `[@PATH]`, v.FieldPath().CleanedPath())
}

func (v *longitudeValidator) Imports() []string {
	return []string{"gitlab.cept.gov.in/it-2.0-common/n-api-validation/validation/validationhelper"}
}

// ValidateLongitude creates a new longitudeValidator for string types.
func ValidateLongitude(input registry.ValidatorInput) validator.Validator {
	typ := input.Pass.TypesInfo.TypeOf(input.Field.Type)

	// Check if it's a string type
	basic, ok := typ.Underlying().(*types.Basic)
	if !ok || basic.Kind() != types.String {
		return nil
	}

	return &longitudeValidator{
		pass:        input.Pass,
		field:       input.Field,
		expressions: input.Expressions,
		structName:  input.StructName,
		ruleName:    input.RuleName,
		parentPath:  input.ParentPath,
	}
}
