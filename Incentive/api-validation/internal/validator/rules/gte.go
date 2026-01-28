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

type gteValidator struct {
	pass        *codegen.Pass
	field       *ast.Field
	gteValue    string
	expressions map[string]string
	structName  string
	ruleName    string
	parentPath  string
}

var _ validator.Validator = (*gteValidator)(nil)

const gteKey = "%s-gte"

func (m *gteValidator) Validate() string {
	fieldName := m.FieldName()
	validation := fmt.Sprintf("!(t.%s >= %s)", fieldName, m.gteValue)

	// Get the field type for omitempty handling
	typ := m.pass.TypesInfo.TypeOf(m.field.Type)

	// Wrap with omitempty logic if needed
	return wrapWithOmitEmpty(m.pass, typ, fieldName, validation, m.expressions)
}

func (m *gteValidator) FieldName() string {
	return m.field.Names[0].Name
}

func (m *gteValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(m.structName, m.parentPath, m.FieldName())
}

func (m *gteValidator) Err() string {
	key := fmt.Sprintf(gteKey, m.structName+m.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the value of the field is less than [@VALUE].
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must be greater than or equal to [@VALUE]", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sGTEValidation", m.structName, m.FieldName())
	currentErrVarName := m.ErrVariable()

	replacer := strings.NewReplacer(
		"[@ERRVARIABLE]", currentErrVarName,
		"[@LEGACYERRVAR]", legacyErrVarName,
		"[@FIELD]", m.FieldName(),
		"[@PATH]", m.FieldPath().String(),
		"[@VALUE]", m.gteValue,
		"[@TYPE]", m.ruleName,
	)

	if currentErrVarName != legacyErrVarName {
		return replacer.Replace(deprecationNoticeTemplate + errTemplate)
	}

	return replacer.Replace(errTemplate)
}

func (m *gteValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]GTEValidation", "[@PATH]", m.FieldPath().CleanedPath())
}

func (m *gteValidator) Imports() []string {
	return []string{}
}

// ValidateGTE creates a new gteValidator if the field type is numeric and the gte marker is present.
func ValidateGTE(input registry.ValidatorInput) validator.Validator {
	typ := input.Pass.TypesInfo.TypeOf(input.Field.Type)
	basic, ok := typ.Underlying().(*types.Basic)

	if !ok || (basic.Info()&types.IsNumeric) == 0 {
		return nil
	}

	gteValue, ok := input.Expressions[markers.GoValidMarkerGte]
	if !ok {
		return nil
	}

	return &gteValidator{
		pass:        input.Pass,
		field:       input.Field,
		gteValue:    gteValue,
		expressions: input.Expressions,
		structName:  input.StructName,
		ruleName:    input.RuleName,
		parentPath:  input.ParentPath,
	}
}
