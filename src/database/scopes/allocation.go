package scopes

import "gorm.io/gorm"

func IsActive() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = ?", true)
	}
}
func IsAllocated() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return IsActive()(db).Where("allocated_at IS NULL")
	}
}
