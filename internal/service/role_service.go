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

// RoleServiceInterface define a interface para a lógica de negócios de roles.
type RoleServiceInterface interface {
	Create(roleDTO dto.CreateRoleDTO) (*dto.RoleDTO, error)
	Update(roleDTO dto.CreateRoleDTO, roleID int64) (*dto.RoleDTO, error)
	Delete(roleID int64) error
	FindByID(roleID int64) (*dto.RoleDTO, error)
	FindAll(filters map[string][]string, page, pageSize int) (*dto.PaginatedResponse[dto.RoleDTO], error)
}

// roleService é a implementação concreta.
type roleService struct {
	repo        database.RoleRepositoryInterface
	repoPerm    database.PermissionRepositoryInterface
	repoCompany database.CompanyGlobalRepositoryInterface
}

// NewRoleService cria uma nova instância do serviço de roles.
func NewRoleService(repo database.RoleRepositoryInterface, repoPerm database.PermissionRepositoryInterface, repoCompany database.CompanyGlobalRepositoryInterface) RoleServiceInterface {
	return &roleService{repo: repo, repoPerm: repoPerm, repoCompany: repoCompany}
}

func (s *roleService) Create(roleDTO dto.CreateRoleDTO) (*dto.RoleDTO, error) {
	existingRole, err := s.repo.FindByName(roleDTO.Name)
	if err != nil {
		return nil, err
	}
	if existingRole != nil {
		return nil, ErrRoleNameInUse
	}

	company, err := s.repoCompany.FindByID(roleDTO.CompanyGlobalID, false)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, ErrCompanyGlobalNotFound
	}

	permissions, err := s.repoPerm.FindByIDs(roleDTO.PermissionIDs)
	if err != nil {
		return nil, err
	}

	if len(permissions) != len(roleDTO.PermissionIDs) {
		return nil, ErrPermissionNotFound // Um novo erro de serviço que você pode criar
	}

	newRole := &model.Role{
		ID:              util.NewSnowflake(),
		Name:            roleDTO.Name,
		Description:     roleDTO.Description,
		CompanyGlobalID: roleDTO.CompanyGlobalID,
		CanEdit:         roleDTO.CanEdit,
		CanDelete:       roleDTO.CanDelete,
		IsAdmin:         roleDTO.IsAdmin,
	}

	if err := s.repo.Create(newRole); err != nil {
		return nil, err
	}

	if err := s.repo.AssociatePermissions(newRole, permissions); err != nil {
		return nil, err
	}

	return mapper.MapToRoleDTO(newRole), nil
}

func (s *roleService) Update(roleDTO dto.CreateRoleDTO, roleID int64) (*dto.RoleDTO, error) {
	existingRole, err := s.repo.FindByID(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if existingRole == nil {
		return nil, ErrNotFound
	}

	if roleDTO.Name != existingRole.Name {
		roleWithNewName, err := s.repo.FindByName(roleDTO.Name)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if roleWithNewName != nil {
			return nil, ErrRoleNameInUse
		}
	}

	existingRole.Name = roleDTO.Name
	existingRole.Description = roleDTO.Description

	if err := s.repo.Update(existingRole); err != nil {
		return nil, err
	}

	return mapper.MapToRoleDTO(existingRole), nil
}

func (s *roleService) Delete(roleID int64) error {
	err := s.repo.Delete(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *roleService) FindByID(roleID int64) (*dto.RoleDTO, error) {
	role, err := s.repo.FindByID(roleID)
	if err != nil {
		return nil, err
	}
	return mapper.MapToRoleDTO(role), nil
}

func (s *roleService) FindAll(filters map[string][]string, page, pageSize int) (*dto.PaginatedResponse[dto.RoleDTO], error) {
	roles, totalItems, err := s.repo.FindAll(filters, page, pageSize)
	if err != nil {
		return nil, err
	}

	rolePtrs := make([]*model.Role, len(roles))
	for i := range roles {
		rolePtrs[i] = &roles[i]
	}
	roleDTOs := mapper.MapToRoleDTOs(rolePtrs)
	totalPages := 0
	if pageSize > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(pageSize)))
	}

	return &dto.PaginatedResponse[dto.RoleDTO]{
		Items: *roleDTOs,
		PageInfo: dto.PageInfo{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}, nil
}
