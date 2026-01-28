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

type required_with_allValidator struct {
	pass       *codegen.Pass
	field      *ast.Field
	fields     []string
	structName string
	ruleName   string
	parentPath string
}

var _ validator.Validator = (*required_with_allValidator)(nil)

const required_with_allKey = "%s-required_with_all"

func (r *required_with_allValidator) Validate() string {
	typ := r.pass.TypesInfo.TypeOf(r.field.Type)
	zero := getZeroValueForType(typ)
	fieldName := r.FieldName()

	// Generate: (field1 != zero && field2 != zero && ...) && thisField == zero
	var conditions []string
	for _, f := range r.fields {
		conditions = append(conditions, fmt.Sprintf(`t.%s != ""`, f))
	}

	return fmt.Sprintf("(%s) && t.%s == %s",
		strings.Join(conditions, " && "), fieldName, zero)
}

func (r *required_with_allValidator) FieldName() string {
	return r.field.Names[0].Name
}

func (r *required_with_allValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(r.structName, r.parentPath, r.FieldName())
}

func (r *required_with_allValidator) Err() string {
	key := fmt.Sprintf(required_with_allKey, r.structName+r.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the field is required because all other fields are present.
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] is required when all of [@FIELDS] are present", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sRequiredWithAllValidation", r.structName, r.FieldName())
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

func (r *required_with_allValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]RequiredWithAllValidation", "[@PATH]", r.FieldPath().CleanedPath())
}

func (r *required_with_allValidator) Imports() []string {
	return []string{}
}

// ValidateRequiredWithAll creates a new required_with_allValidator.
// Format: required_with_all=Field1 Field2 Field3...
func ValidateRequiredWithAll(input registry.ValidatorInput) validator.Validator {
	expr, ok := input.Expressions[markers.GoValidMarkerRequired_with_all]
	if !ok {
		return nil
	}

	fields := strings.Fields(expr)
	if len(fields) == 0 {
		return nil
	}

	return &required_with_allValidator{
		pass:       input.Pass,
		field:      input.Field,
		fields:     fields,
		structName: input.StructName,
		ruleName:   input.RuleName,
		parentPath: input.ParentPath,
	}
}
