package dto

import (
	"time"
)

type CreateUserDTO struct {
	Name            string `json:"name" binding:"required,min=2"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	CompanyGlobalID int64  `json:"companyGlobalId" binding:"required"`

	// Adicionamos o campo para receber os IDs das roles.
	// "dive" diz ao validador para aplicar a regra "uuid" em cada elemento do array.
	RoleIDs []int64 `json:"roleIds" binding:"required,dive"`
}

type UserDTO struct {
	ID            int64            `json:"id"`
	Name          string           `json:"name"`
	Email         string           `json:"email"`
	Password      string           `json:"password"`
	Enabled       bool             `json:"enabled"`
	CompanyGlobal CompanyGlobalDTO `json:"companyGlobal"`
	Roles         []RoleDTO        `json:"roles"`
	CreatedAt     time.Time        `json:"createdAt"`
	UpdatedAt     time.Time        `json:"updatedAt"`
	DeletedAt     time.Time        `json:"deletedAt"`
}

type UpdateUserDTO struct {
	Name    *string `json:"name,omitempty" binding:"omitempty,max=255"`
	Email   *string `json:"email,omitempty" binding:"omitempty,email,max=150"`
	Enabled *bool   `json:"enabled,omitempty"`
}
