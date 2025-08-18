package mapper

import (
	"go-sales/internal/dto"
	"go-sales/internal/model"
	"time"
)

func MapToCompanyGlobalDTO(company *model.CompanyGlobal) *dto.CompanyGlobalDTO {
	var deletedAt *time.Time
	// Verifica se o campo DeletedAt é válido (não é NULL no banco).
	if company.DeletedAt.Valid {
		// Se for válido, pegamos o endereço da variável de tempo.
		deletedAt = &company.DeletedAt.Time
	}

	return &dto.CompanyGlobalDTO{
		ID:          company.ID.String(),
		Name:        company.Name,
		Description: company.Description,
		CGC:         company.CGC,
		Enabled:     company.Enabled,
		CreatedAt:   company.CreatedAt,
		UpdatedAt:   company.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}
