package dto

import (
	"go-sales/pkg/util"
	"time"
)

type CreateUserDTO struct {
	Name            string    `json:"name" binding:"required,min=2"`
	Email           string    `json:"email" binding:"required,email"`
	Password        string    `json:"password" binding:"required,min=8"`
	CompanyGlobalID util.UUID `json:"company_global_id" binding:"required,uuid"`

	// Adicionamos o campo para receber os IDs das roles.
	// "dive" diz ao validador para aplicar a regra "uuid" em cada elemento do array.
	RoleIDs []util.UUID `json:"role_ids" binding:"required,dive,uuid"`
}

type UserDTO struct {
	ID            util.UUID        `json:"id"`
	Name          string           `json:"name"`
	Email         string           `json:"email"`
	Password      string           `json:"password"`
	Enabled       bool             `json:"enabled"`
	CompanyGlobal CompanyGlobalDTO `json:"company_global"`
	Roles         []RoleDTO        `json:"roles"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	DeletedAt     time.Time        `json:"deleted_at"`
}
