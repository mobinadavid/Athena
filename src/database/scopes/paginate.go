package scopes

import (
	"time"

	"gorm.io/gorm"
)

type QueryBuilderModel struct {
	Page          uint
	Limit         uint
	UserID        *uint
	SortBy        string
	SortOrder     string
	Filters       map[string]interface{}
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
}

type PaginatedModel struct {
	Limit       uint        `json:"limit"`
	CurrentPage uint        `json:"current_page"`
	TotalPages  int64       `json:"total_pages"`
	TotalItems  int64       `json:"total_items"`
	Items       interface{} `json:"items"`
}

func NewPaginatedModel(items interface{}, count int64, page, limit uint) *PaginatedModel {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}
	totalPages := int64(0)
	if limit > 0 {
		totalPages = (count + int64(limit) - 1) / int64(limit)
	}
	return &PaginatedModel{
		Limit:       limit,
		CurrentPage: page,
		TotalItems:  count,
		TotalPages:  totalPages,
		Items:       items,
	}
}

func PaginateScope(page uint, limit uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case limit > 100:
			limit = 100
		case limit <= 0:
			limit = 10
		}

		offset := (page - 1) * limit
		return db.Offset(int(offset)).Limit(int(limit))
	}
}
