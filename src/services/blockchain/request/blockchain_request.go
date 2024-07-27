package request

import (
	"gorm.io/datatypes"
)

type CreateBlockchainRequest struct {
	NativeAsset string         `json:"native_asset" validate:"required,max=255"`
	Title       datatypes.JSON ` json:"title" validate:""`
	Name        string         `json:"name" validate:"required,max=255"`
	IsActive    *bool          `json:"is_active" validate:"required,boolean"`
}
