package govalid

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"text/template"

	"github.com/gostaticanalysis/codegen"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/imports"

	// Local package import
	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/analyzers/markers"
	govaliderrors "gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/errors"
	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/validator"
	"gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/validator/registry"
)

const (
	// Name is the name of the govalid generator.
	Name = "govalid"
	// Doc is the documentation for the govalid generator.
	Doc = "govalid generates type-safe validation code for structs based on markers."
)

var (
	// dryRun indicates whether the generator should run in dry-run mode.
	dryRun bool
)

// generator is the main type for the govalid analyzer.
type generator struct{}

// newGenerator creates a new instance of the govalid generator.
func newGenerator() (*codegen.Generator, error) {
	g := &generator{}

	generator := &codegen.Generator{
		Name:     Name,
		Doc:      Doc,
		Run:      g.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer, markers.Analyzer},
	}

	return generator, nil
}

// TemplateData holds the data for the template used to generate validation code.
type TemplateData struct {
	PackageName    string
	TypeName       string
	Metadata       []*AnalyzedMetadata
	ImportPackages map[string]struct{}
}

// run is the main function that runs the govalid analyzer.
func (g *generator) run(pass *codegen.Pass) error {
	inspector, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return govaliderrors.ErrCouldNotGetInspector
	}

	markersInspect, ok := pass.ResultOf[markers.Analyzer].(markers.Markers)
	if !ok {
		return govaliderrors.ErrCouldNotGetInspector
	}

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	// Build a map of named struct types in the package for resolving dive element types.
	typeMap := map[string]*ast.StructType{}
	inspector.Preorder(nodeFilter, func(n ast.Node) {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			return
		}
		for _, spec := range genDecl.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if st, ok := ts.Type.(*ast.StructType); ok {
				typeMap[ts.Name.Name] = st
			}
		}
	})

	tmplList := map[string]TemplateData{}

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			return
		}

		for _, spec := range genDecl.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				return
			}

			typeMarkers := markersInspect.TypeMarkers(ts)

			structType, ok := ts.Type.(*ast.StructType)
			if !ok {
				return
			}

			// Check if struct has any validation tags before skipping
			hasValidationTags := hasAnyValidationTags(pass, markersInspect, structType)

			metadata := analyzeMarker(pass, markersInspect, typeMarkers, structType, "", ts.Name.Name, typeMap)

			// Skip only if there are no validation tags at all
			// If metadata is empty but validation tags exist (e.g., only omitempty), still generate the validator
			if len(metadata) == 0 && !hasValidationTags {
				return
			}

			// Consolidate validators with the same ParentVariable into single loops for performance
			metadata = consolidateMetadata(metadata)

			tmplData := TemplateData{
				PackageName:    pass.Pkg.Name(),
				TypeName:       ts.Name.Name,
				Metadata:       metadata,
				ImportPackages: collectImportPackages(metadata),
			}

			data, ok := tmplList[ts.Name.Name]
			if ok {
				data.Metadata = append(data.Metadata, tmplData.Metadata...)
			}

			if err := writeFile(pass, ts, tmplData); err != nil {
				panic(fmt.Sprintf("failed to write file for %s: %v", ts.Name.Name, err))
			}
		}
	})

	return nil
}

// AnalyzedMetadata holds the metadata for a field in a struct, including its validators and parent variable name.
type AnalyzedMetadata struct {
	Validators     []validator.Validator
	ParentVariable string
}

// consolidateMetadata merges AnalyzedMetadata entries with the same ParentVariable
// to generate a single loop instead of multiple loops for better performance.
// This is especially important for dive directives on collections.
// Only consolidates indexed parents (containing '[i]') to avoid changing behavior
// for non-collection nested structs.
func consolidateMetadata(metadata []*AnalyzedMetadata) []*AnalyzedMetadata {
	if len(metadata) == 0 {
		return metadata
	}

	// Map to collect validators by ParentVariable (only for indexed parents)
	parentMap := make(map[string]*AnalyzedMetadata)
	var order []string // Track insertion order
	result := make([]*AnalyzedMetadata, 0, len(metadata))

	for _, meta := range metadata {
		key := meta.ParentVariable

		// Only consolidate indexed parents (dive on collections)
		// Non-indexed parents should keep their original structure
		isIndexed := strings.Contains(key, "[i]")

		if isIndexed {
			if existing, ok := parentMap[key]; ok {
				// Merge validators into existing entry
				existing.Validators = append(existing.Validators, meta.Validators...)
			} else {
				// Create new entry
				newMeta := &AnalyzedMetadata{
					Validators:     meta.Validators,
					ParentVariable: meta.ParentVariable,
				}
				parentMap[key] = newMeta
				order = append(order, key)
			}
		} else {
			// Keep non-indexed entries as-is in their original position
			result = append(result, meta)
		}
	}

	// Append consolidated indexed entries at the end in their order
	for _, key := range order {
		result = append(result, parentMap[key])
	}

	return result
}

// hasAnyValidationTags checks if a struct has any fields with validate tags,
// including checking embedded structs for their field markers.
func hasAnyValidationTags(pass *codegen.Pass, markersInspect markers.Markers, structType *ast.StructType) bool {
	if structType == nil || structType.Fields == nil {
		return false
	}

	for _, field := range structType.Fields.List {
		// Check if field has validation tag
		if field.Tag != nil {
			tagValue := strings.Trim(field.Tag.Value, "`")
			if tagValue != "" {
				structTag := reflect.StructTag(tagValue)
				validateTag := structTag.Get("validate")
				if validateTag != "" {
					return true
				}
			}
		}

		// Check if this is an embedded struct (no field names)
		if len(field.Names) == 0 {
			// Check markers for the embedded type
			fieldMarkers := markersInspect.FieldMarkers(field)
			if len(fieldMarkers) > 0 {
				return true
			}

			// Also check if the embedded type itself is a struct with validation tags
			if fieldType, ok := field.Type.(*ast.StructType); ok {
				if hasAnyValidationTags(pass, markersInspect, fieldType) {
					return true
				}
			}

			// For embedded structs from other packages (e.g., port.MetadataRequest),
			// we assume they might have validation and return true to generate a validator
			if _, ok := field.Type.(*ast.SelectorExpr); ok {
				return true
			}
			if _, ok := field.Type.(*ast.Ident); ok {
				// Embedded named struct - assume it might need validation
				return true
			}
		}
	}

	return false
}

// makeValidatorInput contains all the input parameters needed for makeValidator function.
type makeValidatorInput struct {
	Pass       *codegen.Pass
	Markers    []markers.Marker
	Field      *ast.Field
	StructName string
	ParentPath string
}

//nolint:funlen // This function is complex but cohesive - it handles complete field analysis including nested structs
func analyzeMarker(pass *codegen.Pass, markersInspect markers.Markers, typeMarkers markers.MarkerSet, structType *ast.StructType, parent, structName string, typeMap map[string]*ast.StructType) []*AnalyzedMetadata {
	analyzed := make([]*AnalyzedMetadata, 0)

	typeMarkersList := make([]markers.Marker, 0, len(typeMarkers))
	for _, marker := range typeMarkers {
		typeMarkersList = append(typeMarkersList, marker)
	}

	sort.SliceStable(typeMarkersList, func(i, j int) bool {
		return typeMarkersList[i].Identifier < typeMarkersList[j].Identifier
	})

	for _, field := range structType.Fields.List {
		validators := make([]validator.Validator, 0)

		// Handle embedded structs (fields with no names)
		if len(field.Names) == 0 {
			// For embedded structs, recursively analyze their fields
			if st, ok := field.Type.(*ast.StructType); ok {
				// Embedded inline struct
				analyzed = append(analyzed, analyzeMarker(pass, markersInspect, typeMarkers, st, parent, structName, typeMap)...)
			} else if ident, ok := field.Type.(*ast.Ident); ok {
				// Embedded named struct from same package
				if embeddedStruct, ok := typeMap[ident.Name]; ok {
					analyzed = append(analyzed, analyzeMarker(pass, markersInspect, nil, embeddedStruct, parent, structName, typeMap)...)
				}
			} else if _, ok := field.Type.(*ast.SelectorExpr); ok {
				// Embedded struct from another package (e.g., port.MetadataRequest)
				// We can't analyze it directly, but we can check if it has field markers
				fieldMarkers := markersInspect.FieldMarkers(field)
				if len(fieldMarkers) > 0 {
					// There are markers on fields of the embedded struct
					// We'll need to process them but can't recurse into external package
					// This is a limitation - external package embedded structs need their own validation
				}
			}
			continue
		}

		// Apply markers to the field
		fieldMarkers := markersInspect.FieldMarkers(field)

		fieldMarkersList := make([]markers.Marker, 0, len(fieldMarkers))
		for _, marker := range fieldMarkers {
			fieldMarkersList = append(fieldMarkersList, marker)
		}

		sort.SliceStable(fieldMarkersList, func(i, j int) bool {
			return fieldMarkersList[i].Identifier < fieldMarkersList[j].Identifier
		})

		markersList := make([]markers.Marker, 0, len(typeMarkersList)+len(fieldMarkersList))
		markersList = append(markersList, typeMarkersList...)
		markersList = append(markersList, fieldMarkersList...)

		input := makeValidatorInput{
			Pass:       pass,
			Markers:    markersList,
			Field:      field,
			StructName: structName,
			ParentPath: parent,
		}

		// Check for dive marker on collection types to validate nested elements.
		hasDive := false
		for _, m := range markersList {
			if m.Identifier == "govalid:dive" {
				hasDive = true
				break
			}
		}

		// If not a struct and no dive, treat as a regular field.
		if _, ok := field.Type.(*ast.StructType); !ok && !hasDive {
			validators = makeValidator(input)
			if len(validators) == 0 {
				continue
			}

			analyzed = append(analyzed, &AnalyzedMetadata{
				Validators:     validators,
				ParentVariable: parent,
			})

			continue
		}

		// If the field is an inline struct, process it as before.
		if st, ok := field.Type.(*ast.StructType); ok && !hasDive {
			for _, f := range st.Fields.List {
				input.Field = f
				validators = append(validators, makeValidator(input)...)
			}

			var parentVariable string
			if parent != "" {
				parentVariable = fmt.Sprintf("%s.%s", parent, field.Names[0].Name)
			} else {
				parentVariable = field.Names[0].Name
			}

			if len(validators) > 0 {
				analyzed = append(analyzed, &AnalyzedMetadata{
					Validators:     validators,
					ParentVariable: parentVariable,
				})
			}

			// Recursively analyze nested inline structs
			analyzed = append(analyzed, analyzeMarker(pass, markersInspect, typeMarkers, st, parentVariable, structName, typeMap)...)
			continue
		}

		// Handle dive for slices/arrays of structs (including named or pointer-to-struct elements),
		// and also for direct/named struct fields when 'dive' is present.
		if hasDive {
			// Keep field-level validators (except dive itself)
			filtered := make([]markers.Marker, 0, len(markersList))
			for _, m := range markersList {
				if m.Identifier == "govalid:dive" {
					continue
				}
				filtered = append(filtered, m)
			}
			input.Markers = filtered
			validators = makeValidator(input)
			if len(validators) > 0 {
				analyzed = append(analyzed, &AnalyzedMetadata{
					Validators:     validators,
					ParentVariable: parent,
				})
			}

			// First, try to resolve element struct type for collections and validate elements.
			elt := resolveElementType(field.Type)
			if elt != nil {
				if elStruct := resolveStructTypeFromExpr(elt, typeMap); elStruct != nil {
					// Parent variable becomes Field[i]
					var base string
					if parent != "" {
						base = fmt.Sprintf("%s.%s", parent, field.Names[0].Name)
					} else {
						base = field.Names[0].Name
					}
					idxParent := fmt.Sprintf("%s[i]", base)
					// Analyze element struct using the parent type name to keep full path
					analyzed = append(analyzed, analyzeMarker(pass, markersInspect, nil, elStruct, idxParent, structName, typeMap)...)
				}
				// Done handling collection dive.
				continue
			}

			// If not a collection, support dive for direct/named struct fields as well.
			if target := resolveStructTypeFromExpr(field.Type, typeMap); target != nil {
				// Parent variable becomes Field (no index)
				var base string
				if parent != "" {
					base = fmt.Sprintf("%s.%s", parent, field.Names[0].Name)
				} else {
					base = field.Names[0].Name
				}
				// Analyze nested struct
				analyzed = append(analyzed, analyzeMarker(pass, markersInspect, nil, target, base, structName, typeMap)...)
				continue
			}

			continue
		}

		// Fallback: treat as regular field if none of the above matched
		validators = makeValidator(input)
		if len(validators) == 0 {
			continue
		}
		analyzed = append(analyzed, &AnalyzedMetadata{
			Validators:     validators,
			ParentVariable: parent,
		})
	}

	return analyzed
}

// resolveElementType extracts the element expression from array/slice/pointer types.
func resolveElementType(expr ast.Expr) ast.Expr {
	switch t := expr.(type) {
	case *ast.ArrayType:
		return t.Elt
	case *ast.StarExpr:
		return resolveElementType(t.X)
	default:
		return nil
	}
}

// resolveStructTypeFromExpr resolves a struct type from an ast.Expr using the provided type map.
func resolveStructTypeFromExpr(expr ast.Expr, typeMap map[string]*ast.StructType) *ast.StructType {
	switch e := expr.(type) {
	case *ast.StructType:
		return e
	case *ast.StarExpr:
		return resolveStructTypeFromExpr(e.X, typeMap)
	case *ast.Ident:
		if st, ok := typeMap[e.Name]; ok {
			return st
		}
		return nil
	default:
		return nil
	}
}

func makeValidator(input makeValidatorInput) []validator.Validator {
	validators := make([]validator.Validator, 0)

	for _, marker := range input.Markers {
		factory, err := registry.Validator(marker.Identifier)
		if err != nil {
			// Validator not found, skip
			continue
		}

		ruleName := strings.TrimPrefix(marker.Identifier, "govalid:")

		validatorInput := registry.ValidatorInput{
			Pass:        input.Pass,
			Field:       input.Field,
			Expressions: marker.Expressions,
			StructName:  input.StructName,
			RuleName:    ruleName,
			ParentPath:  input.ParentPath,
		}
		v := factory(validatorInput)

		if v == nil {
			continue
		}

		validators = append(validators, v)
	}

	return validators
}

// collectImportPackages analyzes validators and collects required import packages.
func collectImportPackages(metadata []*AnalyzedMetadata) map[string]struct{} {
	packages := make(map[string]struct{})

	for _, meta := range metadata {
		for _, validator := range meta.Validators {
			for _, pkg := range validator.Imports() {
				packages[pkg] = struct{}{}
			}
		}
	}

	return packages
}

func writeFile(pass *codegen.Pass, ts *ast.TypeSpec, tmplData TemplateData) error {
	t, err := template.New("validator").Funcs(template.FuncMap{
		"trimDots": func(s string) string {
			return strings.ReplaceAll(s, ".", "")
		},
		"hasIndex": func(s string) bool {
			return strings.Contains(s, "[i]")
		},
		"indexBase": func(s string) string {
			// convert something like Parent.Field[i] to Parent.Field
			return strings.TrimSuffix(s, "[i]")
		},
	}).Parse(ValidationTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, tmplData); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Use goimports to format the source code with proper import grouping
	src, err := imports.Process("", buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("failed to format source code with imports: %w", err)
	}

	src, err = format.Source(src)
	if err != nil {
		return fmt.Errorf("failed to format source code: %w", err)
	}

	if testing.Testing() || dryRun {
		if _, err := pass.Print(string(src)); err != nil {
			return fmt.Errorf("failed to print source code: %w", err)
		}

		return nil
	}

	originalFilePath := pass.Fset.Position(ts.Pos()).Filename
	fileName := strings.TrimSuffix(filepath.Base(originalFilePath), filepath.Ext(originalFilePath))
	typeName := ts.Name.Name
	fileName = fmt.Sprintf("%s_%s_validator.go", fileName, strings.ToLower(typeName))

	file, err := os.Create(filepath.Clean(fileName))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("failed to close file: %v\n", err)
		}
	}()

	if _, err := fmt.Fprint(file, string(src)); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
