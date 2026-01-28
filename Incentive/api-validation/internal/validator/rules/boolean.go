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

type booleanValidator struct {
	pass        *codegen.Pass
	field       *ast.Field
	expressions map[string]string
	structName  string
	ruleName    string
	parentPath  string
}

var _ validator.Validator = (*booleanValidator)(nil)

const booleanKey = "%s-boolean"

func (b *booleanValidator) Validate() string {
	fieldName := b.FieldName()
	// Use external helper function for boolean validation
	validation := fmt.Sprintf("!validationhelper.IsValidBoolean(t.%s)", fieldName)

	// Get the field type for omitempty handling
	typ := b.pass.TypesInfo.TypeOf(b.field.Type)

	// Wrap with omitempty logic if needed
	return wrapWithOmitEmpty(b.pass, typ, fieldName, validation, b.expressions)
}

func (b *booleanValidator) FieldName() string {
	return b.field.Names[0].Name
}

func (b *booleanValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(b.structName, b.parentPath, b.FieldName())
}

func (b *booleanValidator) Err() string {
	key := fmt.Sprintf(booleanKey, b.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the field is not a valid boolean string.
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must be a valid boolean (true, false, 1, 0, yes, no, on, off)", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sBooleanValidation", b.structName, b.FieldName())
	currentErrVarName := b.ErrVariable()

	replacer := strings.NewReplacer(
		"[@ERRVARIABLE]", currentErrVarName,
		"[@LEGACYERRVAR]", legacyErrVarName,
		"[@FIELD]", b.FieldName(),
		"[@PATH]", b.FieldPath().String(),
		"[@TYPE]", b.ruleName,
	)

	if currentErrVarName != legacyErrVarName {
		return replacer.Replace(deprecationNoticeTemplate + errTemplate)
	}

	return replacer.Replace(errTemplate)
}

func (b *booleanValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]BooleanValidation", `[@PATH]`, b.FieldPath().CleanedPath())
}

func (b *booleanValidator) Imports() []string {
	return []string{"gitlab.cept.gov.in/it-2.0-common/n-api-validation/validation/validationhelper"}
}

// ValidateBoolean creates a new booleanValidator for string types.
func ValidateBoolean(input registry.ValidatorInput) validator.Validator {
	typ := input.Pass.TypesInfo.TypeOf(input.Field.Type)

	// Check if it's a string type
	basic, ok := typ.Underlying().(*types.Basic)
	if !ok || basic.Kind() != types.String {
		return nil
	}

	return &booleanValidator{
		pass:        input.Pass,
		field:       input.Field,
		expressions: input.Expressions,
		structName:  input.StructName,
		ruleName:    input.RuleName,
		parentPath:  input.ParentPath,
	}
}
