package dto

import (
	"go-sales/pkg/util"
	"time"
)

type CreatePermissionDTO struct {
	Name        string  `json:"name" binding:"required,max=255"`
	Description *string `json:"description" binding:"max=4000"`
}

type PermissionDTO struct {
	ID          util.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
