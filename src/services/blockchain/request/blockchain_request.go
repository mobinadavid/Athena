package request

import (
	"gorm.io/datatypes"
)

type CreateBlockchainRequest struct {
	NativeAsset string         `json:"native-asset" validate:"required,max=255"`
	Title       datatypes.JSON ` json:"title" validate:""`
	IsActive    *bool          `json:"is_active" validate:"required,boolean"`
}
