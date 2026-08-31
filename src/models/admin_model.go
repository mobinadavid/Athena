package models

import (
	"athena/src/hash"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type AdminModel struct {
	ID             uint                `json:"id" gorm:"primarykey"`
	Uuid           uuid.UUID           `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4(); uniqueIndex" filter:"true"`
	IsActive       *bool               `json:"is_active" gorm:"type:bool; default:false" filter:"true"`
	FirstName      string              `json:"first_name" gorm:"type:varchar(255); default:null" filter:"true" like:"true" sort:"true" search:"true"`
	LastName       string              `json:"last_name" gorm:"type:varchar(255); default:null" filter:"true" like:"true" sort:"true" search:"true"`
	Username       string              `json:"username" gorm:"type:varchar(255);unique; default:null" like:"true" search:"true"`
	Password       []byte              `json:"-" gorm:"type:text; default:null"`
	Mobile         string              `json:"mobile" gorm:"type:varchar(100); uniqueIndex; not null" filter:"true" like:"true"`
	Email          string              `json:"email" gorm:"type:varchar(100); uniqueIndex; default:null" filter:"true" like:"true"`
	Roles          []*RoleModel        `json:"roles" gorm:"many2many:admin_role" join:"true" filter:"true"`
	AccessTokens   []*AccessTokenModel `json:"access_tokens" gorm:"polymorphic:Owner;"`
	TotpSecret     []byte              `json:"-" gorm:"type:bytea; default:null"`
	TotpSecretUrl  []byte              `json:"-" gorm:"type:bytea; default:null"`
	TwoFaEnabled   bool                `json:"-" gorm:"type:bool;default:false"`
	RecoveryCodes  pq.StringArray      `json:"-" gorm:"type:text[];default:null"`
	AdminImageUuid string              `json:"-" gorm:"type:varchar(255)"`
	AdminImage     Attachment          `json:"admin_image" gorm:"-"`
	CreatedAt      time.Time           `json:"created_at" sort:"true"`
	UpdatedAt      time.Time           `json:"updated_at" sort:"true"`
	DeletedAt      gorm.DeletedAt      `json:"deleted_at" gorm:"index" sort:"true"`
}

func (*AdminModel) TableName() string {
	return "admins"
}

func (admin *AdminModel) BeforeSave(tx *gorm.DB) error {
	if len(admin.Password) > 0 {
		if admin.ID > 0 {
			var existing AdminModel
			// Find existing user by ID for update case
			res := tx.First(&existing, admin.ID)
			if res.Error == nil && string(existing.Password) == string(admin.Password) {
				// Password hasn't changed, skip re-hashing
				return nil
			}
		}

		// In this case user is not exist so this is before Create
		hashedPassword, err := hash.GetInstance().Generate(admin.Password)
		if err != nil {
			return err
		}
		admin.Password = hashedPassword
	}
	return nil
}
