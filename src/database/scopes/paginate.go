package scopes

import (
	"gorm.io/gorm"
	"time"
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

type PaginateModel struct {
	Limit       uint        `json:"limit"`
	CurrentPage uint        `json:"current_page"`
	TotalPages  int64       `json:"total_pages"`
	TotalItems  int64       `json:"total_items"`
	Items       interface{} `json:"items"`
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
