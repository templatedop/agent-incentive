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

type excluded_without_allValidator struct {
	pass       *codegen.Pass
	field      *ast.Field
	fields     []string
	structName string
	ruleName   string
	parentPath string
}

var _ validator.Validator = (*excluded_without_allValidator)(nil)

const excluded_without_allKey = "%s-excluded_without_all"

func (e *excluded_without_allValidator) Validate() string {
	typ := e.pass.TypesInfo.TypeOf(e.field.Type)
	zero := getZeroValueForType(typ)
	fieldName := e.FieldName()

	var conditions []string
	for _, f := range e.fields {
		conditions = append(conditions, fmt.Sprintf(`t.%s == ""`, f))
	}

	return fmt.Sprintf("(%s) && t.%s != %s",
		strings.Join(conditions, " && "), fieldName, zero)
}

func (e *excluded_without_allValidator) FieldName() string {
	return e.field.Names[0].Name
}

func (e *excluded_without_allValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(e.structName, e.parentPath, e.FieldName())
}

func (e *excluded_without_allValidator) Err() string {
	key := fmt.Sprintf(excluded_without_allKey, e.structName+e.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the field must be absent because all other fields are absent.
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must be absent when all of [@FIELDS] are absent", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sExcludedWithoutAllValidation", e.structName, e.FieldName())
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

func (e *excluded_without_allValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]ExcludedWithoutAllValidation", "[@PATH]", e.FieldPath().CleanedPath())
}

func (e *excluded_without_allValidator) Imports() []string {
	return []string{}
}

// ValidateExcludedWithoutAll creates a new excluded_without_allValidator.
func ValidateExcludedWithoutAll(input registry.ValidatorInput) validator.Validator {
	expr, ok := input.Expressions[markers.GoValidMarkerExcluded_without_all]
	if !ok {
		return nil
	}

	fields := strings.Fields(expr)
	if len(fields) == 0 {
		return nil
	}

	return &excluded_without_allValidator{
		pass:       input.Pass,
		field:      input.Field,
		fields:     fields,
		structName: input.StructName,
		ruleName:   input.RuleName,
		parentPath: input.ParentPath,
	}
}
