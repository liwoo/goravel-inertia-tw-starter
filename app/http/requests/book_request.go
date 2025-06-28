package requests

import (
	"fmt"
	"players/app/contracts"
	"strings"

	"github.com/goravel/framework/contracts/http"
)

// BookCreateRequest handles book creation validation
type BookCreateRequest struct {
	Title       string   `form:"title" json:"title"`
	Author      string   `form:"author" json:"author"`
	ISBN        string   `form:"isbn" json:"isbn"`
	Description string   `form:"description" json:"description"`
	Price       float64  `form:"price" json:"price"`
	Status      string   `form:"status" json:"status"`
	PublishedAt string   `form:"publishedAt" json:"publishedAt"`
	Tags        []string `form:"tags" json:"tags"`
}

// Rules defines validation rules for book creation
func (r *BookCreateRequest) Rules(ctx http.Context) map[string]string {
	rules := map[string]string{
		"title":       fmt.Sprintf("%s|%s", contracts.Required, fmt.Sprintf(contracts.MaxLength, 255)),
		"author":      fmt.Sprintf("%s|%s", contracts.Required, fmt.Sprintf(contracts.MaxLength, 100)),
		"isbn":        fmt.Sprintf("%s|%s", contracts.Required, fmt.Sprintf(contracts.Regex, "^[0-9-]{10,17}$")),
		"description": fmt.Sprintf(contracts.MaxLength, 1000),
		"price":       fmt.Sprintf("%s|%s|%s", contracts.Required, contracts.Numeric, fmt.Sprintf(contracts.MinValue, 0)),
		"status":      fmt.Sprintf("in:%s", "AVAILABLE,BORROWED,MAINTENANCE"),
		"publishedAt": contracts.Date,
		"tags":        fmt.Sprintf("%s|%s", contracts.Array, fmt.Sprintf(contracts.ArrayMax, 10)),
		"tags.*":      fmt.Sprintf(contracts.MaxLength, 50),
	}
	
	return rules
}

// Messages defines custom validation messages
func (r *BookCreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"title.required":      "Book title is required",
		"title.max":           "Book title cannot exceed 255 characters",
		"author.required":     "Author name is required",
		"author.max":          "Author name cannot exceed 100 characters",
		"isbn.required":       "ISBN is required",
		"isbn.regex":          "ISBN must be 10-13 digits",
		"isbn.unique":         "This ISBN already exists",
		"price.required":      "Price is required",
		"price.numeric":       "Price must be a valid number",
		"price.min":           "Price must be greater than or equal to 0",
		"status.in":           "Status must be one of: AVAILABLE, BORROWED, MAINTENANCE",
		"publishedAt.date":    "Published date must be a valid date",
		"publishedAt.before":  "Published date cannot be in the future",
		"tags.array":          "Tags must be an array",
		"tags.max":            "Maximum 10 tags allowed",
		"tags.*.max":          "Each tag cannot exceed 50 characters",
	}
}

// Attributes defines custom attribute names
func (r *BookCreateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"publishedAt": "publication date",
		"isbn":        "ISBN number",
	}
}

// Authorize determines if the user is authorized to make this request
func (r *BookCreateRequest) Authorize(ctx http.Context) error {
	// This would typically check permissions via Gates
	// For now, assume all authenticated users can create books
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *BookCreateRequest) PrepareForValidation(ctx http.Context) error {
	// Example: Normalize ISBN by removing hyphens
	if r.ISBN != "" {
		// Remove common ISBN separators
		r.ISBN = strings.ReplaceAll(r.ISBN, "-", "")
		r.ISBN = strings.ReplaceAll(r.ISBN, " ", "")
	}

	// Set default status if not provided
	if r.Status == "" {
		r.Status = "AVAILABLE"
	}

	return nil
}

// PassedValidation is called after validation passes
func (r *BookCreateRequest) PassedValidation(ctx http.Context) error {
	// Could log the validation success, perform additional checks, etc.
	return nil
}

// ToCreateData converts the request to create data map
func (r *BookCreateRequest) ToCreateData() map[string]interface{} {
	data := map[string]interface{}{
		"title":       r.Title,
		"author":      r.Author,
		"isbn":        r.ISBN,
		"description": r.Description,
		"price":       r.Price,
		"status":      r.Status,
	}

	// Only include optional fields if they have values
	if r.PublishedAt != "" {
		data["publishedAt"] = r.PublishedAt
	}

	if len(r.Tags) > 0 {
		data["tags"] = r.Tags
	}

	return data
}

// BookUpdateRequest handles book update validation
type BookUpdateRequest struct {
	Title       *string   `form:"title" json:"title"`
	Author      *string   `form:"author" json:"author"`
	ISBN        *string   `form:"isbn" json:"isbn"`
	Description *string   `form:"description" json:"description"`
	Price       *float64  `form:"price" json:"price"`
	Status      *string   `form:"status" json:"status"`
	PublishedAt *string   `form:"publishedAt" json:"publishedAt"`
	Tags        *[]string `form:"tags" json:"tags"`
	ID          uint      `form:"-" json:"-"` // Set by controller
}

// Rules defines validation rules for book updates
func (r *BookUpdateRequest) Rules(ctx http.Context) map[string]string {
	rules := map[string]string{}

	// Get the book ID from the route parameter for unique validation
	bookID := ctx.Request().Route("id")

	// Only validate fields that are provided
	if r.Title != nil {
		rules["title"] = fmt.Sprintf(contracts.MaxLength, 255)
	}
	if r.Author != nil {
		rules["author"] = fmt.Sprintf(contracts.MaxLength, 100)
	}
	if r.ISBN != nil {
		// Fix unique validation to exclude current record
		rules["isbn"] = fmt.Sprintf("%s|unique:books,isbn,%s", fmt.Sprintf(contracts.Regex, "^[0-9]{10,13}$"), bookID)
	}
	if r.Description != nil {
		rules["description"] = fmt.Sprintf(contracts.MaxLength, 1000)
	}
	if r.Price != nil {
		rules["price"] = fmt.Sprintf("%s|%s", contracts.Numeric, fmt.Sprintf(contracts.MinValue, 0))
	}
	if r.Status != nil {
		rules["status"] = "in:AVAILABLE,BORROWED,MAINTENANCE"
	}
	if r.PublishedAt != nil {
		rules["publishedAt"] = fmt.Sprintf("%s|%s", contracts.Date, fmt.Sprintf(contracts.Before, "today"))
	}
	if r.Tags != nil {
		rules["tags"] = fmt.Sprintf("%s|%s", contracts.Array, fmt.Sprintf(contracts.ArrayMax, 10))
		rules["tags.*"] = fmt.Sprintf(contracts.MaxLength, 50)
	}

	// If no rules were added, add a dummy rule to prevent empty rules error
	if len(rules) == 0 {
		rules["_at_least_one_field"] = "sometimes"
	}

	return rules
}

// Messages defines custom validation messages for updates
func (r *BookUpdateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"title.max":           "Book title cannot exceed 255 characters",
		"author.max":          "Author name cannot exceed 100 characters",
		"isbn.regex":          "ISBN must be 10-13 digits",
		"isbn.unique":         "This ISBN already exists",
		"price.numeric":       "Price must be a valid number",
		"price.min":           "Price must be greater than or equal to 0",
		"status.in":           "Status must be one of: AVAILABLE, BORROWED, MAINTENANCE",
		"publishedAt.date":    "Published date must be a valid date",
		"publishedAt.before":  "Published date cannot be in the future",
		"tags.array":          "Tags must be an array",
		"tags.max":            "Maximum 10 tags allowed",
		"tags.*.max":          "Each tag cannot exceed 50 characters",
	}
}

// Attributes defines custom attribute names for updates
func (r *BookUpdateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"publishedAt": "publication date",
		"isbn":        "ISBN number",
	}
}

// Authorize determines if the user is authorized to update this book
func (r *BookUpdateRequest) Authorize(ctx http.Context) error {
	// This would check if user can update this specific book
	// Could use Gates: facades.Gate().Allows("update.books", book)
	return nil
}

// PrepareForValidation allows modification of input before validation
func (r *BookUpdateRequest) PrepareForValidation(ctx http.Context) error {
	// Normalize ISBN if provided
	if r.ISBN != nil && *r.ISBN != "" {
		normalized := strings.ReplaceAll(*r.ISBN, "-", "")
		normalized = strings.ReplaceAll(normalized, " ", "")
		r.ISBN = &normalized
	}

	return nil
}

// PassedValidation is called after validation passes
func (r *BookUpdateRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// ToUpdateData converts the request to update data map
func (r *BookUpdateRequest) ToUpdateData() map[string]interface{} {
	data := map[string]interface{}{}

	// Only include fields that are provided (not nil)
	if r.Title != nil {
		data["title"] = *r.Title
	}
	if r.Author != nil {
		data["author"] = *r.Author
	}
	if r.ISBN != nil {
		data["isbn"] = *r.ISBN
	}
	if r.Description != nil {
		data["description"] = *r.Description
	}
	if r.Price != nil {
		data["price"] = *r.Price
	}
	if r.Status != nil {
		data["status"] = *r.Status
	}
	if r.PublishedAt != nil {
		data["publishedAt"] = *r.PublishedAt
	}
	if r.Tags != nil {
		data["tags"] = *r.Tags
	}

	return data
}

// GetResourceID returns the resource ID for update
func (r *BookUpdateRequest) GetResourceID() interface{} {
	return r.ID
}