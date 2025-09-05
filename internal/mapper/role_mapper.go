package mapper

import (
	"go-sales/internal/dto"
	"go-sales/internal/model"
)

// MapToRoleDTO converte um model.Role para dto.RoleDTO.

func MapToRoleDTO(role *model.Role) *dto.RoleDTO {
	if role == nil {
		return nil
	}

	return &dto.RoleDTO{
		ID:              role.ID,
		Name:            role.Name,
		Description:     role.Description,
		Permissions:     *MapToPermissionDTOs(role.Permissions),
		CompanyGlobalID: role.CompanyGlobalID,
		CanEdit:         role.CanEdit,
		CanDelete:       role.CanDelete,
		IsAdmin:         role.IsAdmin,
		CreatedAt:       role.CreatedAt,
		UpdatedAt:       role.UpdatedAt,
	}
}

func MapToRole(roleDTO *dto.RoleDTO) *model.Role {
	if roleDTO == nil {
		return nil
	}

	return &model.Role{
		ID:              roleDTO.ID,
		Name:            roleDTO.Name,
		Description:     roleDTO.Description,
		CompanyGlobalID: roleDTO.CompanyGlobalID,
		Permissions:     DTOsMapToPermissions(roleDTO.Permissions),
		CanEdit:         roleDTO.CanEdit,
		CanDelete:       roleDTO.CanDelete,
		IsAdmin:         roleDTO.IsAdmin,
		CreatedAt:       roleDTO.CreatedAt,
		UpdatedAt:       roleDTO.UpdatedAt,
	}
}

func MapCreateToRole(dto *dto.CreateRoleDTO) *model.Role {
	if dto == nil {
		return nil
	}

	permissions := make([]*model.Permission, len(dto.Permissions))
	for i, perm := range dto.Permissions {
		permissions[i] = &model.Permission{
			ID: perm.ID,
		}
	}

	return &model.Role{
		Name:            dto.Name,
		Description:     dto.Description,
		CompanyGlobalID: dto.CompanyGlobalID,
		Permissions:     permissions,
		CanEdit:         dto.CanEdit,
		CanDelete:       dto.CanDelete,
		IsAdmin:         dto.IsAdmin,
	}
}

// MapToRoleDTOs converte um slice de *model.Role para um slice de dto.RoleDTO.
func MapToRoleDTOs(roles []*model.Role) *[]dto.RoleDTO {
	if roles == nil {
		empty := make([]dto.RoleDTO, 0)
		return &empty
	}
	result := make([]dto.RoleDTO, 0, len(roles))
	for _, role := range roles {
		if dto := MapToRoleDTO(role); dto != nil {
			result = append(result, *dto)
		}
	}
	return &result
}
