package dto

import (
	"athena/src/models"
	"time"

	"github.com/google/uuid"
)

type AdminListModel struct {
	Uuid         uuid.UUID         `json:"uuid"`
	Name         string            `json:"name"`
	LastName     string            `json:"last_name"`
	Mobile       string            `json:"mobile"`
	ProfileImage models.Attachment `json:"profile_image"`
	Username     string            `json:"username"`
	Roles        []AdminRole       `json:"roles"`
	IsActive     *bool             `json:"is_active"`
	CreatedDate  time.Time         `json:"created_date"`
	UpdatedDate  time.Time         `json:"updated_date"`
	Description  string            `json:"description"`
	Reason       string            `json:"reason"`
}

type AdminRole struct {
	Title string `json:"title"`
}
