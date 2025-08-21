package tests

import (
	"github.com/goravel/framework/contracts/http"
	"players/app/contracts"
	"players/app/http/requests"
	"players/app/models"
	"players/app/services"
	"testing"
)

// This file demonstrates compile-time enforcement
// Uncomment any of the error examples to see compile errors in your editor

func TestCompileTimeEnforcement(t *testing.T) {
	bookService := services.NewBookService()

	// ✅ CORRECT: This compiles
	_ = contracts.NewStaticControllerBuilder[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest](
		"book",
		bookService,
	).
		ValidateCreateRequest().
		ValidateUpdateRequest().
		WithAuthChecker(func(ctx http.Context, action string, resource interface{}) error {
			return nil
		}).
		Build()

	// ❌ ERROR 1: Comment out ValidateCreateRequest()
	// Uncomment below to see error in your editor:
	/*
		_ = contracts.NewStaticControllerBuilder[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest](
			"book",
			bookService,
		).
			// ValidateCreateRequest().  // COMMENTED OUT
			ValidateUpdateRequest().     // ERROR: ValidateUpdateRequest undefined (type *contracts.StaticControllerBuilder[...needsCreateRequest] has no field or method ValidateUpdateRequest)
			WithAuthChecker(func(ctx http.Context, action string, resource interface{}) error {
				return nil
			}).
			Build()
	*/

	// ❌ ERROR 2: Comment out ValidateUpdateRequest()
	// Uncomment below to see error in your editor:
	/*
		_ = contracts.NewStaticControllerBuilder[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest](
			"book",
			bookService,
		).
			ValidateCreateRequest().
			// ValidateUpdateRequest().  // COMMENTED OUT
			WithAuthChecker(func(ctx http.Context, action string, resource interface{}) error {
				return nil
			}).  // ERROR: WithAuthChecker undefined (type *contracts.StaticControllerBuilder[...needsUpdateRequest] has no field or method WithAuthChecker)
			Build()
	*/

	// ❌ ERROR 3: Comment out WithAuthChecker()
	// Uncomment below to see error in your editor:
	/*
		_ = contracts.NewStaticControllerBuilder[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest](
			"book",
			bookService,
		).
			ValidateCreateRequest().
			ValidateUpdateRequest().
			// WithAuthChecker(...).  // COMMENTED OUT
			Build()  // ERROR: Build undefined (type *contracts.StaticControllerBuilder[...needsAuthChecker] has no field or method Build)
	*/

	// ❌ ERROR 4: Wrong order
	// Uncomment below to see error in your editor:
	/*
		_ = contracts.NewStaticControllerBuilder[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest](
			"book",
			bookService,
		).
			WithAuthChecker(func(ctx http.Context, action string, resource interface{}) error {
				return nil
			}).  // ERROR: WithAuthChecker undefined (type *contracts.StaticControllerBuilder[...needsCreateRequest] has no field or method WithAuthChecker)
			ValidateCreateRequest().
			ValidateUpdateRequest().
			Build()
	*/

	// ❌ ERROR 5: Try to Build() immediately
	// Uncomment below to see error in your editor:
	/*
		_ = contracts.NewStaticControllerBuilder[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest](
			"book",
			bookService,
		).
			Build()  // ERROR: Build undefined (type *contracts.StaticControllerBuilder[...needsCreateRequest] has no field or method Build)
	*/
}

// This demonstrates the compile-time errors you'll see:
//
// 1. If you comment out ValidateCreateRequest():
//    "ValidateUpdateRequest undefined (type *StaticControllerBuilder[...needsCreateRequest] has no field or method ValidateUpdateRequest)"
//
// 2. If you comment out ValidateUpdateRequest():
//    "WithAuthChecker undefined (type *StaticControllerBuilder[...needsUpdateRequest] has no field or method WithAuthChecker)"
//
// 3. If you comment out WithAuthChecker():
//    "Build undefined (type *StaticControllerBuilder[...needsAuthChecker] has no field or method Build)"
//
// The type system enforces that ALL methods must be called in the correct order!
