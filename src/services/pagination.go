package services

import "athena/src/database/scopes"

func paginatedResult(params *scopes.QueryBuilderModel, items interface{}, count int64) *scopes.PaginatedModel {
	page := uint(1)
	limit := uint(10)
	if params != nil {
		if params.Page > 0 {
			page = params.Page
		}
		if params.Limit > 0 {
			limit = params.Limit
		}
	}

	totalPages := int64(0)
	if limit > 0 {
		totalPages = (count + int64(limit) - 1) / int64(limit)
	}

	return &scopes.PaginatedModel{
		Limit:       limit,
		CurrentPage: page,
		TotalItems:  count,
		TotalPages:  totalPages,
		Items:       items,
	}
}
