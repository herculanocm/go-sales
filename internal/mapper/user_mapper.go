package mapper

import (
	"go-sales/internal/dto"
	"go-sales/internal/model"
)

func MapToRoleDTO(role *model.Role) *dto.RoleDTO {
	if role == nil {
		return nil
	}

	return &dto.RoleDTO{
		ID:   role.ID,
		Name: role.Name,
	}
}

func MapToRoleDTOs(roles []*model.Role) *[]dto.RoleDTO {
	if roles == nil {
		empty := make([]dto.RoleDTO, 0)
		return &empty
	}

	roleDTOs := make([]dto.RoleDTO, len(roles))
	for i, role := range roles {
		roleDTOs[i] = *MapToRoleDTO(role)
	}
	return &roleDTOs
}

func MapToUserDTO(user *model.User) *dto.UserDTO {
	if user == nil {
		return nil
	}

	return &dto.UserDTO{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		Password:      user.Password,
		CompanyGlobal: *MapToCompanyGlobalDTO(&user.CompanyGlobal),
		Roles:         *MapToRoleDTOs(user.Roles),
	}
}

func MapToUserDTOs(users []*model.User) *[]dto.UserDTO {
	if users == nil {
		empty := make([]dto.UserDTO, 0)
		return &empty
	}

	userDTOs := make([]dto.UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = *MapToUserDTO(user)
	}
	return &userDTOs
}
