package dto

import (
	"go-sales/pkg/util"
	"time"
)

type CreatePermissionDTO struct {
	Name            string    `json:"name" binding:"required,max=255"`
	CompanyGlobalID util.UUID `json:"companyGlobalID" binding:"required,uuid"`
	Description     *string   `json:"description" binding:"max=4000"`
}

type PermissionDTO struct {
	ID              util.UUID  `json:"id"`
	CompanyGlobalID util.UUID  `json:"companyGlobalID"`
	Name            string     `json:"name"`
	Description     *string    `json:"description"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	DeletedAt       *time.Time `json:"deletedAt"`
}
