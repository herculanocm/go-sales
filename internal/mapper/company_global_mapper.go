package mapper

import (
	"go-sales/internal/dto"
	"go-sales/internal/model"
	"go-sales/pkg/util"
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

	return &dto.CompanyGlobalDTO{
		ID:          company.ID,
		Name:        company.Name,
		SocialName:  company.SocialName,
		Description: company.Description,
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

func MapToCreateCompanyGlobal(companyDTO *dto.CreateCompanyGlobalDTO) *model.CompanyGlobal {
	if companyDTO == nil {
		return nil
	}

	return &model.CompanyGlobal{
		ID:          util.New(),
		Name:        companyDTO.Name,
		SocialName:  companyDTO.SocialName,
		Description: companyDTO.Description,
		CGC:         companyDTO.CGC,
		Enabled:     companyDTO.Enabled,
	}
}

func MapToUpdateCompanyGlobal(companyDTO *dto.CreateCompanyGlobalDTO, id util.UUID) *model.CompanyGlobal {
	if companyDTO == nil {
		return nil
	}

	return &model.CompanyGlobal{
		ID:          id,
		Name:        companyDTO.Name,
		SocialName:  companyDTO.SocialName,
		Description: companyDTO.Description,
		CGC:         companyDTO.CGC,
		Enabled:     companyDTO.Enabled,
	}
}
