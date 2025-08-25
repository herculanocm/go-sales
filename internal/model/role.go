package model

import (
	"go-sales/pkg/util"
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          util.UUID `gorm:"column:id;type:uuid;primaryKey"`
	Name        string    `gorm:"column:name;type:varchar(255);unique;not null"`
	Description *string   `gorm:"column:description;type:text"`

	CompanyGlobalID util.UUID     `gorm:"column:company_global_id;type:uuid;not null"`
	Permissions     []*Permission `gorm:"many2many:role_permissions;"` // Correto para muitos-para-muitos

	CanEdit   bool `gorm:"column:can_edit;type:boolean;default:true"`
	CanDelete bool `gorm:"column:can_delete;type:boolean;default:true"`
	IsAdmin   bool `gorm:"column:is_admin;type:boolean;default:false"`

	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz"`
}

func (Role) TableName() string {
	return "roles"
}
