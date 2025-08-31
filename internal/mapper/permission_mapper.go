package mapper

import (
	"go-sales/internal/dto"
	"go-sales/internal/model"
	"time"
)

// MapToPermissionDTO converte um model.Permission para dto.PermissionDTO.
func MapToPermissionDTO(permission *model.Permission) *dto.PermissionDTO {
	if permission == nil {
		return nil
	}

	var deletedAt *time.Time
	if permission.DeletedAt.Valid {
		deletedAt = &permission.DeletedAt.Time
	} else {
		deletedAt = nil
	}

	return &dto.PermissionDTO{
		ID:              permission.ID,
		Name:            permission.Name,
		CompanyGlobalID: permission.CompanyGlobalID,
		Description:     permission.Description,
		CreatedAt:       permission.CreatedAt,
		UpdatedAt:       permission.UpdatedAt,
		DeletedAt:       deletedAt,
	}
}

// MapToPermissionDTOs converte um slice de *model.Permission para um slice de dto.PermissionDTO.
func MapToPermissionDTOs(permissions []*model.Permission) *[]dto.PermissionDTO {
	if permissions == nil {
		empty := make([]dto.PermissionDTO, 0)
		return &empty
	}

	result := make([]dto.PermissionDTO, 0, len(permissions))
	for _, permission := range permissions {
		if dto := MapToPermissionDTO(permission); dto != nil {
			result = append(result, *dto)
		}
	}
	return &result
}
