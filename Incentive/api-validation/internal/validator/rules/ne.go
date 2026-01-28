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

type neValidator struct {
	pass        *codegen.Pass
	field       *ast.Field
	neValue     string
	expressions map[string]string
	structName  string
	ruleName    string
	parentPath  string
}

var _ validator.Validator = (*neValidator)(nil)

const neKey = "%s-ne"

func (n *neValidator) Validate() string {
	fieldName := n.FieldName()
	validation := fmt.Sprintf("!(t.%s != %s)", fieldName, n.neValue)

	// Get the field type for omitempty handling
	typ := n.pass.TypesInfo.TypeOf(n.field.Type)

	// Wrap with omitempty logic if needed
	return wrapWithOmitEmpty(n.pass, typ, fieldName, validation, n.expressions)
}

func (n *neValidator) FieldName() string {
	return n.field.Names[0].Name
}

func (n *neValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(n.structName, n.parentPath, n.FieldName())
}

func (n *neValidator) Err() string {
	key := fmt.Sprintf(neKey, n.structName+n.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the field equals [@VALUE] but should not.
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must not equal [@VALUE]", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sNeValidation", n.structName, n.FieldName())
	currentErrVarName := n.ErrVariable()

	// Escape quotes in the value for error message
	escapedValue := strings.ReplaceAll(n.neValue, `"`, `\"`)

	replacer := strings.NewReplacer(
		"[@ERRVARIABLE]", currentErrVarName,
		"[@LEGACYERRVAR]", legacyErrVarName,
		"[@FIELD]", n.FieldName(),
		"[@PATH]", n.FieldPath().String(),
		"[@VALUE]", escapedValue,
		"[@TYPE]", n.ruleName,
	)

	if currentErrVarName != legacyErrVarName {
		return replacer.Replace(deprecationNoticeTemplate + errTemplate)
	}

	return replacer.Replace(errTemplate)
}

func (n *neValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]NeValidation", "[@PATH]", n.FieldPath().CleanedPath())
}

func (n *neValidator) Imports() []string {
	return []string{}
}

// ValidateNe creates a new neValidator if the field is comparable and the ne marker is present.
func ValidateNe(input registry.ValidatorInput) validator.Validator {
	typ := input.Pass.TypesInfo.TypeOf(input.Field.Type)

	// Ensure the type is comparable (string, numeric, bool)
	if !types.Comparable(typ) {
		return nil
	}

	neValue, ok := input.Expressions[markers.GoValidMarkerNe]
	if !ok {
		return nil
	}

	// For string types, wrap the value in quotes if not already quoted
	if basic, ok := typ.Underlying().(*types.Basic); ok && basic.Kind() == types.String {
		if !strings.HasPrefix(neValue, `"`) && !strings.HasPrefix(neValue, "`") {
			neValue = fmt.Sprintf(`"%s"`, neValue)
		}
	}

	return &neValidator{
		pass:        input.Pass,
		field:       input.Field,
		neValue:     neValue,
		expressions: input.Expressions,
		structName:  input.StructName,
		ruleName:    input.RuleName,
		parentPath:  input.ParentPath,
	}
}
