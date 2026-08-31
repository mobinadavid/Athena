package models

import (
	"athena/src/hash"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserModel struct {
	ID       uint      `json:"id,omitempty" gorm:"primarykey"`
	Uuid     uuid.UUID `json:"uuid,omitempty" gorm:"type:uuid;default:uuid_generate_v4(); uniqueIndex" filter:"true"`
	IsActive bool      `json:"is_active,omitempty" gorm:"type:bool; default:true" filter:"true" like:"true"`
	// Personal information.
	FirstName        string         `json:"first_name,omitempty" gorm:"type:varchar(255); default:null" filter:"true" like:"true" sort:"true"`
	LastName         string         `json:"last_name,omitempty" gorm:"type:varchar(255); default:null" filter:"true" like:"true" sort:"true"`
	FullName         string         `json:"full_name,omitempty" gorm:"-" filter:"true" sort:"true"`
	ProfileImageUuid string         `json:"-" gorm:"type:varchar(255)"`
	ProfileImage     Attachment     `json:"profile_image,omitempty" gorm:"-"`
	Password         []byte         `json:"-" gorm:"type:text; default:null"`
	TotpSecret       []byte         `json:"-" gorm:"type:bytea; default:null"`
	TotpSecretUrl    []byte         `json:"-" gorm:"type:bytea; default:null"`
	TwoFaEnabled     bool           `json:"-" gorm:"type:bool;default:false"`
	RecoveryCodes    pq.StringArray `gorm:"-"`
	// Identities.
	NationalIdentityCode string `json:"national_identity_code,omitempty" gorm:"type:varchar(255); uniqueIndex; default:null" filter:"true" like:"true"`
	// Contact information.
	Mobile string `json:"mobile,omitempty" gorm:"type:varchar(100); uniqueIndex; not null" filter:"true" like:"true"`
	Email  string `json:"email,omitempty" gorm:"type:varchar(100); default:null" filter:"true" like:"true"`
	// Relationships
	AccessTokens []*AccessTokenModel `json:"-" gorm:"polymorphic:Owner;"`
	Roles        []*RoleModel        `json:"roles,omitempty" gorm:"many2many:role_user"`
	// Times
	CreatedAt time.Time      `json:"created_at,omitempty" sort:"true"`
	UpdatedAt time.Time      `json:"updated_at,omitempty" sort:"true"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" sort:"true"`
}

func (*UserModel) TableName() string {
	return "users"
}

func (user *UserModel) BeforeSave(tx *gorm.DB) error {
	if len(user.Password) > 0 {
		if user.ID > 0 {
			var existing UserModel
			// Find existing user by ID for update case
			res := tx.First(&existing, user.ID)
			if res.Error == nil && string(existing.Password) == string(user.Password) {
				// Password hasn't changed, skip re-hashing
				return nil
			}
		}

		// In this case user is not exist so this is before Create
		hashedPassword, err := hash.GetInstance().Generate(user.Password)
		if err != nil {
			return err
		}

		user.Password = hashedPassword
	}
	return nil
}

type Attachment struct {
	Name string `json:"file_name"`
	URL  string `json:"url"`
}
