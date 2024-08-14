package scopes

import "gorm.io/gorm"

func IsActive() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = ?", true)
	}
}

// to check if a record is active and not allocated yet
func IsNotAllocated() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return IsActive()(db).Where("allocated_at IS NULL")
	}
}

// to check if a record is active and allocated
func IsAllocated() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return IsActive()(db).Where("allocated_at IS NOT NULL")
	}
}
