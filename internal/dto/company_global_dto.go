package dto

import (
	"go-sales/pkg/util"
	"time"
)

type CreateCompanyGlobalDTO struct {
	Name        string  `json:"name" binding:"required,max=255"`
	SocialName  string  `json:"social_name" binding:"required,max=255"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=4000"`
	CGC         string  `json:"cgc" binding:"required,max=14"`
	Enabled     bool    `json:"enabled"`
}

type CompanyGlobalDTO struct {
	ID          util.UUID  `json:"id"`
	Name        string     `json:"name"`
	SocialName  string     `json:"social_name"`
	Description *string    `json:"description,omitempty"`
	CGC         string     `json:"cgc"`
	Enabled     bool       `json:"enabled"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}
