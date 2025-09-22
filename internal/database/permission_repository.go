package database

import (
	"go-sales/internal/model"
	"go-sales/pkg/util"

	"gorm.io/gorm"
)

// PermissionRepositoryInterface define os métodos para interagir com os dados de permission.
type PermissionRepositoryInterface interface {
	FindByID(id int64) (*model.Permission, error)
	FindByName(name string, companyGlobalID int64) (*model.Permission, error)
	ExistsByName(name string, companyGlobalID int64) (bool, error)
	Create(permission *model.Permission) error
	Update(permission *model.Permission) error
	Delete(id int64) error
	FindAll(filters map[string][]string, page, pageSize int, companyGlobalID int64) ([]*model.Permission, int64, error)
	FindByIDs(ids []int64, companyGlobalID *int64) ([]*model.Permission, error)
}

// permissionRepository é a implementação concreta que usa o GORM.
type permissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository cria uma nova instância do repositório de permission.
func NewPermissionRepository(db *gorm.DB) PermissionRepositoryInterface {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) FindByID(id int64) (*model.Permission, error) {
	var permission model.Permission
	if err := r.db.Where("id = ?", id).First(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) ExistsByName(name string, companyGlobalID int64) (bool, error) {
	var count int64
	if err := r.db.Model(&model.Permission{}).Where("name = ? AND company_global_id = ?", name, companyGlobalID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *permissionRepository) FindByName(name string, companyGlobalID int64) (*model.Permission, error) {
	var permission model.Permission
	if err := r.db.Where("name = ? AND company_global_id = ?", name, companyGlobalID).First(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) Create(permission *model.Permission) error {
	permission.ID = util.NewSnowflake()
	return r.db.Create(permission).Error
}

func (r *permissionRepository) Update(permission *model.Permission) error {
	return r.db.Save(permission).Error
}

func (r *permissionRepository) Delete(id int64) error {
	result := r.db.Unscoped().Where("id = ?", id).Delete(&model.Permission{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *permissionRepository) FindAll(filters map[string][]string, page, pageSize int, companyGlobalID int64) ([]*model.Permission, int64, error) {
	var permissions []*model.Permission
	var totalItems int64
	query := r.db.Model(&model.Permission{}).Where("company_global_id = ?", companyGlobalID)

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
	if err := query.Offset(offset).Limit(pageSize).Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	return permissions, totalItems, nil
}

func (r *permissionRepository) FindByIDs(ids []int64, companyGlobalID *int64) ([]*model.Permission, error) {
	query := r.db.Where("id IN ?", ids)
	if companyGlobalID != nil {
		query = query.Where("company_global_id = ?", *companyGlobalID)
	}
	var permissions []*model.Permission
	if err := query.Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}
