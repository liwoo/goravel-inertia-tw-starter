package books

import (
	"github.com/goravel/framework/contracts/http"
	"players/app/auth"
	"players/app/contracts"
	"players/app/http/requests"
	"players/app/models"
	"players/app/services"
)

// BrokenController1 - missing ValidateCreateRequest - THIS SHOULD NOT COMPILE
func BrokenController1() {
	bookService := services.NewBookService()
	
	// This should cause a compile error
	staticController := contracts.NewStaticControllerBuilder[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest](
		"book",
		bookService,
	).
		// ValidateCreateRequest().  // MISSING!
		ValidateUpdateRequest().     // This line should error
		WithAuthChecker(func(ctx http.Context, action string, resource interface{}) error {
			return nil
		}).
		Build()
		
	_ = staticController
}

// BrokenController2 - missing WithAuthChecker - THIS SHOULD NOT COMPILE  
func BrokenController2() {
	bookService := services.NewBookService()
	
	// This should cause a compile error
	staticController := contracts.NewStaticControllerBuilder[models.Book, *requests.BookCreateRequest, *requests.BookUpdateRequest](
		"book",
		bookService,
	).
		ValidateCreateRequest().
		ValidateUpdateRequest().
		// WithAuthChecker(...).  // MISSING!
		Build()  // This line should error
		
	_ = staticController
}