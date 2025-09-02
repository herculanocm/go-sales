package mapper

import (
	"go-sales/internal/dto"
	"go-sales/internal/model"
)

// MapToPermissionDTO converte um model.Permission para dto.PermissionDTO.
func MapToPermissionDTO(permission *model.Permission) *dto.PermissionDTO {
	if permission == nil {
		return nil
	}

	return &dto.PermissionDTO{
		ID:   permission.ID,
		Name: permission.Name,

		CompanyGlobalID: permission.CompanyGlobalID,
		Description:     permission.Description,
		CreatedAt:       permission.CreatedAt,
		UpdatedAt:       permission.UpdatedAt,
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

func MapToPermission(dto *dto.CreatePermissionDTO) *model.Permission {
	if dto == nil {
		return nil
	}

	return &model.Permission{
		ID:              0,
		CompanyGlobalID: dto.CompanyGlobalID,
		Name:            dto.Name,
		Description:     dto.Description,
	}
}

func MapToPermissions(dtos []*dto.CreatePermissionDTO) []*model.Permission {
	if dtos == nil {
		return nil
	}

	permissions := make([]*model.Permission, 0, len(dtos))
	for _, dto := range dtos {
		if permission := MapToPermission(dto); permission != nil {
			permissions = append(permissions, permission)
		}
	}
	return permissions
}
