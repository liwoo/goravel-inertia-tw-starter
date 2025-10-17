package requests

import (
	"github.com/goravel/framework/contracts/http"
	"players/app/models"
)

// MessageCreateRequest handles message creation validation
type MessageCreateRequest struct {
	RecipientID uint               `json:"recipient_id" form:"recipient_id"`
	Content     string             `json:"content" form:"content"`
	Type        models.MessageType `json:"type" form:"type"`
}

// Authorize determines if the user can create a message
func (r *MessageCreateRequest) Authorize(ctx http.Context) error {
	return nil // Authentication handled by middleware
}

// Rules returns validation rules
func (r *MessageCreateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"recipient_id": "required|uint",
		"content":      "required|max:5000",
		"type":         "in:direct,system,broadcast",
	}
}

// Messages returns custom validation messages
func (r *MessageCreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"recipient_id.required": "Recipient is required",
		"content.required":      "Message content is required",
		"content.max":           "Message content cannot exceed 5000 characters",
	}
}

// Attributes returns custom attribute names
func (r *MessageCreateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"recipient_id": "recipient",
		"content":      "message content",
	}
}

// PrepareForValidation prepares the data before validation
func (r *MessageCreateRequest) PrepareForValidation(ctx http.Context) error {
	// Set default message type if not provided
	if r.Type == "" {
		r.Type = models.MessageTypeDirect
	}
	return nil
}

// PassedValidation is called after validation passes
func (r *MessageCreateRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// Filters returns validation filters
func (r *MessageCreateRequest) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

// ToCreateData converts the request to create data
func (r *MessageCreateRequest) ToCreateData() map[string]interface{} {
	return map[string]interface{}{
		"recipient_id": r.RecipientID,
		"content":      r.Content,
		"type":         r.Type,
	}
}

// MessageUpdateRequest handles message update validation
type MessageUpdateRequest struct {
	ResourceID uint    `json:"-" form:"-"`
	Content    *string `json:"content" form:"content"`
	IsRead     *bool   `json:"is_read" form:"is_read"`
	Important  *bool   `json:"important" form:"important"`
}

// Authorize determines if the user can update a message
func (r *MessageUpdateRequest) Authorize(ctx http.Context) error {
	return nil // Authorization handled in controller
}

// Rules returns validation rules
func (r *MessageUpdateRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"content": "max:5000",
	}
}

// Messages returns custom validation messages
func (r *MessageUpdateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{
		"content.max": "Message content cannot exceed 5000 characters",
	}
}

// Attributes returns custom attribute names
func (r *MessageUpdateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"content": "message content",
	}
}

// PrepareForValidation prepares the data before validation
func (r *MessageUpdateRequest) PrepareForValidation(ctx http.Context) error {
	return nil
}

// PassedValidation is called after validation passes
func (r *MessageUpdateRequest) PassedValidation(ctx http.Context) error {
	return nil
}

// Filters returns validation filters
func (r *MessageUpdateRequest) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

// ToUpdateData converts the request to update data
func (r *MessageUpdateRequest) ToUpdateData() map[string]interface{} {
	data := make(map[string]interface{})
	if r.Content != nil {
		data["content"] = *r.Content
	}
	if r.IsRead != nil {
		data["is_read"] = *r.IsRead
	}
	if r.Important != nil {
		data["important"] = *r.Important
	}
	return data
}

// GetResourceID returns the resource ID for this update request
func (r *MessageUpdateRequest) GetResourceID() interface{} {
	return r.ResourceID
}
