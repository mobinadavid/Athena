package AuthenticationRequests

type RefreshAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
