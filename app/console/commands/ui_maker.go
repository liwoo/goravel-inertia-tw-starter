package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

// UIMaker generates the entire UI hierarchy for a page
type UIMaker struct {
}

// Signature The name and signature of the console command.
func (receiver *UIMaker) Signature() string {
	return "make:ui"
}

// Description The console command description.
func (receiver *UIMaker) Description() string {
	return "Generate the entire UI hierarchy for a page (Index, sections, types)"
}

// Extend The console command extend.
func (receiver *UIMaker) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:     "page",
				Aliases:  []string{"p"},
				Usage:    "Page name (e.g., Lender, Book, User)",
				Required: true,
			},
			&command.StringFlag{
				Name:     "request",
				Aliases:  []string{"r"},
				Usage:    "Request name for field inference (e.g., Lender, Book)",
				Required: true,
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *UIMaker) Handle(ctx console.Context) error {
	pageName := ctx.Option("page")
	requestName := ctx.Option("request")

	if pageName == "" {
		return fmt.Errorf("page name is required. Use --page=PageName")
	}

	if requestName == "" {
		return fmt.Errorf("request name is required. Use --request=RequestName")
	}

	// Normalize names
	pageName = strings.Title(pageName)
	requestName = strings.Title(requestName)
	pluralPage := receiver.pluralize(strings.ToLower(pageName))

	ctx.Info(fmt.Sprintf("Generating UI hierarchy for %s...", pageName))

	// Infer fields from request struct
	fields, err := receiver.inferFieldsFromRequest(requestName)
	if err != nil {
		ctx.Warning(fmt.Sprintf("Could not infer fields from request: %v", err))
		ctx.Info("Using default fields...")
		fields = receiver.getDefaultFields()
	}

	// Create directory structure
	pageDir := filepath.Join("resources", "js", "Pages", pageName)
	sectionsDir := filepath.Join(pageDir, "sections")

	for _, dir := range []string{pageDir, sectionsDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	// Generate TypeScript type file
	typePath := filepath.Join("resources", "js", "types", fmt.Sprintf("%s.ts", strings.ToLower(pageName)))
	if err := os.WriteFile(typePath, []byte(receiver.generateTypeContent(pageName, pluralPage, fields)), 0644); err != nil {
		return fmt.Errorf("failed to write type file: %v", err)
	}
	ctx.Success(fmt.Sprintf("✓ Generated type file: %s", typePath))

	// Generate Index.tsx
	indexPath := filepath.Join(pageDir, "Index.tsx")
	if err := os.WriteFile(indexPath, []byte(receiver.generateIndexContent(pageName, pluralPage, fields)), 0644); err != nil {
		return fmt.Errorf("failed to write Index file: %v", err)
	}
	ctx.Success(fmt.Sprintf("✓ Generated Index file: %s", indexPath))

	// Generate sections files
	files := map[string]func(string, string, []Field) string{
		fmt.Sprintf("%sColumns.tsx", pageName):    receiver.generateColumnsContent,
		fmt.Sprintf("%sCreateForm.tsx", pageName): receiver.generateCreateFormContent,
		fmt.Sprintf("%sEditForm.tsx", pageName):   receiver.generateEditFormContent,
		fmt.Sprintf("%sDetailView.tsx", pageName): receiver.generateDetailViewContent,
		fmt.Sprintf("%sPageConfig.tsx", pageName): receiver.generatePageConfigContent,
		"index.ts": receiver.generateSectionsIndexContent,
	}

	for filename, generator := range files {
		filePath := filepath.Join(sectionsDir, filename)
		content := generator(pageName, pluralPage, fields)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %v", filename, err)
		}
		ctx.Success(fmt.Sprintf("✓ Generated section file: %s", filePath))
	}

	// Print next steps
	ctx.Info("\n" + strings.Repeat("=", 60))
	ctx.Info("✅ UI files generated successfully!")
	ctx.Info(strings.Repeat("=", 60))
	ctx.Info("\n⚠️  IMPORTANT: To prevent 404 errors, complete these steps:")
	ctx.Info(strings.Repeat("=", 60))

	ctx.Info("\n1. Generate the page controller (REQUIRED):")
	ctx.Info(fmt.Sprintf("   go run . artisan make:page-ctrl --controller=%s", pageName))

	ctx.Info("\n2. Register the route in routes/web.go:")
	ctx.Info(fmt.Sprintf(`   %sPageController := %s.New%sPageController()`, pluralPage, pluralPage, pageName))
	ctx.Info(fmt.Sprintf(`   router.Get("/admin/%s", %sPageController.Index)`, pluralPage, pluralPage))

	ctx.Info("\n3. Add navigation menu item in resources/js/config/navigation.ts:")
	ctx.Info(fmt.Sprintf(`   {
       title: "%s",
       url: "/admin/%s",
       icon: YourIcon, // Import from lucide-react
       requiredService: "%s",
       requiredAction: "read" as const,
   }`, pageName, pluralPage, pluralPage))

	ctx.Info("\n4. Ensure your backend service implements:")
	ctx.Info(fmt.Sprintf("   - Get%sList(request) for paginated data", pageName))
	ctx.Info(fmt.Sprintf("   - Get%sStatistics() for stats (if needed)", pageName))

	ctx.Info("\n5. Review and customize the generated files:")
	ctx.Info(fmt.Sprintf("   - %s (main page)", indexPath))
	ctx.Info(fmt.Sprintf("   - %s (TypeScript types)", typePath))
	ctx.Info(fmt.Sprintf("   - %s/* (sections)", sectionsDir))

	ctx.Info("\n" + strings.Repeat("=", 60))

	return nil
}

type Field struct {
	Name         string
	Type         string
	TSType       string
	Label        string
	Required     bool
	IsPointer    bool
	IsArray      bool
	Placeholder  string
	Icon         string
	FormControl  string // input, textarea, select, date, etc.
	Options      []string
	Sortable     bool
	FilterType   string // text, select, number, date
	ShowInTable  bool
	ShowInDetail bool
}

func (receiver *UIMaker) pluralize(word string) string {
	if strings.HasSuffix(word, "y") && len(word) > 1 && !receiver.isVowel(rune(word[len(word)-2])) {
		return word[:len(word)-1] + "ies"
	}
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") || strings.HasSuffix(word, "z") {
		return word + "es"
	}
	return word + "s"
}

func (receiver *UIMaker) isVowel(r rune) bool {
	vowels := "aeiouAEIOU"
	return strings.ContainsRune(vowels, r)
}

func (receiver *UIMaker) getDefaultFields() []Field {
	return []Field{
		{
			Name: "Name", Type: "string", TSType: "string", Label: "Name", Required: true,
			Placeholder: "Enter name", Icon: "User", FormControl: "input", Sortable: true,
			FilterType: "text", ShowInTable: true, ShowInDetail: true,
		},
		{
			Name: "Email", Type: "string", TSType: "string", Label: "Email", Required: true,
			Placeholder: "Enter email", Icon: "Mail", FormControl: "input", Sortable: true,
			FilterType: "text", ShowInTable: true, ShowInDetail: true,
		},
	}
}

// inferFieldsFromRequest reads the request struct and infers fields
func (receiver *UIMaker) inferFieldsFromRequest(requestName string) ([]Field, error) {
	// Try to read the create request file
	createRequestPath := filepath.Join("app", "http", "requests", fmt.Sprintf("%s_create_request.go", strings.ToLower(requestName)))

	content, err := os.ReadFile(createRequestPath)
	if err != nil {
		// Try alternative naming
		createRequestPath = filepath.Join("app", "http", "requests", fmt.Sprintf("%s_request.go", strings.ToLower(requestName)))
		content, err = os.ReadFile(createRequestPath)
		if err != nil {
			return nil, err
		}
	}

	return receiver.parseRequestFields(string(content)), nil
}

func (receiver *UIMaker) parseRequestFields(content string) []Field {
	fields := []Field{}
	lines := strings.Split(content, "\n")

	inStruct := false
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Detect struct definition
		if strings.Contains(line, "CreateRequest struct") || strings.Contains(line, "Request struct") {
			inStruct = true
			continue
		}

		if inStruct && line == "}" {
			break
		}

		if inStruct && strings.Contains(line, "`") {
			// Parse field line: Name string `form:"name" json:"name"`
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				fieldName := parts[0]
				fieldType := parts[1]

				// Skip if starts with lowercase (embedded structs)
				if len(fieldName) > 0 && fieldName[0] >= 'a' && fieldName[0] <= 'z' {
					continue
				}

				field := receiver.createFieldFromType(fieldName, fieldType)
				fields = append(fields, field)
			}
		}
	}

	return fields
}

func (receiver *UIMaker) createFieldFromType(name, goType string) Field {
	field := Field{
		Name:         name,
		Type:         goType,
		Label:        receiver.humanize(name),
		Placeholder:  "Enter " + strings.ToLower(receiver.humanize(name)),
		ShowInTable:  true,
		ShowInDetail: true,
		Sortable:     true,
	}

	// Handle pointers
	if strings.HasPrefix(goType, "*") {
		field.IsPointer = true
		field.Required = false
		goType = strings.TrimPrefix(goType, "*")
		field.Type = goType
	} else {
		field.Required = true
	}

	// Handle arrays
	if strings.HasPrefix(goType, "[]") {
		field.IsArray = true
		goType = strings.TrimPrefix(goType, "[]")
		field.Type = goType
	}

	// Map Go types to TypeScript and form controls
	switch goType {
	case "string":
		field.TSType = "string"
		field.FormControl = "input"
		field.FilterType = "text"
		field.Icon = receiver.getIconForField(name)
		// Special cases
		if strings.Contains(strings.ToLower(name), "email") {
			field.Icon = "Mail"
			field.Placeholder = "Enter email address"
		} else if strings.Contains(strings.ToLower(name), "phone") {
			field.Icon = "Phone"
			field.Placeholder = "Enter phone number"
		} else if strings.Contains(strings.ToLower(name), "address") {
			field.Icon = "MapPin"
			field.FormControl = "textarea"
			field.Placeholder = "Enter address"
		} else if strings.Contains(strings.ToLower(name), "description") {
			field.Icon = "FileText"
			field.FormControl = "textarea"
		}
	case "int", "int64", "uint", "uint64":
		field.TSType = "number"
		field.FormControl = "input"
		field.FilterType = "number"
		field.Icon = "Hash"
	case "float64", "float32":
		field.TSType = "number"
		field.FormControl = "input"
		field.FilterType = "number"
		field.Icon = "DollarSign"
		if strings.Contains(strings.ToLower(name), "price") || strings.Contains(strings.ToLower(name), "amount") {
			field.Placeholder = "0.00"
		}
	case "bool":
		field.TSType = "boolean"
		field.FormControl = "checkbox"
		field.FilterType = "select"
		field.Icon = "CheckSquare"
	case "time.Time", "Time":
		field.TSType = "string"
		field.FormControl = "date"
		field.FilterType = "date"
		field.Icon = "Calendar"
	default:
		// Check for custom types (enums, status, etc.)
		if strings.Contains(strings.ToLower(name), "status") || strings.Contains(strings.ToLower(name), "gender") {
			field.TSType = "string"
			field.FormControl = "select"
			field.FilterType = "select"
			field.Icon = "Tag"
		} else {
			field.TSType = "string"
			field.FormControl = "input"
			field.FilterType = "text"
			field.Icon = "FileText"
		}
	}

	return field
}

func (receiver *UIMaker) getIconForField(name string) string {
	nameLower := strings.ToLower(name)
	iconMap := map[string]string{
		"name":        "User",
		"title":       "BookOpen",
		"email":       "Mail",
		"phone":       "Phone",
		"address":     "MapPin",
		"description": "FileText",
		"price":       "DollarSign",
		"amount":      "DollarSign",
		"date":        "Calendar",
		"time":        "Clock",
		"status":      "Tag",
		"tag":         "Tag",
		"category":    "FolderOpen",
		"type":        "Tag",
		"gender":      "User",
	}

	for key, icon := range iconMap {
		if strings.Contains(nameLower, key) {
			return icon
		}
	}

	return "FileText"
}

func (receiver *UIMaker) humanize(name string) string {
	// Convert camelCase or PascalCase to Title Case
	var result strings.Builder
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune(' ')
		}
		result.WriteRune(r)
	}
	return result.String()
}
