package markers

import (
	"go/ast"
	"go/types"
	"reflect"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	// Local package import
	govaliderrors "gitlab.cept.gov.in/it-2.0-common/n-api-validation/internal/errors"
)

const (
	// Name is the name of the markers analyzer.
	Name = "markers"
	// Doc is the documentation for the markers analyzer.
	Doc = "markers is a helper for generating govalid validation"
)

// Analyzer is the main entry point for the markers analyzer.
// This variable must be initialized by registry package.
var Analyzer *analysis.Analyzer

// analyzer implements the analysis.Analyzer interface for the markers analyzer.
type analyzer struct{}

// newAnalyzer creates a new instance of the markers analyzer.
func newAnalyzer() *analysis.Analyzer {
	a := &analyzer{}
	analyzer := &analysis.Analyzer{
		Name:       Name,
		Doc:        Doc,
		Run:        a.run,
		Requires:   []*analysis.Analyzer{inspect.Analyzer},
		ResultType: reflect.TypeOf(newMarkers()),
		FactTypes: []analysis.Fact{
			(*MarkerFact)(nil),
		},
	}

	return analyzer
}

// run is the main function that runs the markers analyzer.
func (a *analyzer) run(pass *analysis.Pass) (any, error) {
	inspect, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, govaliderrors.ErrCouldNotGetInspector
	}

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	results, ok := newMarkers().(*markers)
	if !ok {
		return nil, govaliderrors.ErrCouldNotCreateMarkers
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.GenDecl:
			if n == nil {
				return
			}

			if n.Doc != nil && len(n.Doc.List) > 0 {
				collectTypeMarkers(pass, n, results)
			}

			for _, spec := range n.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					return
				}

				collectStructMarkers(pass, st, results)
			}
		default:
		}
	})

	return results, nil
}

// collectTypeMarkers collects markers from a GenDecl node and adds them to the results.
func collectTypeMarkers(pass *analysis.Pass, genDecl *ast.GenDecl, results *markers) {
	if genDecl.Doc == nil || len(genDecl.Doc.List) == 0 {
		return
	}

	for _, doc := range genDecl.Doc.List {
		if !strings.HasPrefix(doc.Text, "// +") {
			continue
		}

		markerContent := strings.TrimPrefix(doc.Text, "// +")

		identifier, expressions := extractMarker(markerContent)
		marker := Marker{
			Identifier:  identifier,
			Expressions: expressions,
		}

		for _, spec := range genDecl.Specs {
			if ts, ok := spec.(*ast.TypeSpec); ok {
				results.insertTypeMarker(ts, marker)

				if obj, ok := pass.TypesInfo.Defs[ts.Name]; ok {
					pass.ExportObjectFact(obj, &MarkerFact{
						Identifier:  identifier,
						Expressions: expressions,
					})
				}
			}
		}
	}
}

// collectStructMarkers collects markers from a TypeSpec node and adds them to the results.
func collectStructMarkers(pass *analysis.Pass, s *ast.StructType, results *markers) {
	if s == nil || s.Fields == nil || len(s.Fields.List) == 0 {
		return
	}

	for _, field := range s.Fields.List {
		fieldMarkers(pass, field, results)

		structType, ok := field.Type.(*ast.StructType)
		if !ok {
			continue
		}

		collectStructMarkers(pass, structType, results)
	}
}

// fieldMarkers extracts markers from a struct field and adds them to the results.
func fieldMarkers(pass *analysis.Pass, field *ast.Field, results *markers) {
	if field == nil || len(field.Names) == 0 {
		return
	}

	// Support legacy comment-based markers for backward compatibility.
	if field.Doc != nil && len(field.Doc.List) > 0 {
		for _, doc := range field.Doc.List {
			if !strings.HasPrefix(doc.Text, "// +") {
				continue
			}

			markerContent := strings.TrimPrefix(doc.Text, "// +")

			identifier, expressions := extractMarker(markerContent)
			marker := Marker{
				Identifier:  identifier,
				Expressions: expressions,
			}
			results.insertFieldMarker(field, marker)

			if obj, ok := pass.TypesInfo.Defs[field.Names[0]]; ok {
				pass.ExportObjectFact(obj, &MarkerFact{
					Identifier:  identifier,
					Expressions: expressions,
				})
			}
		}
	}

	// New tag-based markers: parse the `validate:"..."` struct tag.
	if field.Tag == nil {
		return
	}

	tagValue := strings.Trim(field.Tag.Value, "`")
	if tagValue == "" {
		return
	}

	structTag := reflect.StructTag(tagValue)
	validateRaw := structTag.Get("validate")
	if validateRaw == "" {
		return
	}

	// Check if omitempty is present
	hasOmitEmpty := false
	tokens := strings.Split(validateRaw, ",")
	for _, tok := range tokens {
		if strings.TrimSpace(tok) == "omitempty" {
			hasOmitEmpty = true
			break
		}
	}

	// Split validators by comma, e.g. `required,email,lt=10,max=5`
	for _, tok := range tokens {
		v := strings.TrimSpace(tok)
		if v == "" {
			continue
		}

		// Skip omitempty token itself - it's a modifier, not a validator
		if v == "omitempty" {
			continue
		}

		identifier, expressions := normalizeValidateToken(pass, field, v)
		if identifier == "" {
			continue
		}

		// If omitempty is present, add it to the expressions map
		// This tells the code generator to make this validation conditional
		if hasOmitEmpty {
			if expressions == nil {
				expressions = make(map[string]string)
			}
			expressions["omitempty"] = "true"
		}

		marker := Marker{Identifier: identifier, Expressions: expressions}
		results.insertFieldMarker(field, marker)

		if obj, ok := pass.TypesInfo.Defs[field.Names[0]]; ok {
			fact := &MarkerFact{Identifier: identifier, Expressions: expressions}
			pass.ExportObjectFact(obj, fact)
		}
	}
}

// normalizeValidateToken converts a single validate tag token into a marker identifier and expressions map.
// It supports synonyms and type-based disambiguation (e.g., max -> maxlength or maxitems).
func normalizeValidateToken(pass *analysis.Pass, field *ast.Field, token string) (string, map[string]string) {
	var key string
	var value string

	if strings.Contains(token, "=") {
		parts := strings.SplitN(token, "=", 2)
		key = strings.TrimSpace(parts[0])
		value = strings.TrimSpace(parts[1])
	} else {
		key = token
	}

	keyLower := strings.ToLower(key)

	// Determine field type for contextual mapping of min/max.
	typ := pass.TypesInfo.TypeOf(field.Type)
	underlying := typ
	if typ != nil {
		underlying = typ.Underlying()
	}

	isString := false
	isCollection := false

	switch u := underlying.(type) {
	case *types.Basic:
		isString = u.Kind() == types.String
	case *types.Slice, *types.Map, *types.Chan, *types.Array:
		isCollection = true
	}

	// Map synonyms to internal marker names.
	switch keyLower {
	// Original validators
	case "required", "email", "uuid", "url", "numeric", "ipv4", "ipv6", "alpha", "enum", "cel", "lt", "lte", "gt", "gte", "length", "date", "dive":
		// direct mapping
	// New simple validators
	case "eq", "ne", "isdefault", "boolean", "lowercase", "oneof", "number", "alphanum":
		// direct mapping
	// String pattern validators
	case "containsany", "excludes", "excludesall":
		// direct mapping
	// Collection validators
	case "unique":
		// direct mapping
	// Format validators
	case "uri", "fqdn", "latitude", "longitude", "iscolour", "iscolor":
		if keyLower == "iscolor" {
			keyLower = "iscolour" // US spelling -> UK spelling
		}
		// direct mapping
	// Time/Duration validators
	case "minduration", "maxduration":
		// direct mapping
	// Conditional required validators
	case "required_if", "required_unless", "required_with", "required_with_all", "required_without", "required_without_all":
		// direct mapping
	// Conditional excluded validators
	case "excluded_if", "excluded_unless", "excluded_with", "excluded_with_all", "excluded_without", "excluded_without_all":
		// direct mapping
	// Synonym mappings
	case "len":
		keyLower = "length"
	case "max":
		if isString {
			keyLower = "maxlength"
		} else if isCollection {
			keyLower = "maxitems"
		} else {
			// Fallback: treat as maxlength
			keyLower = "maxlength"
		}
	case "min":
		if isString {
			keyLower = "minlength"
		} else if isCollection {
			keyLower = "minitems"
		} else {
			keyLower = "minlength"
		}
	default:
		// Unknown token; ignore.
		return "", nil
	}

	identifier := "govalid:" + keyLower
	if value == "" {
		return identifier, nil
	}

	expressions := map[string]string{identifier: value}
	return identifier, expressions
}

// extractMarker extracts the identifier and expressions from a marker content string.
// It returns the identifier and a map of expressions if applicable.
// If the content does not contain an identifier or expressions, it returns an empty string and nil.
func extractMarker(content string) (string, map[string]string) {
	if strings.Count(content, "=") == 0 {
		return content, nil
	}

	// Split on the first = only, allowing expressions to contain = characters
	splits := strings.SplitN(content, "=", 2)
	if len(splits) != 2 {
		return "", nil
	}

	expressions := map[string]string{}
	expressions[splits[0]] = splits[1]

	return splits[0], expressions
}
