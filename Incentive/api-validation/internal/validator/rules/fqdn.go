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

type fqdnValidator struct {
	pass        *codegen.Pass
	field       *ast.Field
	expressions map[string]string
	structName  string
	ruleName    string
	parentPath  string
}

var _ validator.Validator = (*fqdnValidator)(nil)

const fqdnKey = "%s-fqdn"

func (v *fqdnValidator) Validate() string {
	fieldName := v.FieldName()
	validation := fmt.Sprintf("!validationhelper.IsValidFQDN(t.%s)", fieldName)

	// Get the field type for omitempty handling
	typ := v.pass.TypesInfo.TypeOf(v.field.Type)

	// Wrap with omitempty logic if needed
	return wrapWithOmitEmpty(v.pass, typ, fieldName, validation, v.expressions)
}

func (v *fqdnValidator) FieldName() string {
	return v.field.Names[0].Name
}

func (v *fqdnValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(v.structName, v.parentPath, v.FieldName())
}

func (v *fqdnValidator) Err() string {
	key := fmt.Sprintf(fqdnKey, v.structName+v.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the field is not a fully qualified domain name.
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must be a fully qualified domain name", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sFQDNValidation", v.structName, v.FieldName())
	currentErrVarName := v.ErrVariable()

	replacer := strings.NewReplacer(
		"[@ERRVARIABLE]", currentErrVarName,
		"[@LEGACYERRVAR]", legacyErrVarName,
		"[@FIELD]", v.FieldName(),
		"[@PATH]", v.FieldPath().String(),
		"[@TYPE]", v.ruleName,
	)

	if currentErrVarName != legacyErrVarName {
		return replacer.Replace(deprecationNoticeTemplate + errTemplate)
	}

	return replacer.Replace(errTemplate)
}

func (v *fqdnValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]FQDNValidation", `[@PATH]`, v.FieldPath().CleanedPath())
}

func (v *fqdnValidator) Imports() []string {
	return []string{"gitlab.cept.gov.in/it-2.0-common/n-api-validation/validation/validationhelper"}
}

// ValidateFQDN creates a new fqdnValidator for string types.
func ValidateFQDN(input registry.ValidatorInput) validator.Validator {
	typ := input.Pass.TypesInfo.TypeOf(input.Field.Type)

	// Check if it's a string type
	basic, ok := typ.Underlying().(*types.Basic)
	if !ok || basic.Kind() != types.String {
		return nil
	}

	return &fqdnValidator{
		pass:        input.Pass,
		field:       input.Field,
		expressions: input.Expressions,
		structName:  input.StructName,
		ruleName:    input.RuleName,
		parentPath:  input.ParentPath,
	}
}
