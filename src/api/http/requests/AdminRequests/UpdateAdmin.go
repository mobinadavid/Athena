package AdminRequests

// UpdateAdminRequest struct for validating incoming request data for updating a user
type UpdateAdminRequest struct {
	FirstName     string   `json:"first_name" validate:"required,omitempty,is-persian-string,max=255"`
	LastName      string   `json:"last_name" validate:"required,omitempty,is-persian-string,max=255"`
	Mobile        string   `json:"mobile" validate:"required,omitempty,iranian-mobile"`
	Username      string   `json:"username" validate:"required,omitempty,max=255,username"`
	Password      string   `json:"password" validate:"required,omitempty,max=255,is-strong-password"`
	ProfileImage  string   `json:"profile_image"`
	Roles         []string `json:"roles" validate:"required,omitempty,max=100,min=1,dive,is-uuid,exists=roles:uuid"`
	Description   string   `json:"description" validate:"omitempty,max=255"`
	AllIPsAllowed *bool    `json:"all_ips_allowed" validate:"required,omitempty"`
	AllowedIPs    []string `json:"allowed_ips" validate:"omitempty,dive,ip"`
	IsActive      *bool    `json:"is_active" validate:"omitempty,boolean"`
}
