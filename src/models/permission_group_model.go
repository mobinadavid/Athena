package models

import (
	"athena/src/api/errs"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	DefaultPermissionGroup = "default-permission-group"
	AllPermissionGroup     = "all-permissions"
)

type PermissionGroupModel struct {
	ID          uint               `json:"id" gorm:"primarykey"`
	Uuid        uuid.UUID          `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4(); uniqueIndex" filter:"true"`
	Name        string             `json:"name" filter:"true" gorm:"type:varchar(255); uniqueIndex" filter:"true" like:"true" sort:"true" search:"true"`
	Title       datatypes.JSON     `json:"title" gorm:"type:json" search:"true"`
	Description string             `json:"description" validate:"required,max=1000"`
	IsActive    *bool              `json:"is_active" gorm:"default:true"`
	Roles       []*RoleModel       `json:"-" gorm:"many2many:role_permission_groups"`
	Permissions []*PermissionModel `json:"permissions,omitempty" gorm:"many2many:permission_group_permissions"`
	CreatedAt   time.Time          `json:"created_at" sort:"true"`
	UpdatedAt   time.Time          `json:"updated_at" sort:"true"`
}

func (*PermissionGroupModel) TableName() string {
	return "permission_groups"
}

func (permissionGroup *PermissionGroupModel) BeforeDelete(tx *gorm.DB) (err error) {
	if permissionGroup.Name == DefaultPermissionGroup || permissionGroup.Name == AllPermissionGroup {
		return errs.CantDeletePermissionGroup
	}
	return nil
}
