package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type PermissionModel struct {
	ID               uint                    `json:"id" gorm:"primarykey"`
	Uuid             uuid.UUID               `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4(); uniqueIndex" filter:"true"`
	PermissionGroups []*PermissionGroupModel `json:"-" gorm:"many2many:permission_group_permissions"`
	Name             string                  `json:"name" filter:"true" gorm:"type:varchar(255); uniqueIndex" filter:"true" like:"true" sort:"true"`
	Title            datatypes.JSON          `json:"title" gorm:"type:json" like:"true"`
	CreatedAt        time.Time               `json:"created_at" sort:"true"`
	UpdatedAt        time.Time               `json:"updated_at" sort:"true"`
}

func (*PermissionModel) TableName() string {
	return "permissions"
}
