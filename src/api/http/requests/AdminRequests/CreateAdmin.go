package AdminRequests

// CreateAdminRequest struct for validating incoming request data for creating a user
type CreateAdminRequest struct {
	FirstName    string   `json:"first_name" validate:"required,is-persian-string,max=255"`
	LastName     string   `json:"last_name" validate:"required,is-persian-string,max=255"`
	Mobile       string   `json:"mobile" validate:"required,iranian-mobile"`
	Username     string   `json:"username" validate:"required,max=255,username"`
	Password     string   `json:"password" validate:"required,max=255,is-strong-password"`
	ProfileImage string   `json:"profile_image"`
	Roles        []string `json:"roles" validate:"required,max=100,min=1,dive,is-uuid,exists=roles:uuid"`
	IsActive     *bool    `json:"is_active" validate:"boolean"`
}
