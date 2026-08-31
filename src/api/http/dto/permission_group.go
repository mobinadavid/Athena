package dto

import (
	"athena/src/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type PermissionGroupListModel struct {
	ID               uint                      `json:"id" `
	Uuid             uuid.UUID                 `json:"uuid"`
	Name             string                    `json:"name"`
	Title            datatypes.JSON            `json:"title"`
	IsActive         bool                      `json:"is_active"`
	Description      string                    `json:"description"`
	Permissions      []*models.PermissionModel `json:"permissions" `
	PermissionsCount int                       `json:"permissions_count"`
	CreatedAt        time.Time                 `json:"created_at"`
	UpdatedAt        time.Time                 `json:"updated_at"`
}
