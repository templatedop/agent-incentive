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

type minLengthValidator struct {
	pass           *codegen.Pass
	field          *ast.Field
	minLengthValue string
	expressions    map[string]string
	structName     string
	ruleName       string
	parentPath     string
}

var _ validator.Validator = (*minLengthValidator)(nil)

const minLengthKey = "%s-minlength"

func (m *minLengthValidator) Validate() string {
	fieldName := m.FieldName()
	validation := fmt.Sprintf("utf8.RuneCountInString(t.%s) < %s", fieldName, m.minLengthValue)

	// Get the field type for omitempty handling
	typ := m.pass.TypesInfo.TypeOf(m.field.Type)

	// Wrap with omitempty logic if needed
	return wrapWithOmitEmpty(m.pass, typ, fieldName, validation, m.expressions)
}

func (m *minLengthValidator) FieldName() string {
	return m.field.Names[0].Name
}

func (m *minLengthValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(m.structName, m.parentPath, m.FieldName())
}

func (m *minLengthValidator) Err() string {
	key := fmt.Sprintf(minLengthKey, m.structName+m.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the length of the field is less than the minimum of [@VALUE].
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must have a minimum length of [@VALUE]", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sMinLengthValidation", m.structName, m.FieldName())
	currentErrVarName := m.ErrVariable()

	replacer := strings.NewReplacer(
		"[@ERRVARIABLE]", currentErrVarName,
		"[@LEGACYERRVAR]", legacyErrVarName,
		"[@FIELD]", m.FieldName(),
		"[@PATH]", m.FieldPath().String(),
		"[@VALUE]", m.minLengthValue,
		"[@TYPE]", m.ruleName,
	)

	if currentErrVarName != legacyErrVarName {
		return replacer.Replace(deprecationNoticeTemplate + errTemplate)
	}

	return replacer.Replace(errTemplate)
}

func (m *minLengthValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]MinLengthValidation", "[@PATH]", m.FieldPath().CleanedPath())
}

func (m *minLengthValidator) Imports() []string {
	return []string{"unicode/utf8"}
}

// ValidateMinLength creates a new minLengthValidator if the field type is string and the minlength marker is present.
func ValidateMinLength(input registry.ValidatorInput) validator.Validator {
	typ := input.Pass.TypesInfo.TypeOf(input.Field.Type)
	basic, ok := typ.Underlying().(*types.Basic)

	if !ok || basic.Kind() != types.String {
		return nil
	}

	minLengthValue, ok := input.Expressions[markers.GoValidMarkerMinlength]
	if !ok {
		return nil
	}

	return &minLengthValidator{
		pass:           input.Pass,
		field:          input.Field,
		minLengthValue: minLengthValue,
		expressions:    input.Expressions,
		structName:     input.StructName,
		ruleName:       input.RuleName,
		parentPath:     input.ParentPath,
	}
}
