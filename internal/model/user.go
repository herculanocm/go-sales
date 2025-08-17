package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"column:id;type:uuid;primary_key"`
	Name      string         `gorm:"column:full_name;type:varchar(100);not null"`
	Email     string         `gorm:"column:email_address;type:varchar(100);uniqueIndex;not null"`
	Password  string         `gorm:"column:password_hash;not null"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (User) TableName() string {
	return "users"
}
