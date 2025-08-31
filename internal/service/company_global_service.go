package service

import (
	"errors"
	"go-sales/internal/database"
	"go-sales/internal/dto"
	"go-sales/internal/mapper"
	"go-sales/internal/model"
	"go-sales/pkg/util"
	"math"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type CompanyGlobalService struct {
	repo database.CompanyGlobalRepositoryInterface
}
type CompanyGlobalServiceInterface interface {
	Create(companyDTO dto.CreateCompanyGlobalDTO) (*dto.CompanyGlobalDTO, ErrorUtil)
	Update(companyDTO dto.CreateCompanyGlobalDTO, id util.UUID) (*dto.CompanyGlobalDTO, ErrorUtil)
	Delete(id util.UUID) ErrorUtil
	FindByID(id util.UUID) (*dto.CompanyGlobalDTO, ErrorUtil)
	FindByCGC(cgc string) (*dto.CompanyGlobalDTO, ErrorUtil)
	FindAll(filters map[string][]string, page, pageSize int) (*dto.PaginatedResponse[dto.CompanyGlobalDTO], ErrorUtil)
	Restore(id util.UUID) ErrorUtil
}

func NewCompanyGlobalService(repo database.CompanyGlobalRepositoryInterface) CompanyGlobalServiceInterface {
	return &CompanyGlobalService{
		repo: repo,
	}
}

func (s *CompanyGlobalService) Restore(id util.UUID) ErrorUtil {
	err := s.repo.Restore(id)
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("failed to restore company")
		return GormDefaultError(err)
	}
	return nil
}

func (s *CompanyGlobalService) Create(companyDTO dto.CreateCompanyGlobalDTO) (*dto.CompanyGlobalDTO, ErrorUtil) {
	// 1. Verificar se o CGC já existe (lógica de negócio).
	existingCompany, err := s.repo.FindByCGC(companyDTO.CGC, true)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().
			Err(err).
			Caller().
			Msg("failed to create company")
		return nil, GormDefaultError(err)
	}
	if existingCompany != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("failed to create company")
		return nil, ErrCGCInUse
	}

	// 2. Mapear o DTO para o modelo do banco de dados.
	newCompany := mapper.MapToCreateCompanyGlobal(&companyDTO)

	// 3. Chamar o repositório para persistir a empresa.
	if err := s.repo.Create(newCompany); err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("failed to create company")
		return nil, GormDefaultError(err)
	}

	return mapper.MapToCompanyGlobalDTO(newCompany), nil
}

func (s *CompanyGlobalService) Update(companyDTO dto.CreateCompanyGlobalDTO, id util.UUID) (*dto.CompanyGlobalDTO, ErrorUtil) {
	// 1. Verificar se a empresa existe.
	original, err := s.repo.FindByID(id, false)
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("failed to update company")
		return nil, GormDefaultError(err)
	}

	otherCompany, err := s.repo.FindByCGC(companyDTO.CGC, true)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().
			Err(err).
			Caller().
			Msg("failed to update company")
		return nil, GormDefaultError(err)
	}
	if otherCompany != nil && *otherCompany.ID != id {
		return nil, ErrCGCInUse
	}

	// 3. Mapear o DTO para o modelo do banco de dados.
	updatedCompany := mapper.MapToUpdateCompanyGlobal(&companyDTO, &id)

	updatedCompany.CreatedAt = original.CreatedAt

	// 4. Chamar o repositório para persistir a empresa.
	if err := s.repo.Update(updatedCompany); err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("failed to update company")

		return nil, GormDefaultError(err)
	}

	return mapper.MapToCompanyGlobalDTO(updatedCompany), nil
}

func (s *CompanyGlobalService) Delete(id util.UUID) ErrorUtil {
	err := s.repo.Delete(id)
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("failed to delete company")
		return GormDefaultError(err)
	}
	return nil
}

func (s *CompanyGlobalService) FindByID(id util.UUID) (*dto.CompanyGlobalDTO, ErrorUtil) {
	company, err := s.repo.FindByID(id, false)
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("failed to find company")
		return nil, GormDefaultError(err)
	}
	return mapper.MapToCompanyGlobalDTO(company), nil
}

func (s *CompanyGlobalService) FindByCGC(cgc string) (*dto.CompanyGlobalDTO, ErrorUtil) {
	company, err := s.repo.FindByCGC(cgc, false)
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("failed to find company")
		return nil, GormDefaultError(err)
	}
	return mapper.MapToCompanyGlobalDTO(company), nil
}

// ...
func (s *CompanyGlobalService) FindAll(filters map[string][]string, page, pageSize int) (*dto.PaginatedResponse[dto.CompanyGlobalDTO], ErrorUtil) {
	companies, totalItems, err := s.repo.FindAll(filters, page, pageSize, false)
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("failed to findAll company")
		return nil, GormDefaultError(err)
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
