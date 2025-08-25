package mapper

import (
	"go-sales/internal/dto"
	"go-sales/internal/model"
	"time"
)

// MapToRoleDTO converte um model.Role para dto.RoleDTO.

func MapToRoleDTO(role *model.Role) *dto.RoleDTO {
	if role == nil {
		return nil
	}
	var deletedAt *time.Time
	if role.DeletedAt.Valid {
		deletedAt = &role.DeletedAt.Time
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
		DeletedAt:       deletedAt,
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
