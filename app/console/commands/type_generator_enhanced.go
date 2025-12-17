package commands

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

// EnhancedTypeGenerator generates TypeScript types from Go request structs and enums
type EnhancedTypeGenerator struct{}

// RequestField represents a parsed field from a Go request struct
type RequestField struct {
	Name       string
	GoType     string
	TSType     string
	IsPointer  bool
	IsArray    bool
	JSONTag    string
	IsRequired bool
}

// GoEnum represents a Go enum type
type GoEnum struct {
	TypeName string
	BaseType string
	Values   []GoEnumValue
}

// GoEnumValue represents a single enum constant
type GoEnumValue struct {
	ConstName string
	Value     string
	Label     string
}

// ParseRequestStruct extracts field information from a Go request struct file
func (g *EnhancedTypeGenerator) ParseRequestStruct(filePath string) ([]RequestField, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	var fields []RequestField

	ast.Inspect(node, func(n ast.Node) bool {
		// Look for struct types
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		// Only process CreateRequest or Request structs
		if !strings.Contains(typeSpec.Name.Name, "Request") {
			return true
		}

		// Parse struct fields
		for _, field := range structType.Fields.List {
			for _, name := range field.Names {
				// Skip unexported fields
				if !ast.IsExported(name.Name) {
					continue
				}

				reqField := g.parseField(name.Name, field)
				fields = append(fields, reqField)
			}
		}

		return true
	})

	return fields, nil
}

// parseField extracts field information from an AST field
func (g *EnhancedTypeGenerator) parseField(fieldName string, field *ast.Field) RequestField {
	reqField := RequestField{
		Name: fieldName,
	}

	// Extract JSON tag
	if field.Tag != nil {
		tag := field.Tag.Value
		reqField.JSONTag = g.extractJSONTag(tag)
	}

	// Parse the field type
	reqField.GoType, reqField.IsPointer, reqField.IsArray = g.parseType(field.Type)
	reqField.TSType = g.goTypeToTS(reqField.GoType, reqField.IsArray)
	reqField.IsRequired = !reqField.IsPointer

	return reqField
}

// parseType extracts type information from an AST expression
func (g *EnhancedTypeGenerator) parseType(expr ast.Expr) (typeName string, isPointer bool, isArray bool) {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name, false, false
	case *ast.StarExpr:
		innerType, _, innerArray := g.parseType(t.X)
		return innerType, true, innerArray
	case *ast.ArrayType:
		innerType, innerPointer, _ := g.parseType(t.Elt)
		return innerType, innerPointer, true
	case *ast.SelectorExpr:
		// Handle qualified types like time.Time
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name + "." + t.Sel.Name, false, false
		}
		return t.Sel.Name, false, false
	default:
		return "interface{}", false, false
	}
}

// goTypeToTS converts Go types to TypeScript types
func (g *EnhancedTypeGenerator) goTypeToTS(goType string, isArray bool) string {
	var baseType string

	switch goType {
	case "string":
		baseType = "string"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		baseType = "number"
	case "bool":
		baseType = "boolean"
	case "time.Time", "Time":
		baseType = "string" // ISO date string
	default:
		// Check if it's a custom type (might be an enum)
		if strings.HasSuffix(goType, "Type") || strings.HasSuffix(goType, "Status") {
			baseType = goType // Keep the custom type name
		} else {
			baseType = "any"
		}
	}

	if isArray {
		return baseType + "[]"
	}
	return baseType
}

// extractJSONTag extracts the json field name from struct tags
func (g *EnhancedTypeGenerator) extractJSONTag(tag string) string {
	// Remove backticks
	tag = strings.Trim(tag, "`")

	// Find json tag
	parts := strings.Split(tag, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, "json:") {
			jsonValue := strings.TrimPrefix(part, "json:")
			jsonValue = strings.Trim(jsonValue, `"`)
			// Take the first part before comma
			jsonParts := strings.Split(jsonValue, ",")
			return jsonParts[0]
		}
	}

	return ""
}

// ParseEnums extracts enum type definitions from a Go file
func (g *EnhancedTypeGenerator) ParseEnums(filePath string) ([]GoEnum, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	enumTypes := make(map[string]*GoEnum)
	var currentEnumType string

	// First pass: find type aliases (type GenderType string)
	ast.Inspect(node, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			return true
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// Check if it's a string-based type
			if ident, ok := typeSpec.Type.(*ast.Ident); ok && ident.Name == "string" {
				enumTypes[typeSpec.Name.Name] = &GoEnum{
					TypeName: typeSpec.Name.Name,
					BaseType: "string",
					Values:   []GoEnumValue{},
				}
			}
		}

		return true
	})

	// Second pass: find const declarations
	ast.Inspect(node, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			return true
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			// Check if the const has an explicit type
			if valueSpec.Type != nil {
				if ident, ok := valueSpec.Type.(*ast.Ident); ok {
					currentEnumType = ident.Name
				}
			}

			// If we're in an enum type, add the values
			if enum, exists := enumTypes[currentEnumType]; exists {
				for i, name := range valueSpec.Names {
					var value string
					if i < len(valueSpec.Values) {
						if lit, ok := valueSpec.Values[i].(*ast.BasicLit); ok {
							value = strings.Trim(lit.Value, `"`)
						}
					}

					label := g.generateEnumLabel(name.Name, currentEnumType)

					enum.Values = append(enum.Values, GoEnumValue{
						ConstName: name.Name,
						Value:     value,
						Label:     label,
					})
				}
			}
		}

		return true
	})

	// Convert map to slice, only include enums with values
	var result []GoEnum
	for _, enum := range enumTypes {
		if len(enum.Values) > 0 {
			result = append(result, *enum)
		}
	}

	return result, nil
}

// generateEnumLabel creates a human-readable label from a constant name
func (g *EnhancedTypeGenerator) generateEnumLabel(constName, typeName string) string {
	// Remove type name prefix (e.g., GenderMale -> Male)
	label := strings.TrimPrefix(constName, typeName)

	// If nothing removed, try common prefixes
	if label == constName {
		prefixes := []string{"Type", "Status", "Kind", "Mode"}
		for _, prefix := range prefixes {
			if strings.HasPrefix(label, prefix) {
				label = strings.TrimPrefix(label, prefix)
				break
			}
		}
	}

	// Add spaces before capital letters
	var result strings.Builder
	for i, r := range label {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune(' ')
		}
		result.WriteRune(r)
	}

	return strings.TrimSpace(result.String())
}

// GenerateTypeScriptInterface generates a TypeScript interface from request fields
func (g *EnhancedTypeGenerator) GenerateTypeScriptInterface(structName string, fields []RequestField) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("export interface %s {\n", structName))

	for _, field := range fields {
		jsonName := field.JSONTag
		if jsonName == "" {
			jsonName = g.toJSONName(field.Name)
		}

		optional := ""
		if !field.IsRequired {
			optional = "?"
		}

		sb.WriteString(fmt.Sprintf("  %s%s: %s;\n", jsonName, optional, field.TSType))
	}

	sb.WriteString("}\n")

	return sb.String()
}

// toJSONName converts PascalCase to camelCase
func (g *EnhancedTypeGenerator) toJSONName(name string) string {
	if len(name) == 0 {
		return name
	}
	return strings.ToLower(name[:1]) + name[1:]
}

// GenerateEnumTypeScript generates TypeScript union types and option arrays for enums
func (g *EnhancedTypeGenerator) GenerateEnumTypeScript(enum GoEnum) string {
	var sb strings.Builder

	// Generate union type
	sb.WriteString(fmt.Sprintf("export type %s = ", enum.TypeName))
	for i, val := range enum.Values {
		if i > 0 {
			sb.WriteString(" | ")
		}
		sb.WriteString(fmt.Sprintf("'%s'", val.Value))
	}
	sb.WriteString(";\n\n")

	// Generate options array
	optionsName := strings.ToUpper(g.toSnakeCase(enum.TypeName)) + "_OPTIONS"
	sb.WriteString(fmt.Sprintf("export const %s: { value: %s; label: string }[] = [\n",
		optionsName, enum.TypeName))

	for i, val := range enum.Values {
		sb.WriteString(fmt.Sprintf("  { value: '%s', label: '%s' }", val.Value, val.Label))
		if i < len(enum.Values)-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteString("\n")
		}
	}
	sb.WriteString("];\n")

	return sb.String()
}

// toSnakeCase converts PascalCase to snake_case
func (g *EnhancedTypeGenerator) toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// ScanDirectoryForEnums scans a directory for all enum type definitions
func (g *EnhancedTypeGenerator) ScanDirectoryForEnums(dirPath string) ([]GoEnum, error) {
	var allEnums []GoEnum

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}

		filePath := dirPath + "/" + entry.Name()
		enums, err := g.ParseEnums(filePath)
		if err != nil {
			// Skip files that can't be parsed
			continue
		}

		allEnums = append(allEnums, enums...)
	}

	return allEnums, nil
}
