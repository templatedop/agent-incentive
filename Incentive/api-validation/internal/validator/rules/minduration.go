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

type mindurationValidator struct {
	pass        *codegen.Pass
	field       *ast.Field
	minDuration string
	structName  string
	ruleName    string
	parentPath  string
}

var _ validator.Validator = (*mindurationValidator)(nil)

const mindurationKey = "%s-minduration"

func (m *mindurationValidator) Validate() string {
	fieldName := m.FieldName()
	return fmt.Sprintf("func() bool { d, _ := time.ParseDuration(%q); return t.%s < d }()", m.minDuration, fieldName)
}

func (m *mindurationValidator) FieldName() string {
	return m.field.Names[0].Name
}

func (m *mindurationValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(m.structName, m.parentPath, m.FieldName())
}

func (m *mindurationValidator) Err() string {
	key := fmt.Sprintf(mindurationKey, m.structName+m.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the duration is less than the minimum.
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must be at least [@VALUE]", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sMindurationValidation", m.structName, m.FieldName())
	currentErrVarName := m.ErrVariable()

	replacer := strings.NewReplacer(
		"[@ERRVARIABLE]", currentErrVarName,
		"[@LEGACYERRVAR]", legacyErrVarName,
		"[@FIELD]", m.FieldName(),
		"[@PATH]", m.FieldPath().String(),
		"[@VALUE]", m.minDuration,
		"[@TYPE]", m.ruleName,
	)

	if currentErrVarName != legacyErrVarName {
		return replacer.Replace(deprecationNoticeTemplate + errTemplate)
	}

	return replacer.Replace(errTemplate)
}

func (m *mindurationValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]MindurationValidation", "[@PATH]", m.FieldPath().CleanedPath())
}

func (m *mindurationValidator) Imports() []string {
	return []string{"time"}
}

// ValidateMinduration creates a new mindurationValidator for time.Duration types.
func ValidateMinduration(input registry.ValidatorInput) validator.Validator {
	typ := input.Pass.TypesInfo.TypeOf(input.Field.Type)

	// Check if it's time.Duration
	if typ.String() != "time.Duration" {
		return nil
	}

	minDuration, ok := input.Expressions[markers.GoValidMarkerMinduration]
	if !ok {
		return nil
	}

	return &mindurationValidator{
		pass:        input.Pass,
		field:       input.Field,
		minDuration: minDuration,
		structName:  input.StructName,
		ruleName:    input.RuleName,
		parentPath:  input.ParentPath,
	}
}
