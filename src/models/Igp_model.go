package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"time"
)

type IGPModel struct {
	ID                    uint           `json:"id" gorm:"primarykey"`
	Uuid                  uuid.UUID      `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();uniqueIndex"`
	Ipg                   string         `json:"ipg" gorm:"type:varchar(255); not null;"`
	IssuerReferenceNumber string         `json:"issuer_reference_number" gorm:"type:varchar(255);default:null;uniqueIndex"`
	Receipt               datatypes.JSON `json:"receipt" gorm:"type:JSON;"`
	Amount                float64        `json:"amount" gorm:"type:numeric(15,6);default:0;check:amount >= 0"`
	Status                string         `json:"status" gorm:"type:varchar(255);not null;"`
	CallbackUrl           string         `json:"callback_url" gorm:"type:varchar(255);not null;"`
	CreatedAt             time.Time      `json:"created_at" gorm:"type:timestamp with time zone;default:current_timestamp"`
	UpdatedAt             time.Time      `json:"updated_at" gorm:"type:timestamp with time zone"`
}

func (*IGPModel) TableName() string {
	return "internet_gateway_payments"
}
