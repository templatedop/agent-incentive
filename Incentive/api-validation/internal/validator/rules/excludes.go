// Package rules implements validation rules for fields in structs.
package rules

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/gostaticanalysis/codegen"

	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/markers"
	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/validator"
	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/validator/registry"
)

type excludesValidator struct {
	pass       *codegen.Pass
	field      *ast.Field
	substr     string
	structName string
	ruleName   string
	parentPath string
}

var _ validator.Validator = (*excludesValidator)(nil)

const excludesKey = "%s-excludes"

func (e *excludesValidator) Validate() string {
	fieldName := e.FieldName()
	return fmt.Sprintf("strings.Contains(t.%s, %q)", fieldName, e.substr)
}

func (e *excludesValidator) FieldName() string {
	return e.field.Names[0].Name
}

func (e *excludesValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(e.structName, e.parentPath, e.FieldName())
}

func (e *excludesValidator) Err() string {
	key := fmt.Sprintf(excludesKey, e.structName+e.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the field contains the excluded substring.
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must not contain: [@SUBSTR]", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sExcludesValidation", e.structName, e.FieldName())
	currentErrVarName := e.ErrVariable()

	replacer := strings.NewReplacer(
		"[@ERRVARIABLE]", currentErrVarName,
		"[@LEGACYERRVAR]", legacyErrVarName,
		"[@FIELD]", e.FieldName(),
		"[@PATH]", e.FieldPath().String(),
		"[@SUBSTR]", e.substr,
		"[@TYPE]", e.ruleName,
	)

	if currentErrVarName != legacyErrVarName {
		return replacer.Replace(deprecationNoticeTemplate + errTemplate)
	}

	return replacer.Replace(errTemplate)
}

func (e *excludesValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]ExcludesValidation", "[@PATH]", e.FieldPath().CleanedPath())
}

func (e *excludesValidator) Imports() []string {
	return []string{"strings"}
}

// ValidateExcludes creates a new excludesValidator for string types.
func ValidateExcludes(input registry.ValidatorInput) validator.Validator {
	typ := input.Pass.TypesInfo.TypeOf(input.Field.Type)

	// Check if it's a string type
	basic, ok := typ.Underlying().(*types.Basic)
	if !ok || basic.Kind() != types.String {
		return nil
	}

	substr, ok := input.Expressions[markers.GoValidMarkerExcludes]
	if !ok {
		return nil
	}

	return &excludesValidator{
		pass:       input.Pass,
		field:      input.Field,
		substr:     substr,
		structName: input.StructName,
		ruleName:   input.RuleName,
		parentPath: input.ParentPath,
	}
}
