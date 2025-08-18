package model

import (
	"go-sales/pkg/util"
	"time"

	"gorm.io/gorm"
)

type CompanyGlobal struct {
	ID          util.UUID `gorm:"column:id;type:uuid;primary_key"`
	Name        string    `gorm:"column:name;type:varchar(255);not null"`
	Description string    `gorm:"column:description;type:text"`
	CGC         string    `gorm:"column:cgc;type:varchar(14);not null;index:idx_company_globals_cgc"`
	Enabled     bool      `gorm:"column:enabled;type:boolean;not null;default:false"`

	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz;index:idx_company_globals_deleted_at"`
}

func (CompanyGlobal) TableName() string {
	return "company_globals"
}
