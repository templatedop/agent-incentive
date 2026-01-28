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

type excluded_withValidator struct {
	pass       *codegen.Pass
	field      *ast.Field
	fields     []string
	structName string
	ruleName   string
	parentPath string
}

var _ validator.Validator = (*excluded_withValidator)(nil)

const excluded_withKey = "%s-excluded_with"

func (e *excluded_withValidator) Validate() string {
	typ := e.pass.TypesInfo.TypeOf(e.field.Type)
	zero := getZeroValueForType(typ)
	fieldName := e.FieldName()

	// Generate: (field1 != zero || field2 != zero || ...) && thisField != zero
	var conditions []string
	for _, f := range e.fields {
		conditions = append(conditions, fmt.Sprintf(`t.%s != ""`, f))
	}

	return fmt.Sprintf("(%s) && t.%s != %s",
		strings.Join(conditions, " || "), fieldName, zero)
}

func (e *excluded_withValidator) FieldName() string {
	return e.field.Names[0].Name
}

func (e *excluded_withValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(e.structName, e.parentPath, e.FieldName())
}

func (e *excluded_withValidator) Err() string {
	key := fmt.Sprintf(excluded_withKey, e.structName+e.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the field must be absent because other fields are present.
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must be absent when any of [@FIELDS] are present", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sExcludedWithValidation", e.structName, e.FieldName())
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

func (e *excluded_withValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]ExcludedWithValidation", "[@PATH]", e.FieldPath().CleanedPath())
}

func (e *excluded_withValidator) Imports() []string {
	return []string{}
}

// ValidateExcludedWith creates a new excluded_withValidator.
// Format: excluded_with=Field1 Field2 Field3...
func ValidateExcludedWith(input registry.ValidatorInput) validator.Validator {
	expr, ok := input.Expressions[markers.GoValidMarkerExcluded_with]
	if !ok {
		return nil
	}

	fields := strings.Fields(expr)
	if len(fields) == 0 {
		return nil
	}

	return &excluded_withValidator{
		pass:       input.Pass,
		field:      input.Field,
		fields:     fields,
		structName: input.StructName,
		ruleName:   input.RuleName,
		parentPath: input.ParentPath,
	}
}
