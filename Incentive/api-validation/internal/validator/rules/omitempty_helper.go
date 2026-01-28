package rules

import (
	"fmt"
	"go/types"

	"github.com/gostaticanalysis/codegen"
)

// wrapWithOmitEmpty wraps a validation expression with omitempty logic if the omitempty flag is set.
// It checks if the field has a zero value before applying the validation.
func wrapWithOmitEmpty(pass *codegen.Pass, fieldType types.Type, fieldName string, validation string, expressions map[string]string) string {
	// Check if omitempty is enabled
	if expressions == nil || expressions["omitempty"] != "true" {
		return validation
	}

	// Determine the zero-value check based on the field type
	underlying := fieldType.Underlying()

	var zeroCheck string
	switch t := underlying.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.String:
			zeroCheck = fmt.Sprintf("t.%s != \"\"", fieldName)
		case types.Bool:
			// For bool, omitempty means validate only if true
			zeroCheck = fmt.Sprintf("t.%s", fieldName)
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
			zeroCheck = fmt.Sprintf("t.%s != 0", fieldName)
		case types.Float32, types.Float64:
			zeroCheck = fmt.Sprintf("t.%s != 0.0", fieldName)
		default:
			// Fallback for other basic types
			zeroCheck = fmt.Sprintf("t.%s != \"\"", fieldName)
		}
	case *types.Pointer:
		zeroCheck = fmt.Sprintf("t.%s != nil", fieldName)
	case *types.Slice, *types.Map, *types.Chan:
		zeroCheck = fmt.Sprintf("len(t.%s) != 0", fieldName)
	case *types.Interface:
		zeroCheck = fmt.Sprintf("t.%s != nil", fieldName)
	default:
		// For structs and other complex types, we can't easily check for zero value
		// So we just apply the validation (omitempty doesn't make much sense for these)
		return validation
	}

	// Wrap the validation with the zero-value check
	return fmt.Sprintf("%s && %s", zeroCheck, validation)
}
