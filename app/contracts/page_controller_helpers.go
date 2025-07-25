package contracts

// BuildStandardPageProps is a helper function that builds the standard PageProps structure
// to reduce boilerplate in page controllers
func BuildStandardPageProps(
	result *PaginatedResult,
	req *ListRequest,
	permissions map[string]bool,
	stats map[string]interface{},
	component string,
	resourceType string,
) PageProps {
	return PageProps{
		Data: PaginatedData{
			Data:        result.Data,
			Total:       result.Total,
			CurrentPage: result.CurrentPage,
			LastPage:    result.LastPage,
			PerPage:     result.PerPage,
			From:        (result.CurrentPage - 1) * result.PerPage + 1,
			To:          min(result.CurrentPage*result.PerPage, int(result.Total)),
			HasNext:     result.CurrentPage < result.LastPage,
			HasPrev:     result.CurrentPage > 1,
		},
		Filters: FilterState{
			Page:      req.Page,
			PageSize:  req.PageSize,
			Search:    req.Search,
			Sort:      req.Sort,
			Direction: req.Direction,
			Filters:   req.Filters,
		},
		Permissions: convertToPermissionsMap(permissions),
		Meta: PageMetadata{
			Component:    component,
			ResourceType: resourceType,
			PaginationConfig: PaginationMetadata{
				DefaultPageSize: 20,
				AllowedSizes:    []int{10, 20, 50, 100},
				MaxPageSize:     100,
			},
		},
		Stats: stats,
	}
}

// BuildStandardPagePropsWithPagination allows custom pagination config
func BuildStandardPagePropsWithPagination(
	result *PaginatedResult,
	req *ListRequest,
	permissions map[string]bool,
	stats map[string]interface{},
	component string,
	resourceType string,
	paginationConfig PaginationMetadata,
) PageProps {
	props := BuildStandardPageProps(result, req, permissions, stats, component, resourceType)
	props.Meta.PaginationConfig = paginationConfig
	return props
}