package models

import (
	"athena/src/api/errs"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	SuperAdminRole = "super-admin"
	DefaultRole    = "default-role"
	AdminRole      = "admin"
	UserRole       = "user"
)

type RoleModel struct {
	ID                   uint                    `json:"id" gorm:"primarykey"`
	Uuid                 uuid.UUID               `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4(); uniqueIndex" filter:"true"`
	Name                 string                  `json:"name" gorm:"type:varchar(255); uniqueIndex" filter:"true" like:"true" sort:"true" search:"true"`
	Title                string                  `json:"title" gorm:"type:varchar(255)" search:"true"`
	Description          string                  `json:"description" validate:"required,max=1000"`
	IsActive             *bool                   `json:"is_active" gorm:"default:true"`
	PermissionGroups     []*PermissionGroupModel `json:"permission_groups" gorm:"many2many:role_permission_groups"`
	PermissionGroupCount int                     `json:"permission_group_count" gorm:"-"`
	Admins               []*AdminModel           `json:"admins,omitempty" gorm:"many2many:admin_role"`
	TotalAssignees       int64                   `json:"total_assignees" gorm:"-"`
	CreatedAt            time.Time               `json:"created_at" sort:"true"`
	UpdatedAt            time.Time               `json:"updated_at" sort:"true"`
}

func (*RoleModel) TableName() string {
	return "roles"
}

func (role *RoleModel) BeforeDelete(tx *gorm.DB) (err error) {
	if role.Name == SuperAdminRole || role.Name == DefaultRole {
		return errs.CantDeleteThisRole
	}
	return nil
}
