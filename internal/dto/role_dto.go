package dto

import (
	"time"
)

type CreateRoleDTO struct {
	Name            string  `json:"name" binding:"required,max=255"`
	Description     *string `json:"description,omitempty" binding:"omitempty,max=4000"`
	CompanyGlobalID int64   `json:"companyGlobalId" binding:"required"`
	PermissionIDs   []int64 `json:"permissionIds" binding:"dive"`
	CanEdit         bool    `json:"canEdit"`
	CanDelete       bool    `json:"canDelete"`
	IsAdmin         bool    `json:"isAdmin"`
}

type RoleDTO struct {
	ID              int64           `json:"id"`
	Name            string          `json:"name"`
	Description     *string         `json:"description,omitempty"`
	Permissions     []PermissionDTO `json:"permissions"`
	CompanyGlobalID int64           `json:"companyGlobalId"`

	CanEdit   bool `json:"canEdit"`
	CanDelete bool `json:"canDelete"`
	IsAdmin   bool `json:"isAdmin"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
