package books

import (
	"smedi-sme-db/app/auth"
	"smedi-sme-db/app/contracts"
	"smedi-sme-db/app/services"
)

// BooksPageController handles the books page
type BooksPageController struct {
	*contracts.GenericPageController
	bookService *services.BookService
}

// NewBooksPageController creates a new books page controller
func NewBooksPageController() *BooksPageController {
	bookService := services.NewBookService()

	return &BooksPageController{
		GenericPageController: contracts.NewGenericPageController(contracts.GenericPageConfig{
			ResourceType:      "books",
			PageComponent:     "Books/Index",
			Service:           bookService,
			ServiceIdentifier: auth.ServiceBooks,
			StatsEnabled:      true,
			StatsBuilder: func(controller *contracts.GenericPageController) map[string]interface{} {
				stats, _ := bookService.GetBookStatistics()
				return stats
			},
		}),
		bookService: bookService,
	}
}
