package database

import (
	"go-sales/internal/model"
	"go-sales/pkg/util"

	"gorm.io/gorm"
)

// RoleRepositoryInterface define os métodos para interagir com os dados de role.
type RoleRepositoryInterface interface {
	FindAllByIDs(roleIDs []int64) ([]*model.Role, error)
	FindByID(id int64) (*model.Role, error)
	FindByName(name string, companyGlobalID int64) (*model.Role, error)
	ExistsByName(name string, companyGlobalID int64) (bool, error)
	Create(role *model.Role) error
	Update(role *model.Role) error
	Delete(id int64) error
	FindAll(filters map[string][]string, page, pageSize int, companyGlobalID int64) ([]model.Role, int64, error)
	AssociatePermissions(role *model.Role, permissions []*model.Permission) error
}

// roleRepository é a implementação concreta que usa o GORM.
type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository cria uma nova instância do repositório de role.
func NewRoleRepository(db *gorm.DB) RoleRepositoryInterface {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindByID(id int64) (*model.Role, error) {
	var role model.Role
	if err := r.db.Preload("Permissions").Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) ExistsByName(name string, companyGlobalID int64) (bool, error) {
	var count int64
	if err := r.db.Model(&model.Role{}).Where("name = ? AND company_global_id = ?", name, companyGlobalID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *roleRepository) FindByName(name string, companyGlobalID int64) (*model.Role, error) {
	var role model.Role
	if err := r.db.Preload("Permissions").Where("name = ? AND company_global_id = ?", name, companyGlobalID).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) Create(role *model.Role) error {
	role.ID = util.NewSnowflake()
	return r.db.Create(role).Error
}

func (r *roleRepository) Update(role *model.Role) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(id int64) error {
	result := r.db.Where("id = ?", id).Delete(&model.Role{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *roleRepository) FindAll(filters map[string][]string, page, pageSize int, companyGlobalID int64) ([]model.Role, int64, error) {
	var roles []model.Role
	var totalItems int64
	query := r.db.Preload("Permissions").Model(&model.Role{}).Where("company_global_id = ?", companyGlobalID)

	allowedFilters := map[string]bool{
		"name":            true,
		"companyGlobalId": true,
	}

	for key, value := range filters {
		if allowed, ok := allowedFilters[key]; ok && allowed && len(value) > 0 {
			if key == "name" {
				query = query.Where("name ILIKE ?", "%"+value[0]+"%")
			} else if key == "companyGlobalId" {
				// companyGlobalId é a chave estrangeira, então usamos o nome da coluna correta.
				query = query.Where("company_global_id = ?", value[0])

			} else {
				query = query.Where(key+" = ?", value[0])
			}
		}
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("Permissions").Offset(offset).Limit(pageSize).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, totalItems, nil
}

// Busca todas as roles pelos IDs informados
func (r *roleRepository) FindAllByIDs(ids []int64) ([]*model.Role, error) {
	var roles []*model.Role
	if err := r.db.Preload("Permissions").Where("id IN ?", ids).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) AssociatePermissions(role *model.Role, permissions []*model.Permission) error {
	return r.db.Model(role).Association("Permissions").Append(permissions)
}
