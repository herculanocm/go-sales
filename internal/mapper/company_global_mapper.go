package mapper

import (
	"go-sales/internal/dto"
	"go-sales/internal/model"
	"time"
)

func MapToCompanyGlobalDTO(company *model.CompanyGlobal) *dto.CompanyGlobalDTO {

	if company == nil {
		return nil
	}

	var deletedAt *time.Time
	// Verifica se o campo DeletedAt é válido (não é NULL no banco).
	if company.DeletedAt.Valid {
		// Se for válido, pegamos o endereço da variável de tempo.
		deletedAt = &company.DeletedAt.Time
	}

	var description *string
	if company.Description != "" {
		description = &company.Description
	}

	return &dto.CompanyGlobalDTO{
		ID:          company.ID,
		Name:        company.Name,
		Description: description,
		CGC:         company.CGC,
		Enabled:     company.Enabled,
		CreatedAt:   company.CreatedAt,
		UpdatedAt:   company.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}

func MapToCompanyGlobalDTOs(companies []*model.CompanyGlobal) *[]dto.CompanyGlobalDTO {
	if companies == nil {
		empty := make([]dto.CompanyGlobalDTO, 0)
		return &empty
	}

	dtos := make([]dto.CompanyGlobalDTO, 0, len(companies))
	for _, company := range companies {
		if dto := MapToCompanyGlobalDTO(company); dto != nil {
			dtos = append(dtos, *dto)
		}
	}
	return &dtos
}
