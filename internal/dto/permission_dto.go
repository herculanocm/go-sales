package dto

import (
	"time"
)

type CreatePermissionDTO struct {
	Name            string  `json:"name" binding:"required,max=255"`
	CompanyGlobalID int64   `json:"companyGlobalID" binding:"required"`
	Description     *string `json:"description" binding:"omitempty,max=4000"`
}

type PermissionDTO struct {
	ID              int64     `json:"id" binding:"required,snowflake"`
	CompanyGlobalID int64     `json:"companyGlobalID" binding:"required,snowflake"`
	Name            string    `json:"name" binding:"required,max=255"`
	Description     *string   `json:"description,omitempty" binding:"omitempty,max=4000"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
