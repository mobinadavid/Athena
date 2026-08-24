package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccessTokenModel struct {
	ID                    uint           `json:"id" gorm:"primarykey"`
	Uuid                  uuid.UUID      `json:"uuid" gorm:"type:uuid; uniqueIndex" filter:"true"`
	OwnerID               uint           `json:"-"`
	OwnerType             string         `json:"-"`
	AccessToken           []byte         `json:"-" gorm:"type:text;not null"`
	AccessTokenExpiresAt  time.Time      `json:"access_token_expires_at"`
	RefreshToken          []byte         `json:"-" gorm:"type:text;not null"`
	RefreshTokenExpiresAt time.Time      `json:"refresh_token_expires_at" sort:"true"`
	IP                    string         `json:"ip"`
	UserAgent             string         `json:"user_agent"`
	LastUsedAt            *time.Time     `json:"last_used_at" sort:"true"`
	CreatedAt             time.Time      `json:"created_at" sort:"true"`
	UpdatedAt             time.Time      `json:"updated_at" sort:"true"`
	DeletedAt             gorm.DeletedAt `json:"deleted_at" gorm:"index" sort:"true"`
}

func (*AccessTokenModel) TableName() string {
	return "access_tokens"
}
