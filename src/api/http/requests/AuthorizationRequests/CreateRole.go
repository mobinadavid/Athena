package AuthorizationRequests

type CreateRoleRequest struct {
	Name             string   `json:"name" validate:"required,english-only"`
	Title            string   `json:"title" validate:"required"`
	PermissionGroups []string `json:"permission_groups" validate:"required,min=1,max=100,dive,is-uuid,exists=permission_groups:uuid"`
	Description      string   `json:"description" validate:"omitempty,max=1000"`
}
