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

type maxdurationValidator struct {
	pass        *codegen.Pass
	field       *ast.Field
	maxDuration string
	structName  string
	ruleName    string
	parentPath  string
}

var _ validator.Validator = (*maxdurationValidator)(nil)

const maxdurationKey = "%s-maxduration"

func (m *maxdurationValidator) Validate() string {
	fieldName := m.FieldName()
	return fmt.Sprintf("func() bool { d, _ := time.ParseDuration(%q); return t.%s > d }()", m.maxDuration, fieldName)
}

func (m *maxdurationValidator) FieldName() string {
	return m.field.Names[0].Name
}

func (m *maxdurationValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(m.structName, m.parentPath, m.FieldName())
}

func (m *maxdurationValidator) Err() string {
	key := fmt.Sprintf(maxdurationKey, m.structName+m.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the duration exceeds the maximum.
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must not exceed [@VALUE]", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sMaxdurationValidation", m.structName, m.FieldName())
	currentErrVarName := m.ErrVariable()

	replacer := strings.NewReplacer(
		"[@ERRVARIABLE]", currentErrVarName,
		"[@LEGACYERRVAR]", legacyErrVarName,
		"[@FIELD]", m.FieldName(),
		"[@PATH]", m.FieldPath().String(),
		"[@VALUE]", m.maxDuration,
		"[@TYPE]", m.ruleName,
	)

	if currentErrVarName != legacyErrVarName {
		return replacer.Replace(deprecationNoticeTemplate + errTemplate)
	}

	return replacer.Replace(errTemplate)
}

func (m *maxdurationValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]MaxdurationValidation", "[@PATH]", m.FieldPath().CleanedPath())
}

func (m *maxdurationValidator) Imports() []string {
	return []string{"time"}
}

// ValidateMaxduration creates a new maxdurationValidator for time.Duration types.
func ValidateMaxduration(input registry.ValidatorInput) validator.Validator {
	typ := input.Pass.TypesInfo.TypeOf(input.Field.Type)

	// Check if it's time.Duration
	if typ.String() != "time.Duration" {
		return nil
	}

	maxDuration, ok := input.Expressions[markers.GoValidMarkerMaxduration]
	if !ok {
		return nil
	}

	return &maxdurationValidator{
		pass:        input.Pass,
		field:       input.Field,
		maxDuration: maxDuration,
		structName:  input.StructName,
		ruleName:    input.RuleName,
		parentPath:  input.ParentPath,
	}
}
