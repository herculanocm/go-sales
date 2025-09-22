package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID   int64  `gorm:"column:id;type:bigint;primaryKey"`
	Name string `gorm:"column:full_name;type:varchar(255);not null"`

	Email           string     `gorm:"column:email_address;type:varchar(255)"`
	EmailRecovery   string     `gorm:"column:email_recovery;type:varchar(255)"`
	EmailVerified   bool       `gorm:"column:email_verified;type:boolean;not null;default:false"`
	EmailVerifiedAt *time.Time `gorm:"column:email_verified_at;type:timestamptz;"`

	Phone           string     `gorm:"column:phone_number;type:varchar(20)"`
	PhoneVerified   bool       `gorm:"column:phone_verified;type:boolean;not null;default:false"`
	PhoneVerifiedAt *time.Time `gorm:"column:phone_verified_at;type:timestamptz;"`

	Password string `gorm:"column:password_hash;not null"`
	Enabled  bool   `gorm:"column:enabled;type:boolean;not null;default:false"`

	Actived bool `gorm:"column:actived;type:boolean;not null;default:false"`

	ActivationKey *string `gorm:"column:activation_key;type:varchar(255);"`
	ActivatedAt   *string `gorm:"column:activated_at;type:timestamptz;"`

	ResetKey       *string    `gorm:"column:reset_key;type:varchar(255);"`
	ResetRequested *string    `gorm:"column:reset_requested;type:timestamptz;"`
	ResetAt        *time.Time `gorm:"column:reset_at;type:timestamptz;"`

	// 1. Adiciona o campo para a chave estrangeira.
	// O GORM usará isso para criar a coluna `company_global_id` na tabela `users`.
	CompanyGlobalID int64 `gorm:"column:company_global_id;type:bigint;not null"`

	// 2. Adiciona a struct para a associação.
	// Isso permite que você acesse os dados da empresa associada, por exemplo, com `user.CompanyGlobal`.
	// A tag `foreignKey` diz ao GORM qual campo nesta struct (`User`) usar para a junção.
	CompanyGlobal CompanyGlobal `gorm:"foreignKey:CompanyGlobalID"`

	Roles []*Role `gorm:"many2many:user_roles;"`

	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz"`
}

func (User) TableName() string {
	return "users"
}
