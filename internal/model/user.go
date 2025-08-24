package model

import (
	"go-sales/pkg/util"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       util.UUID `gorm:"column:id;type:uuid;primaryKey"`
	Name     string    `gorm:"column:full_name;type:varchar(255);not null"`
	Email    string    `gorm:"column:email_address;type:varchar(255)"`
	Password string    `gorm:"column:password_hash;not null"`
	Enabled  bool      `gorm:"column:enabled;type:boolean;not null;default:false"`

	// 1. Adiciona o campo para a chave estrangeira.
	// O GORM usará isso para criar a coluna `company_global_id` na tabela `users`.
	CompanyGlobalID util.UUID `gorm:"column:company_global_id;type:uuid;not null"`

	// 2. Adiciona a struct para a associação.
	// Isso permite que você acesse os dados da empresa associada, por exemplo, com `user.CompanyGlobal`.
	// A tag `foreignKey` diz ao GORM qual campo nesta struct (`User`) usar para a junção.
	CompanyGlobal CompanyGlobal `gorm:"foreignKey:CompanyGlobalID"`

	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz"`
}

func (User) TableName() string {
	return "users"
}
