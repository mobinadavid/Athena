package AdminRequests

type ChangeActivationRequest struct {
	IsActive *bool  `json:"is_active" validate:"required"`
	Reason   string `json:"reason" validate:"omitempty,max=255"`
}
