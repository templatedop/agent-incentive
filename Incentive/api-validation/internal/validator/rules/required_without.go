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

type required_withoutValidator struct {
	pass       *codegen.Pass
	field      *ast.Field
	fields     []string
	structName string
	ruleName   string
	parentPath string
}

var _ validator.Validator = (*required_withoutValidator)(nil)

const required_withoutKey = "%s-required_without"

func (r *required_withoutValidator) Validate() string {
	typ := r.pass.TypesInfo.TypeOf(r.field.Type)
	zero := getZeroValueForType(typ)
	fieldName := r.FieldName()

	// Generate: (field1 == zero || field2 == zero || ...) && thisField == zero
	var conditions []string
	for _, f := range r.fields {
		conditions = append(conditions, fmt.Sprintf(`t.%s == ""`, f))
	}

	return fmt.Sprintf("(%s) && t.%s == %s",
		strings.Join(conditions, " || "), fieldName, zero)
}

func (r *required_withoutValidator) FieldName() string {
	return r.field.Names[0].Name
}

func (r *required_withoutValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(r.structName, r.parentPath, r.FieldName())
}

func (r *required_withoutValidator) Err() string {
	key := fmt.Sprintf(required_withoutKey, r.structName+r.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the field is required because other fields are absent.
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] is required when any of [@FIELDS] are absent", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sRequiredWithoutValidation", r.structName, r.FieldName())
	currentErrVarName := r.ErrVariable()

	replacer := strings.NewReplacer(
		"[@ERRVARIABLE]", currentErrVarName,
		"[@LEGACYERRVAR]", legacyErrVarName,
		"[@FIELD]", r.FieldName(),
		"[@PATH]", r.FieldPath().String(),
		"[@FIELDS]", strings.Join(r.fields, ", "),
		"[@TYPE]", r.ruleName,
	)

	if currentErrVarName != legacyErrVarName {
		return replacer.Replace(deprecationNoticeTemplate + errTemplate)
	}

	return replacer.Replace(errTemplate)
}

func (r *required_withoutValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]RequiredWithoutValidation", "[@PATH]", r.FieldPath().CleanedPath())
}

func (r *required_withoutValidator) Imports() []string {
	return []string{}
}

// ValidateRequiredWithout creates a new required_withoutValidator.
// Format: required_without=Field1 Field2 Field3...
func ValidateRequiredWithout(input registry.ValidatorInput) validator.Validator {
	expr, ok := input.Expressions[markers.GoValidMarkerRequired_without]
	if !ok {
		return nil
	}

	fields := strings.Fields(expr)
	if len(fields) == 0 {
		return nil
	}

	return &required_withoutValidator{
		pass:       input.Pass,
		field:      input.Field,
		fields:     fields,
		structName: input.StructName,
		ruleName:   input.RuleName,
		parentPath: input.ParentPath,
	}
}
