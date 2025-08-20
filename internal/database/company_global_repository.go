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
	FindAll(filters map[string][]string, page, pageSize int) ([]model.CompanyGlobal, int64, error)
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
func (r *CompanyGlobalRepository) FindAll(filters map[string][]string, page, pageSize int) ([]model.CompanyGlobal, int64, error) {
	var companies []model.CompanyGlobal
	var totalItems int64
	query := r.db.Model(&model.CompanyGlobal{})

	// Lista de colunas permitidas para filtragem para evitar injeção de SQL.
	allowedFilters := map[string]bool{
		"name":    true,
		"cgc":     true,
		"enabled": true,
	}

	for key, value := range filters {
		// Verifica se o filtro é permitido e se tem um valor.
		if allowed, ok := allowedFilters[key]; ok && allowed && len(value) > 0 {
			// Para campos de texto, usamos 'LIKE' para buscas parciais.
			// Para outros, usamos '=' para buscas exatas.
			if key == "name" {
				query = query.Where("name ILIKE ?", "%"+value[0]+"%") // ILIKE é case-insensitive (PostgreSQL)
			} else {
				query = query.Where(key+" = ?", value[0])
			}
		}
	}

	// 1. Faz a contagem do total de itens que correspondem ao filtro, ANTES de aplicar limit/offset.
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// 2. Calcula o offset para a paginação.
	offset := (page - 1) * pageSize

	// 3. Aplica a paginação (limit e offset) e busca os itens da página.
	if err := query.Offset(offset).Limit(pageSize).Find(&companies).Error; err != nil {
		return nil, 0, err
	}

	return companies, totalItems, nil
}
