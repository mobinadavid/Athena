package AuthorizationRequests

import "gorm.io/datatypes"

type CreatePermissionGroupRequest struct {
	Name        string         `json:"name" validate:"required,english-only"`
	Title       datatypes.JSON `json:"title" validate:"required,required-fa-title"`
	Permissions []string       `json:"permissions" validate:"required,min=1,max=100,dive,is-uuid,exists=permissions:uuid"`
	Description string         `json:"description" validate:"omitempty,max=1000"`
}
