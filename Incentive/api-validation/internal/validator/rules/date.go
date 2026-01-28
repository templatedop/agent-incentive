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

// dateValidator validates a date string with format dd/mm/yy.
type dateValidator struct {
	pass        *codegen.Pass
	field       *ast.Field
	expressions map[string]string
	structName  string
	parentPath  string
	ruleName    string
}

var _ validator.Validator = (*dateValidator)(nil)

const dateKey = "%s-date"

func (d *dateValidator) Validate() string {
	fieldName := d.FieldName()
	// use helper for pattern check
	validation := fmt.Sprintf("!validationhelper.IsValidDateDDMMYY(t.%s)", fieldName)

	// Get the field type for omitempty handling
	typ := d.pass.TypesInfo.TypeOf(d.field.Type)

	// Wrap with omitempty logic if needed
	return wrapWithOmitEmpty(d.pass, typ, fieldName, validation, d.expressions)
}

func (d *dateValidator) FieldName() string { return d.field.Names[0].Name }

func (d *dateValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(d.structName, d.parentPath, d.FieldName())
}

func (d *dateValidator) Err() string {
	key := fmt.Sprintf(dateKey, d.FieldPath().CleanedPath())
	if validator.GeneratorMemory[key] {
		return ""
	}
	validator.GeneratorMemory[key] = true

	const errTemplate = `
		// [@ERRVARIABLE] is the error returned when the field is not a valid date (dd/mm/yy).
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must be a valid date (dd/mm/yy)", Path: "[@PATH]", Type: "[@TYPE]"}
	`
	replacer := strings.NewReplacer(
		"[@ERRVARIABLE]", d.ErrVariable(),
		"[@FIELD]", d.FieldName(),
		"[@PATH]", d.FieldPath().String(),
		"[@TYPE]", d.ruleName,
	)
	return replacer.Replace(errTemplate)
}

func (d *dateValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]DateValidation", "[@PATH]", d.FieldPath().CleanedPath())
}

func (d *dateValidator) Imports() []string {
	return []string{"gitlab.cept.gov.in/it-2.0-common/n-api-validation/validation/validationhelper"}
}

// ValidateDate constructs a dateValidator if the field is a string.
func ValidateDate(input registry.ValidatorInput) validator.Validator {
	typ := input.Pass.TypesInfo.TypeOf(input.Field.Type)
	basic, ok := typ.Underlying().(*types.Basic)
	if !ok || basic.Kind() != types.String {
		return nil
	}
	return &dateValidator{pass: input.Pass, field: input.Field, expressions: input.Expressions, structName: input.StructName, parentPath: input.ParentPath, ruleName: input.RuleName}
}
