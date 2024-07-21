package scopes

import (
	"gorm.io/gorm"
)

type PaginateModel struct {
	TotalItems int64       `json:"total_items"`
	Items      interface{} `json:"items"`
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
