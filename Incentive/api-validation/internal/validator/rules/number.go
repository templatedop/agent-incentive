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

type numberValidator struct {
	pass        *codegen.Pass
	field       *ast.Field
	expressions map[string]string
	structName  string
	ruleName    string
	parentPath  string
}

var _ validator.Validator = (*numberValidator)(nil)

const numberKey = "%s-number"

func (n *numberValidator) Validate() string {
	fieldName := n.FieldName()
	validation := fmt.Sprintf("!validationhelper.IsNumber(t.%s)", fieldName)

	// Get the field type for omitempty handling
	typ := n.pass.TypesInfo.TypeOf(n.field.Type)

	// Wrap with omitempty logic if needed
	return wrapWithOmitEmpty(n.pass, typ, fieldName, validation, n.expressions)
}

func (n *numberValidator) FieldName() string {
	return n.field.Names[0].Name
}

func (n *numberValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(n.structName, n.parentPath, n.FieldName())
}

func (n *numberValidator) Err() string {
	key := fmt.Sprintf(numberKey, n.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the field contains non-numeric characters.
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must contain only numbers", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sNumberValidation", n.structName, n.FieldName())
	currentErrVarName := n.ErrVariable()

	replacer := strings.NewReplacer(
		"[@ERRVARIABLE]", currentErrVarName,
		"[@LEGACYERRVAR]", legacyErrVarName,
		"[@FIELD]", n.FieldName(),
		"[@PATH]", n.FieldPath().String(),
		"[@TYPE]", n.ruleName,
	)

	if currentErrVarName != legacyErrVarName {
		return replacer.Replace(deprecationNoticeTemplate + errTemplate)
	}

	return replacer.Replace(errTemplate)
}

func (n *numberValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]NumberValidation", `[@PATH]`, n.FieldPath().CleanedPath())
}

func (n *numberValidator) Imports() []string {
	return []string{"gitlab.cept.gov.in/it-2.0-common/n-api-validation/validation/validationhelper"}
}

// ValidateNumber creates a new numberValidator for string types.
func ValidateNumber(input registry.ValidatorInput) validator.Validator {
	typ := input.Pass.TypesInfo.TypeOf(input.Field.Type)

	// Check if it's a string type
	basic, ok := typ.Underlying().(*types.Basic)
	if !ok || basic.Kind() != types.String {
		return nil
	}

	return &numberValidator{
		pass:        input.Pass,
		field:       input.Field,
		expressions: input.Expressions,
		structName:  input.StructName,
		ruleName:    input.RuleName,
		parentPath:  input.ParentPath,
	}
}
