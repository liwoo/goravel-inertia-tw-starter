package contracts

import (
	"github.com/goravel/framework/contracts/http"
	"testing"
)

// Test types that implement the required interfaces
type TestModel struct{}

type TestCreateRequest struct{}

func (t *TestCreateRequest) Rules(ctx http.Context) map[string]string      { return nil }
func (t *TestCreateRequest) Messages(ctx http.Context) map[string]string   { return nil }
func (t *TestCreateRequest) Attributes(ctx http.Context) map[string]string { return nil }
func (t *TestCreateRequest) Authorize(ctx http.Context) error              { return nil }
func (t *TestCreateRequest) PrepareForValidation(ctx http.Context) error   { return nil }
func (t *TestCreateRequest) PassedValidation(ctx http.Context) error       { return nil }
func (t *TestCreateRequest) ToCreateData() map[string]interface{}          { return nil }

type TestUpdateRequest struct{}

func (t *TestUpdateRequest) Rules(ctx http.Context) map[string]string      { return nil }
func (t *TestUpdateRequest) Messages(ctx http.Context) map[string]string   { return nil }
func (t *TestUpdateRequest) Attributes(ctx http.Context) map[string]string { return nil }
func (t *TestUpdateRequest) Authorize(ctx http.Context) error              { return nil }
func (t *TestUpdateRequest) PrepareForValidation(ctx http.Context) error   { return nil }
func (t *TestUpdateRequest) PassedValidation(ctx http.Context) error       { return nil }
func (t *TestUpdateRequest) ToUpdateData() map[string]interface{}          { return nil }
func (t *TestUpdateRequest) GetResourceID() interface{}                    { return nil }

func TestStaticBuilderCompileTimeCheck(t *testing.T) {
	// This function should NOT compile:
	builder := NewStaticControllerBuilder[TestModel, *TestCreateRequest, *TestUpdateRequest](
		"test",
		nil,
	)
	// Trying to call Build() without going through all required steps
	_ = builder.Build() // This line should cause a compile error

	// This SHOULD compile:
	completeBuilder := NewStaticControllerBuilder[TestModel, *TestCreateRequest, *TestUpdateRequest](
		"test",
		nil,
	).
		ValidateCreateRequest().
		ValidateUpdateRequest().
		WithAuthChecker(nil)

	// Now Build() should be available
	_ = completeBuilder.Build()
}
