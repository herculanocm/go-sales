package mapper

import (
	"go-sales/internal/dto"
	"go-sales/internal/model"
	"time"

	"gorm.io/gorm"
)

func MapToUserDTO(user *model.User) *dto.UserDTO {
	if user == nil {
		return nil
	}

	return &dto.UserDTO{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		EmailRecovery:   user.EmailRecovery,
		EmailVerified:   user.EmailVerified,
		EmailVerifiedAt: user.EmailVerifiedAt,
		Phone:           user.Phone,
		PhoneVerified:   user.PhoneVerified,
		PhoneVerifiedAt: user.PhoneVerifiedAt,
		Password:        nil,
		Enabled:         user.Enabled,
		Actived:         user.Actived,
		ActivationKey:   nil,
		ActivatedAt:     parseTimePtr(user.ActivatedAt),
		ResetKey:        nil,
		ResetRequested:  parseTimePtr(user.ResetRequested),
		ResetAt:         user.ResetAt,
		CompanyGlobalID: user.CompanyGlobalID,
		CompanyGlobal:   *MapToCompanyGlobalDTO(&user.CompanyGlobal),
		Roles:           *MapToRoleDTOs(user.Roles),
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		DeletedAt:       parseDeletedAt(user.DeletedAt),
	}
}

func parseTimePtr(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return nil
	}
	return &t
}

func parseDeletedAt(deleted gorm.DeletedAt) *time.Time {
	if deleted.Valid {
		return &deleted.Time
	}
	return nil
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

func MapCreateUserDTOToUser(userDTO *dto.CreateUserDTO, hashedPassword string, companyGlobal *model.CompanyGlobal, roles []*model.Role) *model.User {
	if userDTO == nil {
		return nil
	}

	return &model.User{
		Name:            userDTO.Name,
		Email:           userDTO.Email,
		Password:        hashedPassword,
		Enabled:         false,
		CompanyGlobalID: userDTO.CompanyGlobalID,
		CompanyGlobal:   *companyGlobal,
		Roles:           roles,
	}
}
