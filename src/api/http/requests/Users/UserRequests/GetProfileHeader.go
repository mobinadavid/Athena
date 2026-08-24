package UserRequests

type GetProfileHeaderRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
