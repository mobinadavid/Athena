package UserRequests

// CreateRequest struct for validating incoming request data for creating a user
type CreateRequest struct {
	FirstName            string `json:"first_name" validate:"required,max=255"`
	LastName             string `json:"last_name" validate:"required,max=255"`
	ProfileImage         string `json:"profile_image" validate:"required"`
	NationalIdentityCode string `json:"national_identity_code" validate:"required,national-code-or-company-id"`
	Mobile               string `json:"mobile" validate:"required,iranian-mobile"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,max=255,is-strong-password"`
	IsActive             bool   `json:"is_active" validate:"required,boolean"`
}
