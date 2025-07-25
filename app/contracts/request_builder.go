package contracts

import (
	"github.com/goravel/framework/contracts/http"
)

// ============================================================================
// Create Request Builder
// ============================================================================

// CreateRequestBuilder uses a step-by-step builder pattern for create requests
type CreateRequestBuilder struct {
	request *baseCreateRequest
}

// Step 1: Required - Set validation rules
type CreateRequestBuilderWithRules struct {
	builder *CreateRequestBuilder
}

// Step 2: Required - Set data transformer
type CreateRequestBuilderWithTransform struct {
	builder *CreateRequestBuilder
}

// Step 3: Optional configurations and build
type CreateRequestBuilderComplete struct {
	builder *CreateRequestBuilder
}

// baseCreateRequest implements the CreateRequestContract
type baseCreateRequest struct {
	rules              map[string]string
	messages           map[string]string
	attributes         map[string]string
	authorize          func(ctx http.Context) error
	prepareValidation  func(ctx http.Context) error
	passedValidation   func(ctx http.Context) error
	toCreateData       func() map[string]interface{}
	context            http.Context
}

// NewCreateRequestBuilder creates a new create request builder
func NewCreateRequestBuilder() *CreateRequestBuilder {
	return &CreateRequestBuilder{
		request: &baseCreateRequest{
			messages:   make(map[string]string),
			attributes: make(map[string]string),
		},
	}
}

// Step 1: Set validation rules (REQUIRED)
func (b *CreateRequestBuilder) WithRules(rules map[string]string) *CreateRequestBuilderWithRules {
	if len(rules) == 0 {
		panic("At least one validation rule must be specified")
	}
	b.request.rules = rules
	return &CreateRequestBuilderWithRules{builder: b}
}

// Step 2: Set data transformer (REQUIRED)
func (b *CreateRequestBuilderWithRules) WithDataTransformer(
	transform func() map[string]interface{},
) *CreateRequestBuilderComplete {
	if transform == nil {
		panic("Data transformer function cannot be nil")
	}
	b.builder.request.toCreateData = transform
	return &CreateRequestBuilderComplete{builder: b.builder}
}

// Optional configurations
func (b *CreateRequestBuilderComplete) WithMessages(messages map[string]string) *CreateRequestBuilderComplete {
	b.builder.request.messages = messages
	return b
}

func (b *CreateRequestBuilderComplete) WithAttributes(attributes map[string]string) *CreateRequestBuilderComplete {
	b.builder.request.attributes = attributes
	return b
}

func (b *CreateRequestBuilderComplete) WithAuthorization(
	authorize func(ctx http.Context) error,
) *CreateRequestBuilderComplete {
	b.builder.request.authorize = authorize
	return b
}

func (b *CreateRequestBuilderComplete) WithPrepareValidation(
	prepare func(ctx http.Context) error,
) *CreateRequestBuilderComplete {
	b.builder.request.prepareValidation = prepare
	return b
}

func (b *CreateRequestBuilderComplete) WithPassedValidation(
	passed func(ctx http.Context) error,
) *CreateRequestBuilderComplete {
	b.builder.request.passedValidation = passed
	return b
}

// Build ensures the request implements CreateRequestContract
func (b *CreateRequestBuilderComplete) Build() CreateRequestContract {
	return b.builder.request
}

// Implement CreateRequestContract methods
func (r *baseCreateRequest) Rules(ctx http.Context) map[string]string {
	r.context = ctx
	return r.rules
}

func (r *baseCreateRequest) Messages(ctx http.Context) map[string]string {
	return r.messages
}

func (r *baseCreateRequest) Attributes(ctx http.Context) map[string]string {
	return r.attributes
}

func (r *baseCreateRequest) Authorize(ctx http.Context) error {
	if r.authorize != nil {
		return r.authorize(ctx)
	}
	return nil
}

func (r *baseCreateRequest) PrepareForValidation(ctx http.Context) error {
	if r.prepareValidation != nil {
		return r.prepareValidation(ctx)
	}
	return nil
}

func (r *baseCreateRequest) PassedValidation(ctx http.Context) error {
	if r.passedValidation != nil {
		return r.passedValidation(ctx)
	}
	return nil
}

func (r *baseCreateRequest) ToCreateData() map[string]interface{} {
	if r.toCreateData != nil {
		return r.toCreateData()
	}
	return make(map[string]interface{})
}

// ============================================================================
// Update Request Builder
// ============================================================================

// UpdateRequestBuilder uses a step-by-step builder pattern for update requests
type UpdateRequestBuilder struct {
	request *baseUpdateRequest
}

// Step 1: Required - Set validation rules
type UpdateRequestBuilderWithRules struct {
	builder *UpdateRequestBuilder
}

// Step 2: Required - Set data transformer
type UpdateRequestBuilderWithTransform struct {
	builder *UpdateRequestBuilder
}

// Step 3: Optional configurations and build
type UpdateRequestBuilderComplete struct {
	builder *UpdateRequestBuilder
}

// baseUpdateRequest implements the UpdateRequestContract
type baseUpdateRequest struct {
	resourceID         interface{}
	rules              map[string]string
	messages           map[string]string
	attributes         map[string]string
	authorize          func(ctx http.Context) error
	prepareValidation  func(ctx http.Context) error
	passedValidation   func(ctx http.Context) error
	toUpdateData       func() map[string]interface{}
	context            http.Context
}

// NewUpdateRequestBuilder creates a new update request builder
func NewUpdateRequestBuilder(resourceID interface{}) *UpdateRequestBuilder {
	return &UpdateRequestBuilder{
		request: &baseUpdateRequest{
			resourceID: resourceID,
			messages:   make(map[string]string),
			attributes: make(map[string]string),
		},
	}
}

// Step 1: Set validation rules (REQUIRED)
func (b *UpdateRequestBuilder) WithRules(rules map[string]string) *UpdateRequestBuilderWithRules {
	if len(rules) == 0 {
		panic("At least one validation rule must be specified")
	}
	b.request.rules = rules
	return &UpdateRequestBuilderWithRules{builder: b}
}

// Step 2: Set data transformer (REQUIRED)
func (b *UpdateRequestBuilderWithRules) WithDataTransformer(
	transform func() map[string]interface{},
) *UpdateRequestBuilderComplete {
	if transform == nil {
		panic("Data transformer function cannot be nil")
	}
	b.builder.request.toUpdateData = transform
	return &UpdateRequestBuilderComplete{builder: b.builder}
}

// Optional configurations (same as create)
func (b *UpdateRequestBuilderComplete) WithMessages(messages map[string]string) *UpdateRequestBuilderComplete {
	b.builder.request.messages = messages
	return b
}

func (b *UpdateRequestBuilderComplete) WithAttributes(attributes map[string]string) *UpdateRequestBuilderComplete {
	b.builder.request.attributes = attributes
	return b
}

func (b *UpdateRequestBuilderComplete) WithAuthorization(
	authorize func(ctx http.Context) error,
) *UpdateRequestBuilderComplete {
	b.builder.request.authorize = authorize
	return b
}

func (b *UpdateRequestBuilderComplete) WithPrepareValidation(
	prepare func(ctx http.Context) error,
) *UpdateRequestBuilderComplete {
	b.builder.request.prepareValidation = prepare
	return b
}

func (b *UpdateRequestBuilderComplete) WithPassedValidation(
	passed func(ctx http.Context) error,
) *UpdateRequestBuilderComplete {
	b.builder.request.passedValidation = passed
	return b
}

// Build ensures the request implements UpdateRequestContract
func (b *UpdateRequestBuilderComplete) Build() UpdateRequestContract {
	return b.builder.request
}

// Implement UpdateRequestContract methods
func (r *baseUpdateRequest) Rules(ctx http.Context) map[string]string {
	r.context = ctx
	return r.rules
}

func (r *baseUpdateRequest) Messages(ctx http.Context) map[string]string {
	return r.messages
}

func (r *baseUpdateRequest) Attributes(ctx http.Context) map[string]string {
	return r.attributes
}

func (r *baseUpdateRequest) Authorize(ctx http.Context) error {
	if r.authorize != nil {
		return r.authorize(ctx)
	}
	return nil
}

func (r *baseUpdateRequest) PrepareForValidation(ctx http.Context) error {
	if r.prepareValidation != nil {
		return r.prepareValidation(ctx)
	}
	return nil
}

func (r *baseUpdateRequest) PassedValidation(ctx http.Context) error {
	if r.passedValidation != nil {
		return r.passedValidation(ctx)
	}
	return nil
}

func (r *baseUpdateRequest) ToUpdateData() map[string]interface{} {
	if r.toUpdateData != nil {
		return r.toUpdateData()
	}
	return make(map[string]interface{})
}

func (r *baseUpdateRequest) GetResourceID() interface{} {
	return r.resourceID
}

// ============================================================================
// Example usage that enforces compile-time checks:
// ============================================================================

// CreateRequest example:
//
// request := NewCreateRequestBuilder().
//     WithRules(map[string]string{                    // REQUIRED
//         "title": "required|string|max:255",
//         "author": "required|string|max:100",
//     }).
//     WithDataTransformer(func() map[string]interface{} { // REQUIRED
//         return map[string]interface{}{
//             "title": formData.Title,
//             "author": formData.Author,
//         }
//     }).
//     WithMessages(customMessages).                    // Optional
//     WithAuthorization(authFunc).                     // Optional
//     Build()

// UpdateRequest example:
//
// request := NewUpdateRequestBuilder(resourceID).
//     WithRules(map[string]string{                    // REQUIRED
//         "title": "string|max:255",
//         "author": "string|max:100",
//     }).
//     WithDataTransformer(func() map[string]interface{} { // REQUIRED
//         return map[string]interface{}{
//             "title": formData.Title,
//             "author": formData.Author,
//         }
//     }).
//     WithAuthorization(authFunc).                     // Optional
//     Build()