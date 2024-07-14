package request

import (
	"gorm.io/datatypes"
)

type CreateBlockchainRequest struct {
	Symbol   string         `json:"symbol" validate:"required,max=255"`
	Title    datatypes.JSON ` json:"title" validate:""`
	IsActive *bool          `json:"is_active" validate:"required,boolean"`
}
