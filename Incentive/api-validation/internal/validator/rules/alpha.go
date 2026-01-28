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

type alphaValidator struct {
	pass        *codegen.Pass
	field       *ast.Field
	expressions map[string]string
	structName  string
	ruleName    string
	parentPath  string
}

var _ validator.Validator = (*alphaValidator)(nil)

const alphaKey = "%s-alpha"

func (v *alphaValidator) Validate() string {
	// Use external helper function for better maintainability
	fieldName := v.FieldName()
	validation := fmt.Sprintf(`!validationhelper.IsValidAlpha(t.%s)`, fieldName)

	// Get the field type for omitempty handling
	typ := v.pass.TypesInfo.TypeOf(v.field.Type)

	// Wrap with omitempty logic if needed
	return wrapWithOmitEmpty(v.pass, typ, fieldName, validation, v.expressions)
}

func (v *alphaValidator) FieldName() string {
	return v.field.Names[0].Name
}

func (v *alphaValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(v.structName, v.parentPath, v.FieldName())
}

func (v *alphaValidator) Err() string {
	key := fmt.Sprintf(alphaKey, v.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when field [@FIELD] is not alphabetic.
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must be alphabetic", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sAlphaValidation", v.structName, v.FieldName())
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

func (v *alphaValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]AlphaValidation", "[@PATH]", v.FieldPath().CleanedPath())
}

func (v *alphaValidator) Imports() []string {
	// Import validation helper package
	return []string{"gitlab.cept.gov.in/it-2.0-common/n-api-validation/validation/validationhelper"}
}

// ValidateAlpha creates a new alphaValidator for string types.
func ValidateAlpha(input registry.ValidatorInput) validator.Validator {
	typ := input.Pass.TypesInfo.TypeOf(input.Field.Type)

	// Check if it's a string type
	basic, ok := typ.Underlying().(*types.Basic)

	if !ok || basic.Kind() != types.String {
		return nil
	}

	return &alphaValidator{
		pass:        input.Pass,
		field:       input.Field,
		expressions: input.Expressions,
		structName:  input.StructName,
		ruleName:    input.RuleName,
		parentPath:  input.ParentPath,
	}
}
