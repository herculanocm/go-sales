package database

import (
	"go-sales/internal/model"

	"gorm.io/gorm"
)

type CompanyGlobalRepository struct {
	db *gorm.DB
}

type CompanyGlobalRepositoryInterface interface {
	Create(company *model.CompanyGlobal) error
	FindByID(id string) (*model.CompanyGlobal, error)
	FindByCGC(cgc string) (*model.CompanyGlobal, error)
	Update(company *model.CompanyGlobal) error
	Delete(id string) error
}

func NewCompanyGlobalRepository(db *gorm.DB) CompanyGlobalRepositoryInterface {
	return &CompanyGlobalRepository{
		db: db,
	}
}

func (r *CompanyGlobalRepository) Create(company *model.CompanyGlobal) error {
	return r.db.Create(company).Error
}

func (r *CompanyGlobalRepository) FindByID(id string) (*model.CompanyGlobal, error) {
	var company model.CompanyGlobal
	if err := r.db.Where("id = ?", id).First(&company).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *CompanyGlobalRepository) FindByCGC(cgc string) (*model.CompanyGlobal, error) {
	var company model.CompanyGlobal
	if err := r.db.Where("cgc = ?", cgc).First(&company).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *CompanyGlobalRepository) Update(company *model.CompanyGlobal) error {
	return r.db.Save(company).Error
}

func (r *CompanyGlobalRepository) Delete(id string) error {
	// Executa a operação de delete e armazena o resultado.
	result := r.db.Where("id = ?", id).Delete(&model.CompanyGlobal{})

	// Primeiro, verifica se houve um erro real na execução da query.
	if result.Error != nil {
		return result.Error
	}

	// Se não houve erro, verifica se alguma linha foi realmente afetada.
	// Se RowsAffected for 0, significa que o registro não foi encontrado.
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
