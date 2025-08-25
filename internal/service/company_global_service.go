package service

import (
	"errors"
	"go-sales/internal/database"
	"go-sales/internal/dto"
	"go-sales/internal/mapper"
	"go-sales/internal/model"
	"go-sales/pkg/util"
	"math"

	"gorm.io/gorm"
)

type CompanyGlobalService struct {
	repo database.CompanyGlobalRepositoryInterface
}
type CompanyGlobalServiceInterface interface {
	Create(companyDTO dto.CreateCompanyGlobalDTO) (*dto.CompanyGlobalDTO, error)
	Update(companyDTO dto.CreateCompanyGlobalDTO, id util.UUID) (*dto.CompanyGlobalDTO, error)
	Delete(id util.UUID) error
	FindByID(id util.UUID) (*dto.CompanyGlobalDTO, error)
	FindByCGC(cgc string) (*dto.CompanyGlobalDTO, error)
	FindAll(filters map[string][]string, page, pageSize int) (*dto.PaginatedResponse[dto.CompanyGlobalDTO], error)
}

func NewCompanyGlobalService(repo database.CompanyGlobalRepositoryInterface) CompanyGlobalServiceInterface {
	return &CompanyGlobalService{
		repo: repo,
	}
}

func (s *CompanyGlobalService) Create(companyDTO dto.CreateCompanyGlobalDTO) (*dto.CompanyGlobalDTO, error) {
	// 1. Verificar se o CGC já existe (lógica de negócio).
	existingCompany, err := s.repo.FindByCGC(companyDTO.CGC)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingCompany != nil {
		return nil, ErrCGCInUse
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

	return mapper.MapToCompanyGlobalDTO(newCompany), nil
}

func (s *CompanyGlobalService) Update(companyDTO dto.CreateCompanyGlobalDTO, id util.UUID) (*dto.CompanyGlobalDTO, error) {
	// 1. Verificar se a empresa existe.
	existingCompany, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
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
			return nil, ErrCGCInUse
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

	return mapper.MapToCompanyGlobalDTO(updatedCompany), nil
}

func (s *CompanyGlobalService) Delete(id util.UUID) error {
	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *CompanyGlobalService) FindByID(id util.UUID) (*dto.CompanyGlobalDTO, error) {
	company, err := s.repo.FindByID(id)
	if err != nil {
		// AQUI ESTÁ A CORREÇÃO:
		// Traduzimos o erro da camada de banco de dados para um erro da camada de serviço.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		// Para qualquer outro erro, nós o repassamos.
		return nil, err
	}
	return mapper.MapToCompanyGlobalDTO(company), nil
}

func (s *CompanyGlobalService) FindByCGC(cgc string) (*dto.CompanyGlobalDTO, error) {
	company, err := s.repo.FindByCGC(cgc)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return mapper.MapToCompanyGlobalDTO(company), nil
}

// ...
func (s *CompanyGlobalService) FindAll(filters map[string][]string, page, pageSize int) (*dto.PaginatedResponse[dto.CompanyGlobalDTO], error) {
	companies, totalItems, err := s.repo.FindAll(filters, page, pageSize)
	if err != nil {
		return nil, err
	}

	// Convert []model.CompanyGlobal to []*model.CompanyGlobal
	companyPtrs := make([]*model.CompanyGlobal, len(companies))
	for i := range companies {
		companyPtrs[i] = &companies[i]
	}
	companyDTOs := mapper.MapToCompanyGlobalDTOs(companyPtrs)

	totalPages := 0
	if pageSize > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(pageSize)))
	}

	// Instancia a struct genérica, especificando que T é dto.CompanyGlobalDTO.
	paginatedResponse := &dto.PaginatedResponse[dto.CompanyGlobalDTO]{
		Items: *companyDTOs,
		PageInfo: dto.PageInfo{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}

	return paginatedResponse, nil
}
