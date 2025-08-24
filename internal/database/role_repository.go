package database

import (
	"go-sales/internal/model"
	"go-sales/pkg/util"

	"gorm.io/gorm"
)

type RoleRepositoryInterface interface {
	FindAllByIDs(roleIDs []util.UUID) ([]*model.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepositoryInterface {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindAllByIDs(roleIDs []util.UUID) ([]*model.Role, error) {
	var roles []*model.Role
	if err := r.db.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
