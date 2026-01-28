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

type isdefaultValidator struct {
	pass       *codegen.Pass
	field      *ast.Field
	structName string
	ruleName   string
	parentPath string
}

var _ validator.Validator = (*isdefaultValidator)(nil)

const isdefaultKey = "%s-isdefault"

func (i *isdefaultValidator) Validate() string {
	typ := i.pass.TypesInfo.TypeOf(i.field.Type)
	zero := getZeroValue(typ)
	name := i.FieldName()

	return fmt.Sprintf("t.%s != %s", name, zero)
}

func (i *isdefaultValidator) FieldName() string {
	return i.field.Names[0].Name
}

func (i *isdefaultValidator) FieldPath() validator.FieldPath {
	return validator.NewFieldPath(i.structName, i.parentPath, i.FieldName())
}

func (i *isdefaultValidator) Err() string {
	key := fmt.Sprintf(isdefaultKey, i.structName+i.FieldPath().CleanedPath())

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
		// [@ERRVARIABLE] is the error returned when the field is not at its default/zero value.
		[@ERRVARIABLE] = govaliderrors.ValidationError{Reason: "field [@FIELD] must be at its default value", Path: "[@PATH]", Type: "[@TYPE]"}
	`

	legacyErrVarName := fmt.Sprintf("Err%s%sIsdefaultValidation", i.structName, i.FieldName())
	currentErrVarName := i.ErrVariable()

	replacer := strings.NewReplacer(
		"[@ERRVARIABLE]", currentErrVarName,
		"[@LEGACYERRVAR]", legacyErrVarName,
		"[@FIELD]", i.FieldName(),
		"[@PATH]", i.FieldPath().String(),
		"[@TYPE]", i.ruleName,
	)

	if currentErrVarName != legacyErrVarName {
		return replacer.Replace(deprecationNoticeTemplate + errTemplate)
	}

	return replacer.Replace(errTemplate)
}

func (i *isdefaultValidator) ErrVariable() string {
	return strings.ReplaceAll("Err[@PATH]IsdefaultValidation", "[@PATH]", i.FieldPath().CleanedPath())
}

func (i *isdefaultValidator) Imports() []string {
	return []string{}
}

// ValidateIsdefault creates a new isdefaultValidator for the given field.
func ValidateIsdefault(input registry.ValidatorInput) validator.Validator {
	return &isdefaultValidator{
		pass:       input.Pass,
		field:      input.Field,
		structName: input.StructName,
		ruleName:   input.RuleName,
		parentPath: input.ParentPath,
	}
}

// getZeroValue returns the zero value representation for a given type.
func getZeroValue(typ types.Type) string {
	switch t := typ.Underlying().(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.String:
			return `""`
		case types.Bool:
			return "false"
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
			types.Float32, types.Float64, types.Complex64, types.Complex128:
			return "0"
		}
	case *types.Pointer, *types.Slice, *types.Map, *types.Chan, *types.Interface:
		return "nil"
	case *types.Struct:
		// For structs, we use a zero value comparison
		return fmt.Sprintf("%s{}", typ.String())
	}
	return "nil"
}
