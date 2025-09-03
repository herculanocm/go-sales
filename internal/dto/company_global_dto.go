package dto

import (
	"time"
)

type CreateCompanyGlobalDTO struct {
	Name        string  `json:"name" binding:"required,max=255"`
	SocialName  string  `json:"socialName" binding:"required,max=255"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=4000"`
	CGC         string  `json:"cgc" binding:"required,max=14"`
	Enabled     bool    `json:"enabled"`
	Email       *string `json:"email" binding:"required,max=150"`

	Address  *CreateCompanyGlobalAddressDTO   `json:"address,omitempty" binding:"required"`
	Contacts []*CreateCompanyGlobalContactDTO `json:"contacts,omitempty" binding:"required,min=1,dive"`
}

type CreateCompanyGlobalAddressDTO struct {
	Street           string  `json:"street" binding:"required,max=255"`
	StreetNumber     *string `json:"streetNumber,omitempty" binding:"omitempty,max=50"`
	StreetComplement *string `json:"streetComplement,omitempty" binding:"omitempty,max=255"`
	City             string  `json:"city" binding:"required,max=100"`
	State            string  `json:"state" binding:"required,max=100"`
	PostalCode       string  `json:"postalCode" binding:"required,max=20"`
	Country          string  `json:"country" binding:"required,max=100"`
}

type CreateCompanyGlobalContactDTO struct {
	Name  string  `json:"name" binding:"required,max=255"`
	Email *string `json:"email" binding:"required,max=150"`
	Phone *string `json:"phone" binding:"required,max=20"`
	CGC   *string `json:"cgc" binding:"required,max=40"`
}

type CompanyGlobalAddressDTO struct {
	ID               int64   `json:"id"`
	Street           string  `json:"street"`
	StreetNumber     *string `json:"streetNumber,omitempty"`
	StreetComplement *string `json:"streetComplement,omitempty"`
	City             string  `json:"city"`
	State            string  `json:"state"`
	PostalCode       string  `json:"postalCode"`
	Country          string  `json:"country"`
}

type CompanyGlobalContactDTO struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Email *string `json:"email,omitempty"`
	Phone *string `json:"phone,omitempty"`
	CGC   *string `json:"cgc,omitempty"`
}

type CompanyGlobalDTO struct {
	ID          int64                      `json:"id"`
	Name        string                     `json:"name"`
	SocialName  string                     `json:"socialName"`
	Description *string                    `json:"description,omitempty"`
	CGC         string                     `json:"cgc"`
	Enabled     bool                       `json:"enabled"`
	Email       *string                    `json:"email,omitempty"`
	CreatedAt   time.Time                  `json:"createdAt"`
	UpdatedAt   time.Time                  `json:"updatedAt"`
	DeletedAt   *time.Time                 `json:"deletedAt,omitempty"`
	Address     *CompanyGlobalAddressDTO   `json:"address,omitempty"`
	Contacts    []*CompanyGlobalContactDTO `json:"contacts,omitempty"`
}

func (dto *CreateCompanyGlobalDTO) ValidateContacts() bool {
	contacts := dto.Contacts
	if len(contacts) == 0 || len(contacts) == 1 {
		return true
	}

	emailSet := make(map[string]struct{})
	phoneSet := make(map[string]struct{})
	cgcSet := make(map[string]struct{})

	for _, contact := range contacts {
		if contact.Email != nil {
			email := *contact.Email
			if email != "" {
				if _, exists := emailSet[email]; exists {
					return false
				}
				emailSet[email] = struct{}{}
			}
		}
		if contact.Phone != nil {
			phone := *contact.Phone
			if phone != "" {
				if _, exists := phoneSet[phone]; exists {
					return false
				}
				phoneSet[phone] = struct{}{}
			}
		}
		if contact.CGC != nil {
			cgc := *contact.CGC
			if cgc != "" {
				if _, exists := cgcSet[cgc]; exists {
					return false
				}
				cgcSet[cgc] = struct{}{}
			}
		}
	}
	return true
}
