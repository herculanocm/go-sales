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

// PermissionServiceInterface define a interface para a lógica de negócios de permissões.
type PermissionServiceInterface interface {
	Create(permissionDTO dto.CreatePermissionDTO) (*dto.PermissionDTO, error)
	Update(permissionDTO dto.CreatePermissionDTO, permissionID string) (*dto.PermissionDTO, error)
	Delete(permissionID string) error
	FindByID(permissionID string) (*dto.PermissionDTO, error)
	FindAll(filters map[string][]string, page, pageSize int) (*dto.PaginatedResponse[dto.PermissionDTO], error)
}

// permissionService é a implementação concreta.
type permissionService struct {
	repo database.PermissionRepositoryInterface
}

// NewPermissionService cria uma nova instância do serviço de permissões.
func NewPermissionService(repo database.PermissionRepositoryInterface) PermissionServiceInterface {
	return &permissionService{repo: repo}
}

func (s *permissionService) Create(permissionDTO dto.CreatePermissionDTO) (*dto.PermissionDTO, error) {
	// Verificar se já existe uma permissão com o mesmo nome
	existingPermission, err := s.repo.FindByName(permissionDTO.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingPermission != nil {
		return nil, ErrPermissionNameInUse
	}

	newPermission := &model.Permission{
		ID:          util.New(),
		Name:        permissionDTO.Name,
		Description: permissionDTO.Description,
	}

	if err := s.repo.Create(newPermission); err != nil {
		return nil, err
	}

	return mapper.MapToPermissionDTO(newPermission), nil
}

func (s *permissionService) Update(permissionDTO dto.CreatePermissionDTO, permissionID string) (*dto.PermissionDTO, error) {
	existingPermission, err := s.repo.FindByID(permissionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if existingPermission == nil {
		return nil, ErrNotFound
	}

	// Verificar se o novo nome já está em uso por outra permissão
	if permissionDTO.Name != existingPermission.Name {
		permissionWithNewName, err := s.repo.FindByName(permissionDTO.Name)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if permissionWithNewName != nil {
			return nil, ErrPermissionNameInUse
		}
	}

	existingPermission.Name = permissionDTO.Name
	existingPermission.Description = permissionDTO.Description

	if err := s.repo.Update(existingPermission); err != nil {
		return nil, err
	}

	return mapper.MapToPermissionDTO(existingPermission), nil
}

func (s *permissionService) Delete(permissionID string) error {
	err := s.repo.Delete(permissionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *permissionService) FindByID(permissionID string) (*dto.PermissionDTO, error) {
	permission, err := s.repo.FindByID(permissionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return mapper.MapToPermissionDTO(permission), nil
}

func (s *permissionService) FindAll(filters map[string][]string, page, pageSize int) (*dto.PaginatedResponse[dto.PermissionDTO], error) {
	permissions, totalItems, err := s.repo.FindAll(filters, page, pageSize)
	if err != nil {
		return nil, err
	}

	permissionPtrs := make([]*model.Permission, len(permissions))
	for i := range permissions {
		permissionPtrs[i] = permissions[i]
	}
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
