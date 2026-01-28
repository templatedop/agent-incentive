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

type excluded_withoutValidator struct {
	pass       *codegen.Pass
	field      *ast.Field
	fields     []string
	structName string
	ruleName   string
	parentPath string
}

var _ validator.Validator = (*excluded_withoutValidator)(nil)

const excluded_withoutKey = "%s-excluded_without"

func (e *excluded_withoutValidator) Validate() string {
	typ := e.pass.TypesInfo.TypeOf(e.field.Type)
	zero := getZeroValueForType(typ)
	fieldName := e.FieldName()

	var conditions []string
	for _, f := range e.fields {
		conditions = append(conditions, fmt.Sprintf(`t.%s == ""`, f))
	}

	return fmt.Sprintf("(%s) && t.%s != %s",
		strings.Join(conditions, " || "), fieldName, zero)
}

func (e *excluded_withoutValidator) FieldName() string {
	return e.field.Names[0].Name
}

func (e *excluded_withoutValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(e.structName, e.parentPath, e.FieldName())
}

func (e *excluded_withoutValidator) Err() string {
	key := fmt.Sprintf(excluded_withoutKey, e.structName+e.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the field must be absent because other fields are absent.
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must be absent when any of [@FIELDS] are absent", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sExcludedWithoutValidation", e.structName, e.FieldName())
	currentErrVarName := e.ErrVariable()

	replacer := strings.NewReplacer(
		"[@ERRVARIABLE]", currentErrVarName,
		"[@LEGACYERRVAR]", legacyErrVarName,
		"[@FIELD]", e.FieldName(),
		"[@PATH]", e.FieldPath().String(),
		"[@FIELDS]", strings.Join(e.fields, ", "),
		"[@TYPE]", e.ruleName,
	)

	if currentErrVarName != legacyErrVarName {
		return replacer.Replace(deprecationNoticeTemplate + errTemplate)
	}

	return replacer.Replace(errTemplate)
}

func (e *excluded_withoutValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]ExcludedWithoutValidation", "[@PATH]", e.FieldPath().CleanedPath())
}

func (e *excluded_withoutValidator) Imports() []string {
	return []string{}
}

// ValidateExcludedWithout creates a new excluded_withoutValidator.
func ValidateExcludedWithout(input registry.ValidatorInput) validator.Validator {
	expr, ok := input.Expressions[markers.GoValidMarkerExcluded_without]
	if !ok {
		return nil
	}

	fields := strings.Fields(expr)
	if len(fields) == 0 {
		return nil
	}

	return &excluded_withoutValidator{
		pass:       input.Pass,
		field:      input.Field,
		fields:     fields,
		structName: input.StructName,
		ruleName:   input.RuleName,
		parentPath: input.ParentPath,
	}
}
