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

type lowercaseValidator struct {
	pass        *codegen.Pass
	field       *ast.Field
	expressions map[string]string
	structName  string
	ruleName    string
	parentPath  string
}

var _ validator.Validator = (*lowercaseValidator)(nil)

const lowercaseKey = "%s-lowercase"

func (l *lowercaseValidator) Validate() string {
	fieldName := l.FieldName()
	validation := fmt.Sprintf("!validationhelper.IsLowercase(t.%s)", fieldName)

	// Get the field type for omitempty handling
	typ := l.pass.TypesInfo.TypeOf(l.field.Type)

	// Wrap with omitempty logic if needed
	return wrapWithOmitEmpty(l.pass, typ, fieldName, validation, l.expressions)
}

func (l *lowercaseValidator) FieldName() string {
	return l.field.Names[0].Name
}

func (l *lowercaseValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(l.structName, l.parentPath, l.FieldName())
}

func (l *lowercaseValidator) Err() string {
	key := fmt.Sprintf(lowercaseKey, l.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the field is not all lowercase.
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must be lowercase", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sLowercaseValidation", l.structName, l.FieldName())
	currentErrVarName := l.ErrVariable()

	replacer := strings.NewReplacer(
		"[@ERRVARIABLE]", currentErrVarName,
		"[@LEGACYERRVAR]", legacyErrVarName,
		"[@FIELD]", l.FieldName(),
		"[@PATH]", l.FieldPath().String(),
		"[@TYPE]", l.ruleName,
	)

	if currentErrVarName != legacyErrVarName {
		return replacer.Replace(deprecationNoticeTemplate + errTemplate)
	}

	return replacer.Replace(errTemplate)
}

func (l *lowercaseValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]LowercaseValidation", `[@PATH]`, l.FieldPath().CleanedPath())
}

func (l *lowercaseValidator) Imports() []string {
	return []string{"gitlab.cept.gov.in/it-2.0-common/n-api-validation/validation/validationhelper"}
}

// ValidateLowercase creates a new lowercaseValidator for string types.
func ValidateLowercase(input registry.ValidatorInput) validator.Validator {
	typ := input.Pass.TypesInfo.TypeOf(input.Field.Type)

	// Check if it's a string type
	basic, ok := typ.Underlying().(*types.Basic)
	if !ok || basic.Kind() != types.String {
		return nil
	}

	return &lowercaseValidator{
		pass:        input.Pass,
		field:       input.Field,
		expressions: input.Expressions,
		structName:  input.StructName,
		ruleName:    input.RuleName,
		parentPath:  input.ParentPath,
	}
}
