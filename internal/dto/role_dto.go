package dto

import (
	"go-sales/pkg/util"
	"time"
)

type CreateRoleDTO struct {
	Name            string      `json:"name" binding:"required,max=255"`
	Description     *string     `json:"description,omitempty" binding:"omitempty,max=4000"`
	CompanyGlobalID util.UUID   `json:"companyGlobalId" binding:"required,uuid"`
	PermissionIDs   []util.UUID `json:"permissionIds" binding:"dive,uuid"`
	CanEdit         bool        `json:"canEdit"`
	CanDelete       bool        `json:"canDelete"`
	IsAdmin         bool        `json:"isAdmin"`
}

type RoleDTO struct {
	ID              util.UUID       `json:"id"`
	Name            string          `json:"name"`
	Description     *string         `json:"description,omitempty"`
	Permissions     []PermissionDTO `json:"permissions"`
	CompanyGlobalID util.UUID       `json:"companyGlobalId"`

	CanEdit   bool `json:"canEdit"`
	CanDelete bool `json:"canDelete"`
	IsAdmin   bool `json:"isAdmin"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
