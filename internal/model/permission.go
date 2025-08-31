package model

import (
	"go-sales/pkg/util"
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID              util.UUID `gorm:"column:id;type:uuid;primaryKey"`
	CompanyGlobalID util.UUID `gorm:"column:company_global_id;type:uuid;not null;index"`
	Name            string    `gorm:"column:name;unique;not null"`
	Description     *string   `gorm:"column:description"`

	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz"`
}

func (Permission) TableName() string {
	return "permissions"
}
