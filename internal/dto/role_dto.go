package dto

import (
	"go-sales/pkg/util"
	"time"
)

type CreateRoleDTO struct {
	Name            string      `json:"name" binding:"required,max=255"`
	Description     *string     `json:"description" binding:"max=4000"`
	CompanyGlobalID util.UUID   `json:"company_global_id" binding:"required,uuid"`
	PermissionIDs   []util.UUID `json:"permission_ids" binding:"dive,uuid"`
	CanEdit         bool        `json:"can_edit"`
	CanDelete       bool        `json:"can_delete"`
	IsAdmin         bool        `json:"is_admin"`
}

type RoleDTO struct {
	ID              util.UUID       `json:"id"`
	Name            string          `json:"name"`
	Description     *string         `json:"description"`
	Permissions     []PermissionDTO `json:"permissions"`
	CompanyGlobalID util.UUID       `json:"company_global_id"`

	CanEdit   bool `json:"can_edit"`
	CanDelete bool `json:"can_delete"`
	IsAdmin   bool `json:"is_admin"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
