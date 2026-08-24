package Authentication

type RefreshAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"omitempty"`
}
