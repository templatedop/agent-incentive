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

type iscolourValidator struct {
	pass        *codegen.Pass
	field       *ast.Field
	expressions map[string]string
	structName  string
	ruleName    string
	parentPath  string
}

var _ validator.Validator = (*iscolourValidator)(nil)

const iscolourKey = "%s-iscolour"

func (v *iscolourValidator) Validate() string {
	fieldName := v.FieldName()
	validation := fmt.Sprintf("!validationhelper.IsValidColour(t.%s)", fieldName)

	// Get the field type for omitempty handling
	typ := v.pass.TypesInfo.TypeOf(v.field.Type)

	// Wrap with omitempty logic if needed
	return wrapWithOmitEmpty(v.pass, typ, fieldName, validation, v.expressions)
}

func (v *iscolourValidator) FieldName() string {
	return v.field.Names[0].Name
}

func (v *iscolourValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(v.structName, v.parentPath, v.FieldName())
}

func (v *iscolourValidator) Err() string {
	key := fmt.Sprintf(iscolourKey, v.structName+v.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the field is not a valid color format.
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must be a valid color format", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sIscolourValidation", v.structName, v.FieldName())
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

func (v *iscolourValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]IscolourValidation", `[@PATH]`, v.FieldPath().CleanedPath())
}

func (v *iscolourValidator) Imports() []string {
	return []string{"gitlab.cept.gov.in/it-2.0-common/n-api-validation/validation/validationhelper"}
}

// ValidateIscolour creates a new iscolourValidator for string types.
func ValidateIscolour(input registry.ValidatorInput) validator.Validator {
	typ := input.Pass.TypesInfo.TypeOf(input.Field.Type)

	// Check if it's a string type
	basic, ok := typ.Underlying().(*types.Basic)
	if !ok || basic.Kind() != types.String {
		return nil
	}

	return &iscolourValidator{
		pass:        input.Pass,
		field:       input.Field,
		expressions: input.Expressions,
		structName:  input.StructName,
		ruleName:    input.RuleName,
		parentPath:  input.ParentPath,
	}
}
