package model

import (
	"go-sales/pkg/util"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       util.UUID `gorm:"column:id;type:uuid;primary_key"`
	Name     string    `gorm:"column:full_name;type:varchar(255);not null"`
	Email    string    `gorm:"column:email_address;type:varchar(255);index:idx_users_email_address;not null"`
	Password string    `gorm:"column:password_hash;not null"`
	Enabled  bool      `gorm:"column:enabled;type:boolean;not null;default:false"`

	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz;index:idx_users_deleted_at"`
}

func (User) TableName() string {
	return "users"
}
