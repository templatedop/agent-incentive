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

type ltValidator struct {
	pass        *codegen.Pass
	field       *ast.Field
	ltValue     string
	expressions map[string]string
	structName  string
	ruleName    string
	parentPath  string
}

var _ validator.Validator = (*ltValidator)(nil)

const ltKey = "%s-lt"

func (m *ltValidator) Validate() string {
	fieldName := m.FieldName()
	validation := fmt.Sprintf("!(t.%s < %s)", fieldName, m.ltValue)

	// Get the field type for omitempty handling
	typ := m.pass.TypesInfo.TypeOf(m.field.Type)

	// Wrap with omitempty logic if needed
	return wrapWithOmitEmpty(m.pass, typ, fieldName, validation, m.expressions)
}

func (m *ltValidator) FieldName() string {
	return m.field.Names[0].Name
}

func (m *ltValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(m.structName, m.parentPath, m.FieldName())
}

func (m *ltValidator) Err() string {
	key := fmt.Sprintf(ltKey, m.structName+m.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the value of the field is greater than the [@VALUE].
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must be less than [@VALUE]", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sLTValidation", m.structName, m.FieldName())
	currentErrVarName := m.ErrVariable()

	replacer := strings.NewReplacer(
		"[@ERRVARIABLE]", currentErrVarName,
		"[@LEGACYERRVAR]", legacyErrVarName,
		"[@FIELD]", m.FieldName(),
		"[@PATH]", m.FieldPath().String(),
		"[@VALUE]", m.ltValue,
		"[@TYPE]", m.ruleName,
	)

	if currentErrVarName != legacyErrVarName {
		return replacer.Replace(deprecationNoticeTemplate + errTemplate)
	}

	return replacer.Replace(errTemplate)
}

func (m *ltValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]LTValidation", "[@PATH]", m.FieldPath().CleanedPath())
}

func (m *ltValidator) Imports() []string {
	return []string{}
}

// ValidateLT creates a new ltValidator if the field type is numeric and the min marker is present.
func ValidateLT(input registry.ValidatorInput) validator.Validator {
	typ := input.Pass.TypesInfo.TypeOf(input.Field.Type)
	basic, ok := typ.Underlying().(*types.Basic)

	if !ok || (basic.Info()&types.IsNumeric) == 0 {
		return nil
	}

	ltValue, ok := input.Expressions[markers.GoValidMarkerLt]
	if !ok {
		return nil
	}

	return &ltValidator{
		pass:        input.Pass,
		field:       input.Field,
		ltValue:     ltValue,
		expressions: input.Expressions,
		structName:  input.StructName,
		ruleName:    input.RuleName,
		parentPath:  input.ParentPath,
	}
}
