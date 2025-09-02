package dto

import (
	"time"
)

type CreatePermissionDTO struct {
	Name            string  `json:"name" binding:"required,max=255"`
	CompanyGlobalID int64   `json:"companyGlobalID" binding:"required"`
	Description     *string `json:"description" binding:"max=4000"`
}

type PermissionDTO struct {
	ID              int64     `json:"id"`
	CompanyGlobalID int64     `json:"companyGlobalID"`
	Name            string    `json:"name"`
	Description     *string   `json:"description"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
