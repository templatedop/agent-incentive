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

type containsanyValidator struct {
	pass        *codegen.Pass
	field       *ast.Field
	chars       string
	expressions map[string]string
	structName  string
	ruleName    string
	parentPath  string
}

var _ validator.Validator = (*containsanyValidator)(nil)

const containsanyKey = "%s-containsany"

func (c *containsanyValidator) Validate() string {
	fieldName := c.FieldName()
	validation := fmt.Sprintf("!strings.ContainsAny(t.%s, %q)", fieldName, c.chars)

	// Get the field type for omitempty handling
	typ := c.pass.TypesInfo.TypeOf(c.field.Type)

	// Wrap with omitempty logic if needed
	return wrapWithOmitEmpty(c.pass, typ, fieldName, validation, c.expressions)
}

func (c *containsanyValidator) FieldName() string {
	return c.field.Names[0].Name
}

func (c *containsanyValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(c.structName, c.parentPath, c.FieldName())
}

func (c *containsanyValidator) Err() string {
	key := fmt.Sprintf(containsanyKey, c.structName+c.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the field does not contain any of the specified characters.
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must contain at least one of these characters: [@CHARS]", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sContainsanyValidation", c.structName, c.FieldName())
	currentErrVarName := c.ErrVariable()

	replacer := strings.NewReplacer(
		"[@ERRVARIABLE]", currentErrVarName,
		"[@LEGACYERRVAR]", legacyErrVarName,
		"[@FIELD]", c.FieldName(),
		"[@PATH]", c.FieldPath().String(),
		"[@CHARS]", c.chars,
		"[@TYPE]", c.ruleName,
	)

	if currentErrVarName != legacyErrVarName {
		return replacer.Replace(deprecationNoticeTemplate + errTemplate)
	}

	return replacer.Replace(errTemplate)
}

func (c *containsanyValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]ContainsanyValidation", "[@PATH]", c.FieldPath().CleanedPath())
}

func (c *containsanyValidator) Imports() []string {
	return []string{"strings"}
}

// ValidateContainsany creates a new containsanyValidator for string types.
func ValidateContainsany(input registry.ValidatorInput) validator.Validator {
	typ := input.Pass.TypesInfo.TypeOf(input.Field.Type)

	// Check if it's a string type
	basic, ok := typ.Underlying().(*types.Basic)
	if !ok || basic.Kind() != types.String {
		return nil
	}

	chars, ok := input.Expressions[markers.GoValidMarkerContainsany]
	if !ok {
		return nil
	}

	return &containsanyValidator{
		pass:        input.Pass,
		field:       input.Field,
		chars:       chars,
		expressions: input.Expressions,
		structName:  input.StructName,
		ruleName:    input.RuleName,
		parentPath:  input.ParentPath,
	}
}
