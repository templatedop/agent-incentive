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

type excludesallValidator struct {
	pass       *codegen.Pass
	field      *ast.Field
	chars      string
	structName string
	ruleName   string
	parentPath string
}

var _ validator.Validator = (*excludesallValidator)(nil)

const excludesallKey = "%s-excludesall"

func (e *excludesallValidator) Validate() string {
	fieldName := e.FieldName()
	return fmt.Sprintf("strings.ContainsAny(t.%s, %q)", fieldName, e.chars)
}

func (e *excludesallValidator) FieldName() string {
	return e.field.Names[0].Name
}

func (e *excludesallValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(e.structName, e.parentPath, e.FieldName())
}

func (e *excludesallValidator) Err() string {
	key := fmt.Sprintf(excludesallKey, e.structName+e.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the field contains any of the excluded characters.
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must not contain any of these characters: [@CHARS]", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sExcludesallValidation", e.structName, e.FieldName())
	currentErrVarName := e.ErrVariable()

	replacer := strings.NewReplacer(
		"[@ERRVARIABLE]", currentErrVarName,
		"[@LEGACYERRVAR]", legacyErrVarName,
		"[@FIELD]", e.FieldName(),
		"[@PATH]", e.FieldPath().String(),
		"[@CHARS]", e.chars,
		"[@TYPE]", e.ruleName,
	)

	if currentErrVarName != legacyErrVarName {
		return replacer.Replace(deprecationNoticeTemplate + errTemplate)
	}

	return replacer.Replace(errTemplate)
}

func (e *excludesallValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]ExcludesallValidation", "[@PATH]", e.FieldPath().CleanedPath())
}

func (e *excludesallValidator) Imports() []string {
	return []string{"strings"}
}

// ValidateExcludesall creates a new excludesallValidator for string types.
func ValidateExcludesall(input registry.ValidatorInput) validator.Validator {
	typ := input.Pass.TypesInfo.TypeOf(input.Field.Type)

	// Check if it's a string type
	basic, ok := typ.Underlying().(*types.Basic)
	if !ok || basic.Kind() != types.String {
		return nil
	}

	chars, ok := input.Expressions[markers.GoValidMarkerExcludesall]
	if !ok {
		return nil
	}

	return &excludesallValidator{
		pass:       input.Pass,
		field:      input.Field,
		chars:      chars,
		structName: input.StructName,
		ruleName:   input.RuleName,
		parentPath: input.ParentPath,
	}
}
