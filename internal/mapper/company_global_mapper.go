package mapper

import (
	"go-sales/internal/dto"
	"go-sales/internal/model"
	"go-sales/pkg/util"
	"time"
)

func MapToCompanyGlobalAddress(addressDTO *dto.CreateCompanyGlobalAddressDTO) *model.CompanyGlobalAddress {
	if addressDTO == nil {
		return nil
	}
	return &model.CompanyGlobalAddress{
		ID:               nil,
		CompanyID:        nil,
		Street:           addressDTO.Street,
		StreetNumber:     addressDTO.StreetNumber,
		StreetComplement: addressDTO.StreetComplement,
		City:             addressDTO.City,
		State:            addressDTO.State,
		PostalCode:       addressDTO.PostalCode,
		Country:          addressDTO.Country,
	}
}

func MapToCompanyGlobalContact(contactDTO *dto.CreateCompanyGlobalContactDTO) *model.CompanyGlobalContact {
	if contactDTO == nil {
		return nil
	}
	return &model.CompanyGlobalContact{
		ID:        nil,
		CompanyID: nil,
		Name:      contactDTO.Name,
		Email:     contactDTO.Email,
		Phone:     contactDTO.Phone,
		CGC:       contactDTO.CGC,
	}
}

func MapToCompanyGlobalContacts(contactDTOs []*dto.CreateCompanyGlobalContactDTO) []*model.CompanyGlobalContact {
	if contactDTOs == nil {
		return nil
	}
	contacts := make([]*model.CompanyGlobalContact, 0, len(contactDTOs))
	for _, contactDTO := range contactDTOs {
		if contact := MapToCompanyGlobalContact(contactDTO); contact != nil {
			contacts = append(contacts, contact)
		}
	}
	return contacts
}

func MapToCompanyGlobalAddressDTO(address *model.CompanyGlobalAddress) *dto.CompanyGlobalAddressDTO {
	if address == nil {
		return nil
	}
	return &dto.CompanyGlobalAddressDTO{
		ID:               address.ID,
		Street:           address.Street,
		StreetNumber:     address.StreetNumber,
		StreetComplement: address.StreetComplement,
		City:             address.City,
		State:            address.State,
		PostalCode:       address.PostalCode,
		Country:          address.Country,
	}
}

func MapToCompanyGlobalContactDTO(contact *model.CompanyGlobalContact) *dto.CompanyGlobalContactDTO {
	if contact == nil {
		return nil
	}
	return &dto.CompanyGlobalContactDTO{
		ID:    contact.ID,
		Name:  contact.Name,
		Email: contact.Email,
		Phone: contact.Phone,
		CGC:   contact.CGC,
	}
}

func MapToCompanyGlobalContactDTOs(contacts []*model.CompanyGlobalContact) []*dto.CompanyGlobalContactDTO {
	if contacts == nil {
		return nil
	}
	dtos := make([]*dto.CompanyGlobalContactDTO, 0, len(contacts))
	for _, contact := range contacts {
		if dto := MapToCompanyGlobalContactDTO(contact); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}

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
		Email:       company.Email,
		Address:     MapToCompanyGlobalAddressDTO(company.Address),
		Contacts:    MapToCompanyGlobalContactDTOs(company.Contacts),
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
		ID:          nil,
		Name:        companyDTO.Name,
		SocialName:  companyDTO.SocialName,
		Description: companyDTO.Description,
		CGC:         companyDTO.CGC,
		Enabled:     companyDTO.Enabled,
		Email:       companyDTO.Email,
		Address:     MapToCompanyGlobalAddress(companyDTO.Address),
		Contacts:    MapToCompanyGlobalContacts(companyDTO.Contacts),
	}
}

func MapToUpdateCompanyGlobal(companyDTO *dto.CreateCompanyGlobalDTO, id *util.UUID) *model.CompanyGlobal {
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
		Email:       companyDTO.Email,
		Address:     MapToCompanyGlobalAddress(companyDTO.Address),
		Contacts:    MapToCompanyGlobalContacts(companyDTO.Contacts),
	}
}
