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
	Email       *string `json:"email,omitempty" binding:"omitempty,max=150"`

	Address  *CreateCompanyGlobalAddressDTO   `json:"address,omitempty" binding:"required"`
	Contacts []*CreateCompanyGlobalContactDTO `json:"contacts,omitempty" binding:"required,min=1,dive"`
}

type CreateCompanyGlobalAddressDTO struct {
	Street           string  `json:"street" binding:"required,max=255"`
	StreetNumber     *string `json:"street_number,omitempty" binding:"omitempty,max=50"`
	StreetComplement *string `json:"street_complement,omitempty" binding:"omitempty,max=255"`
	City             string  `json:"city" binding:"required,max=100"`
	State            string  `json:"state" binding:"required,max=100"`
	PostalCode       string  `json:"postal_code" binding:"required,max=20"`
	Country          string  `json:"country" binding:"required,max=100"`
}

type CreateCompanyGlobalContactDTO struct {
	Name  string  `json:"name" binding:"required,max=255"`
	Email *string `json:"email,omitempty" binding:"omitempty,max=150"`
	Phone *string `json:"phone,omitempty" binding:"omitempty,max=20"`
	CGC   *string `json:"cgc,omitempty" binding:"omitempty,max=40"`
}

type CompanyGlobalAddressDTO struct {
	ID               *util.UUID `json:"id"`
	Street           string     `json:"street"`
	StreetNumber     *string    `json:"street_number,omitempty"`
	StreetComplement *string    `json:"street_complement,omitempty"`
	City             string     `json:"city"`
	State            string     `json:"state"`
	PostalCode       string     `json:"postal_code"`
	Country          string     `json:"country"`
}

type CompanyGlobalContactDTO struct {
	ID    *util.UUID `json:"id"`
	Name  string     `json:"name"`
	Email *string    `json:"email,omitempty"`
	Phone *string    `json:"phone,omitempty"`
	CGC   *string    `json:"cgc,omitempty"`
}

type CompanyGlobalDTO struct {
	ID          *util.UUID                 `json:"id"`
	Name        string                     `json:"name"`
	SocialName  string                     `json:"social_name"`
	Description *string                    `json:"description,omitempty"`
	CGC         string                     `json:"cgc"`
	Enabled     bool                       `json:"enabled"`
	Email       *string                    `json:"email,omitempty"`
	CreatedAt   time.Time                  `json:"created_at"`
	UpdatedAt   time.Time                  `json:"updated_at"`
	DeletedAt   *time.Time                 `json:"deleted_at,omitempty"`
	Address     *CompanyGlobalAddressDTO   `json:"address,omitempty"`
	Contacts    []*CompanyGlobalContactDTO `json:"contacts,omitempty"`
}
