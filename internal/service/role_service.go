package service

import (
	"errors"
	"go-sales/internal/database"
	"go-sales/internal/dto"
	"go-sales/internal/mapper"
	"go-sales/internal/model"
	"math"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// RoleServiceInterface define a interface para a lógica de negócios de roles.
type RoleServiceInterface interface {
	Create(roleDTO dto.CreateRoleDTO) (*dto.RoleDTO, ErrorUtil)
	Update(roleDTO dto.RoleDTO, roleID int64) (*dto.RoleDTO, ErrorUtil)
	Delete(roleID int64) error
	FindByID(roleID int64) (*dto.RoleDTO, error)
	FindAll(filters map[string][]string, page, pageSize int, companyGlobalID int64) (*dto.PaginatedResponse[dto.RoleDTO], error)
}

// roleService é a implementação concreta.
type roleService struct {
	repo              database.RoleRepositoryInterface
	repoPerm          database.PermissionRepositoryInterface
	repoCompanyGlobal database.CompanyGlobalRepositoryInterface
}

// NewRoleService cria uma nova instância do serviço de roles.
func NewRoleService(repo database.RoleRepositoryInterface, repoPerm database.PermissionRepositoryInterface, repoCompanyGlobal database.CompanyGlobalRepositoryInterface) RoleServiceInterface {
	return &roleService{repo: repo, repoPerm: repoPerm, repoCompanyGlobal: repoCompanyGlobal}
}

func (s *roleService) Create(roleDTO dto.CreateRoleDTO) (*dto.RoleDTO, ErrorUtil) {
	roleDTO.Name = strings.ToUpper(strings.TrimSpace(roleDTO.Name))
	companyGlobalExists, err := s.repoCompanyGlobal.Exists(roleDTO.CompanyGlobalID, false)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().
			Err(err).
			Str("company_global_id", strconv.FormatInt(roleDTO.CompanyGlobalID, 10)).
			Str("company_global_id", strconv.FormatInt(roleDTO.CompanyGlobalID, 10)).
			Msg("failed to check if company global exists")
		return nil, GormDefaultError(err)
	}
	if !companyGlobalExists {
		log.Error().
			Err(err).
			Caller().
			Str("company_global_id", strconv.FormatInt(roleDTO.CompanyGlobalID, 10)).
			Msg("failed to find existing company global")
		return nil, ErrCompanyGlobalNotFound
	}

	existingRole, err := s.repo.ExistsByName(roleDTO.Name, roleDTO.CompanyGlobalID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().
			Err(err).
			Caller().
			Str("role_name", roleDTO.Name).
			Msg("failed to find existing role")
		return nil, GormDefaultError(err)
	}
	if existingRole {
		return nil, ErrRoleNameInUse
	}

	if len(roleDTO.Permissions) == 0 {
		return nil, ErrRoleMustHavePermissions
	}

	permIds := make([]int64, len(roleDTO.Permissions))
	for i, perm := range roleDTO.Permissions {
		permIds[i] = perm.ID
	}
	permissions, err := s.repoPerm.FindByIDs(permIds, &roleDTO.CompanyGlobalID) // Ensure permissions belong to the same company global
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Str("role_name", roleDTO.Name).
			Msg("failed to find role permissions")
		return nil, GormDefaultError(err)
	}

	if len(permissions) == 0 {
		log.Error().
			Err(err).
			Caller().
			Str("role_name", roleDTO.Name).
			Msg("failed to find role permissions")
		return nil, ErrPermissionsNotFound
	}

	newRole := mapper.MapCreateToRole(&roleDTO)
	newRole.Permissions = permissions

	if err := s.repo.Create(newRole); err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("failed to create new role")
		return nil, GormDefaultError(err)
	}

	return mapper.MapToRoleDTO(newRole), nil
}

func (s *roleService) Update(roleDTO dto.RoleDTO, roleID int64) (*dto.RoleDTO, ErrorUtil) {
	roleDTO.Name = strings.ToUpper(strings.TrimSpace(roleDTO.Name))
	companyGlobalExists, err := s.repoCompanyGlobal.Exists(roleDTO.CompanyGlobalID, false)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().
			Err(err).
			Str("company_global_id", strconv.FormatInt(roleDTO.CompanyGlobalID, 10)).
			Str("company_global_id", strconv.FormatInt(roleDTO.CompanyGlobalID, 10)).
			Msg("failed to check if company global exists")
		return nil, GormDefaultError(err)
	}
	if !companyGlobalExists {
		log.Error().
			Err(err).
			Caller().
			Str("company_global_id", strconv.FormatInt(roleDTO.CompanyGlobalID, 10)).
			Msg("failed to find existing company global")
		return nil, ErrCompanyGlobalNotFound
	}

	existingRole, err := s.repo.FindByID(roleID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().
			Err(err).
			Caller().
			Str("role_id", strconv.FormatInt(roleID, 10)).
			Msg("failed to find existing role")
		return nil, GormDefaultError(err)
	}
	if existingRole == nil {
		log.Error().
			Caller().
			Str("role_id", strconv.FormatInt(roleID, 10)).
			Msg("role not found")
		return nil, ErrNotFound
	}

	if roleDTO.Name != existingRole.Name {
		roleWithNewName, err := s.repo.ExistsByName(roleDTO.Name, roleDTO.CompanyGlobalID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().
				Err(err).
				Caller().
				Str("role_name", roleDTO.Name).
				Msg("failed to find existing role")
			return nil, GormDefaultError(err)
		}
		if roleWithNewName {
			log.Error().
				Caller().
				Str("role_name", roleDTO.Name).
				Msg("role name is already in use")
			return nil, ErrRoleNameInUse
		}
	}

	if len(roleDTO.Permissions) == 0 {
		return nil, ErrRoleMustHavePermissions
	}

	permIds := make([]int64, len(roleDTO.Permissions))
	for i, perm := range roleDTO.Permissions {
		permIds[i] = perm.ID
	}
	permissions, err := s.repoPerm.FindByIDs(permIds, &roleDTO.CompanyGlobalID) // Ensure permissions belong to the same company global
	if err != nil {
		log.Error().
			Err(err).
			Caller().
			Str("role_name", roleDTO.Name).
			Msg("failed to find role permissions")
		return nil, GormDefaultError(err)
	}

	if len(permissions) == 0 {
		log.Error().
			Err(err).
			Caller().
			Str("role_name", roleDTO.Name).
			Msg("failed to find role permissions")
		return nil, ErrPermissionsNotFound
	}

	updateRole := mapper.MapToRole(&roleDTO)

	if err := s.repo.Update(updateRole); err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("failed to update role")
		return nil, GormDefaultError(err)
	}

	return mapper.MapToRoleDTO(updateRole), nil
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

func (s *roleService) FindAll(filters map[string][]string, page, pageSize int, companyGlobalID int64) (*dto.PaginatedResponse[dto.RoleDTO], error) {
	roles, totalItems, err := s.repo.FindAll(filters, page, pageSize, companyGlobalID)
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
