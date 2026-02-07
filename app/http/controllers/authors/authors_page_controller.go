package authors

import (
	"books-database/app/auth"
	"books-database/app/contracts"
	"books-database/app/services"
)

// AuthorsPageController handles the authors page
type AuthorsPageController struct {
	*contracts.GenericPageController
	authorService *services.AuthorService
}

// NewAuthorsPageController creates a new authors page controller
func NewAuthorsPageController() *AuthorsPageController {
	authorService := services.NewAuthorService()

	return &AuthorsPageController{
		GenericPageController: contracts.NewGenericPageController(contracts.GenericPageConfig{
			ResourceType:      "authors",
			PageComponent:     "Authors/Index",
			Service:           authorService,
			ServiceIdentifier: auth.ServiceAuthors,
			StatsEnabled:      true,
			StatsBuilder: func(controller *contracts.GenericPageController) map[string]interface{} {
				stats, _ := authorService.GetAuthorStatistics()
				return stats
			},
		}),
		authorService: authorService,
	}
}
