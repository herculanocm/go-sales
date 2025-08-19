package service

import (
	"errors"
	"go-sales/internal/database"
	"go-sales/internal/dto"
	"go-sales/internal/model"
	"go-sales/pkg/util"

	"gorm.io/gorm"
)

type CompanyGlobalService struct {
	repo database.CompanyGlobalRepositoryInterface
}
type CompanyGlobalServiceInterface interface {
	Create(companyDTO dto.CreateCompanyGlobalDTO) (*model.CompanyGlobal, error)
	Update(companyDTO dto.CreateCompanyGlobalDTO, id string) (*model.CompanyGlobal, error)
	Delete(id string) error
	FindByID(id string) (*model.CompanyGlobal, error)
	FindByCGC(cgc string) (*model.CompanyGlobal, error)
	FindAll(filters map[string][]string) ([]model.CompanyGlobal, error)
}

func NewCompanyGlobalService(repo database.CompanyGlobalRepositoryInterface) CompanyGlobalServiceInterface {
	return &CompanyGlobalService{
		repo: repo,
	}
}

func (s *CompanyGlobalService) Create(companyDTO dto.CreateCompanyGlobalDTO) (*model.CompanyGlobal, error) {
	// 1. Verificar se o CGC já existe (lógica de negócio).
	existingCompany, err := s.repo.FindByCGC(companyDTO.CGC)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingCompany != nil {
		return nil, ErrCGCAlreadyExists
	}

	// 2. Mapear o DTO para o modelo do banco de dados.
	newCompany := &model.CompanyGlobal{
		ID:          util.New(),
		Name:        companyDTO.Name,
		Description: companyDTO.Description,
		CGC:         companyDTO.CGC,
		Enabled:     companyDTO.Enabled,
	}

	// 3. Chamar o repositório para persistir a empresa.
	if err := s.repo.Create(newCompany); err != nil {
		return nil, err
	}

	return newCompany, nil
}

func (s *CompanyGlobalService) Update(companyDTO dto.CreateCompanyGlobalDTO, id string) (*model.CompanyGlobal, error) {
	// 1. Verificar se a empresa existe.
	existingCompany, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, EntityNotFound
		}

		return nil, err
	}

	// 2. Verificar se o CGC já existe (lógica de negócio).
	if existingCompany.CGC != companyDTO.CGC {
		otherCompany, err := s.repo.FindByCGC(companyDTO.CGC)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if otherCompany != nil {
			return nil, ErrCGCAlreadyExists
		}
	}

	// 3. Mapear o DTO para o modelo do banco de dados.
	updatedCompany := &model.CompanyGlobal{
		ID:          util.New(),
		Name:        companyDTO.Name,
		Description: companyDTO.Description,
		CGC:         companyDTO.CGC,
		Enabled:     companyDTO.Enabled,
	}

	// 4. Chamar o repositório para persistir a empresa.
	if err := s.repo.Update(updatedCompany); err != nil {
		return nil, err
	}

	return updatedCompany, nil
}

func (s *CompanyGlobalService) Delete(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return EntityNotFound
		}
		return err
	}
	return nil
}

func (s *CompanyGlobalService) FindByID(id string) (*model.CompanyGlobal, error) {
	return s.repo.FindByID(id)
}

func (s *CompanyGlobalService) FindByCGC(cgc string) (*model.CompanyGlobal, error) {
	return s.repo.FindByCGC(cgc)
}
func (s *CompanyGlobalService) FindAll(filters map[string][]string) ([]model.CompanyGlobal, error) {
	// O serviço simplesmente repassa os filtros para o repositório.
	// Lógica de negócio mais complexa sobre os filtros poderia ser adicionada aqui se necessário.
	return s.repo.FindAll(filters)
}
