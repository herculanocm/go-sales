package model

import (
	"time"
)

type Permission struct {
	ID              int64   `gorm:"column:id;type:bigint;primaryKey"`
	CompanyGlobalID int64   `gorm:"column:company_global_id;type:bigint;not null"`
	Name            string  `gorm:"column:name;unique;not null"`
	Description     *string `gorm:"column:description"`

	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamptz"`
}

func (Permission) TableName() string {
	return "permissions"
}
