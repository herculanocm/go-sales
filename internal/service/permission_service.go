package service

import (
	"errors"
	"go-sales/internal/database"
	"go-sales/internal/dto"
	"go-sales/internal/mapper"
	"go-sales/internal/model"
	"math"
	"strings"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// PermissionServiceInterface define a interface para a lógica de negócios de permissões.
type PermissionServiceInterface interface {
	Create(permissionDTO dto.CreatePermissionDTO) (*dto.PermissionDTO, ErrorUtil)
	Update(permissionDTO dto.CreatePermissionDTO, permissionID int64) (*dto.PermissionDTO, ErrorUtil)
	Delete(permissionID int64) ErrorUtil
	FindByID(permissionID int64) (*dto.PermissionDTO, ErrorUtil)
	FindAll(filters map[string][]string, page, pageSize int) (*dto.PaginatedResponse[dto.PermissionDTO], ErrorUtil)
}

// permissionService é a implementação concreta.
type permissionService struct {
	repo              database.PermissionRepositoryInterface
	repoCompanyGlobal database.CompanyGlobalRepositoryInterface
}

// NewPermissionService cria uma nova instância do serviço de permissões.
func NewPermissionService(repo database.PermissionRepositoryInterface, repoCompanyGlobal database.CompanyGlobalRepositoryInterface) PermissionServiceInterface {
	return &permissionService{repo: repo, repoCompanyGlobal: repoCompanyGlobal}
}

func (s *permissionService) Create(permissionDTO dto.CreatePermissionDTO) (*dto.PermissionDTO, ErrorUtil) {
	// Verificar se já existe uma permissão com o mesmo nome
	permissionDTO.Name = strings.ToUpper(strings.TrimSpace(permissionDTO.Name))

	companyGlobalExists, err := s.repoCompanyGlobal.Exists(permissionDTO.CompanyGlobalID, false)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().
			Err(err).
			Caller().
			Str("company_global_id", string(permissionDTO.CompanyGlobalID)).
			Msg("failed to check if company global exists")
		return nil, GormDefaultError(err)
	}
	if !companyGlobalExists {
		log.Error().
			Err(err).
			Caller().
			Str("company_global_id", string(permissionDTO.CompanyGlobalID)).
			Msg("failed to find existing company global")
		return nil, ErrCompanyGlobalNotFound
	}

	existingPermission, err := s.repo.FindByName(permissionDTO.Name, permissionDTO.CompanyGlobalID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().
			Err(err).
			Caller().
			Str("permission_name", permissionDTO.Name).
			Str("company_global_id", string(permissionDTO.CompanyGlobalID)).
			Msg("failed to find existing permission")
		return nil, GormDefaultError(err)
	}
	if existingPermission != nil {
		log.Error().
			Err(err).
			Caller().
			Str("permission_name", permissionDTO.Name).
			Str("company_global_id", string(permissionDTO.CompanyGlobalID)).
			Msg("permission name is already in use")
		return nil, ErrPermissionNameInUse
	}

	newPermission := mapper.MapToPermission(&permissionDTO)

	if err := s.repo.Create(newPermission); err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("failed to create permission")
		return nil, GormDefaultError(err)
	}

	return mapper.MapToPermissionDTO(newPermission), nil
}

func (s *permissionService) Update(permissionDTO dto.CreatePermissionDTO, permissionID int64) (*dto.PermissionDTO, ErrorUtil) {
	permissionDTO.Name = strings.ToUpper(strings.TrimSpace(permissionDTO.Name))
	existingPermission, err := s.repo.FindByID(permissionID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().
			Err(err).
			Caller().
			Str("permission_id", string(permissionID)).
			Msg("failed to find existing permission")
		return nil, GormDefaultError(err)
	}
	if existingPermission == nil {
		log.Error().
			Err(err).
			Caller().
			Str("permission_id", string(permissionID)).
			Msg("permission not found")
		return nil, ErrNotFound
	}

	existingPermissionName, err := s.repo.FindByName(permissionDTO.Name, permissionDTO.CompanyGlobalID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().
			Err(err).
			Caller().
			Str("permission_name", permissionDTO.Name).
			Str("company_global_id", string(permissionDTO.CompanyGlobalID)).
			Msg("failed to find existing permission")
		return nil, GormDefaultError(err)
	}
	if existingPermission != nil && existingPermissionName.ID != permissionID {
		log.Error().
			Err(err).
			Caller().
			Str("permission_name", permissionDTO.Name).
			Str("company_global_id", string(permissionDTO.CompanyGlobalID)).
			Msg("permission name is already in use")
		return nil, ErrPermissionNameInUse
	}

	existingPermission.Name = permissionDTO.Name
	existingPermission.Description = permissionDTO.Description

	if err := s.repo.Update(existingPermission); err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("failed to update permission")
		return nil, GormDefaultError(err)
	}

	return mapper.MapToPermissionDTO(existingPermission), nil
}

func (s *permissionService) Delete(permissionID int64) ErrorUtil {
	err := s.repo.Delete(permissionID)
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("failed to delete permission")
		return GormDefaultError(err)
	}
	return nil
}

func (s *permissionService) FindByID(permissionID int64) (*dto.PermissionDTO, ErrorUtil) {
	permission, err := s.repo.FindByID(permissionID)
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("failed to find permission by ID")
		return nil, GormDefaultError(err)
	}
	return mapper.MapToPermissionDTO(permission), nil
}

func (s *permissionService) FindAll(filters map[string][]string, page, pageSize int) (*dto.PaginatedResponse[dto.PermissionDTO], ErrorUtil) {
	permissions, totalItems, err := s.repo.FindAll(filters, page, pageSize)
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("failed to findAll permission")
		return nil, GormDefaultError(err)
	}

	permissionPtrs := make([]*model.Permission, len(permissions))
	copy(permissionPtrs, permissions)
	permissionDTOs := mapper.MapToPermissionDTOs(permissionPtrs)
	totalPages := 0
	if pageSize > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(pageSize)))
	}

	return &dto.PaginatedResponse[dto.PermissionDTO]{
		Items: *permissionDTOs,
		PageInfo: dto.PageInfo{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}, nil
}
