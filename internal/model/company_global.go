package model

import (
	"go-sales/pkg/util"
	"time"

	"gorm.io/gorm"
)

type CompanyGlobal struct {
	ID          util.UUID `gorm:"column:id;type:uuid;primary_key"`
	Name        string    `gorm:"column:name;type:varchar(255)"`
	SocialName  string    `gorm:"column:social_name;type:varchar(255)"`
	Description *string   `gorm:"column:description;type:text"`
	CGC         string    `gorm:"column:cgc;type:varchar(14)"`
	Enabled     bool      `gorm:"column:enabled;type:boolean"`

	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz;autoCreateTime;<-:create"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz"`
}

func (CompanyGlobal) TableName() string {
	return "company_globals"
}
