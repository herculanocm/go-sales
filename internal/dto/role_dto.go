package dto

import (
	"time"
)

type PermissionUtilHelper struct {
	ID int64 `json:"id" binding:"required,snowflake"`
}

type CreateRoleDTO struct {
	Name            string                 `json:"name" binding:"required,max=255"`
	Description     *string                `json:"description,omitempty" binding:"omitempty,max=4000"`
	CompanyGlobalID int64                  `json:"companyGlobalId" binding:"required,snowflake"`
	Permissions     []PermissionUtilHelper `json:"permissions" binding:"required,min=1,dive"`
	CanEdit         bool                   `json:"canEdit"`
	CanDelete       bool                   `json:"canDelete"`
	IsAdmin         bool                   `json:"isAdmin"`
}

type RoleDTO struct {
	ID              int64           `json:"id" binding:"required,snowflake"`
	Name            string          `json:"name" binding:"required,max=255"`
	Description     *string         `json:"description,omitempty" binding:"omitempty,max=4000"`
	Permissions     []PermissionDTO `json:"permissions" binding:"required,min=1,dive"`
	CompanyGlobalID int64           `json:"companyGlobalId" binding:"required,snowflake"`

	CanEdit   bool `json:"canEdit"`
	CanDelete bool `json:"canDelete"`
	IsAdmin   bool `json:"isAdmin"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
