package mapper

import (
	"go-sales/internal/dto"
	"go-sales/internal/model"
)

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

	userDTOs := make([]dto.UserDTO, 0, len(users))
	for _, user := range users {
		if dto := MapToUserDTO(user); dto != nil {
			userDTOs = append(userDTOs, *dto)
		}
	}
	return &userDTOs
}
