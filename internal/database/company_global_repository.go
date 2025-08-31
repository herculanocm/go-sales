package database

import (
	"go-sales/internal/model"
	"go-sales/pkg/util"

	"gorm.io/gorm"
)

type CompanyGlobalRepository struct {
	db *gorm.DB
}

type CompanyGlobalRepositoryInterface interface {
	Create(company *model.CompanyGlobal) error
	FindByID(id util.UUID, useUnscoped bool) (*model.CompanyGlobal, error)
	FindByCGC(cgc string, useUnscoped bool) (*model.CompanyGlobal, error)
	Update(company *model.CompanyGlobal) error
	Delete(id util.UUID) error
	FindAll(filters map[string][]string, page, pageSize int, useUnscoped bool) ([]model.CompanyGlobal, int64, error)
	Restore(id util.UUID) error
}

func NewCompanyGlobalRepository(db *gorm.DB) CompanyGlobalRepositoryInterface {
	return &CompanyGlobalRepository{
		db: db,
	}
}

func (r *CompanyGlobalRepository) Create(company *model.CompanyGlobal) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Garantir que os CompanyID dos filhos estejam definidos para que o GORM
		// possa persistir as associações ao criar a company (evita inserts duplicados).

		company.ID = util.NewPtr()

		if company.Address != nil {
			company.Address.ID = util.NewPtr()
			company.Address.CompanyID = company.ID
		}

		for _, contact := range company.Contacts {
			if contact == nil {
				continue
			}
			contact.ID = util.NewPtr()
			contact.CompanyID = company.ID
		}

		// Cria a company (o GORM irá criar as associações já presentes na struct).
		if err := tx.Create(company).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *CompanyGlobalRepository) FindByID(id util.UUID, useUnscoped bool) (*model.CompanyGlobal, error) {
	var company model.CompanyGlobal
	dbQuery := r.db
	if useUnscoped {
		dbQuery = dbQuery.Unscoped()
	}
	if err := dbQuery.Preload("Address").Preload("Contacts").Where("id = ?", id).First(&company).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *CompanyGlobalRepository) FindByCGC(cgc string, useUnscoped bool) (*model.CompanyGlobal, error) {
	var company model.CompanyGlobal
	dbQuery := r.db
	if useUnscoped {
		dbQuery = dbQuery.Unscoped()
	}
	if err := dbQuery.Preload("Address").Preload("Contacts").Where("cgc = ?", cgc).First(&company).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *CompanyGlobalRepository) Update(company *model.CompanyGlobal) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Carrega company e associações atuais
		var existing model.CompanyGlobal
		if err := tx.Preload("Address").Preload("Contacts").Where("id = ?", company.ID).First(&existing).Error; err != nil {
			return err
		}

		if company.Address != nil {
			company.Address.ID = util.NewPtr()
			company.Address.CompanyID = company.ID
		}

		for _, contact := range company.Contacts {
			if contact == nil {
				continue
			}
			contact.ID = util.NewPtr()
			contact.CompanyID = company.ID
		}

		// Atualiza os campos da tabela principal
		if err := tx.Model(&existing).Select(
			"name", "social_name", "description", "cgc", "email", "enabled",
		).Updates(company).Error; err != nil {
			return err
		}

		// Delete all addresses (hard delete) to avoid GORM setting FK to NULL
		if err := tx.Unscoped().Where("company_id = ?", company.ID).Delete(&model.CompanyGlobalAddress{}).Error; err != nil {
			return err
		}
		// Recreate address if veio no payload
		if company.Address != nil {
			if err := tx.Create(company.Address).Error; err != nil {
				return err
			}
		}

		// Delete all contacts (hard delete) and recreate from payload
		if err := tx.Unscoped().Where("company_id = ?", company.ID).Delete(&model.CompanyGlobalContact{}).Error; err != nil {
			return err
		}
		for _, contact := range company.Contacts {
			if contact == nil {
				continue
			}
			if err := tx.Create(contact).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// ...existing code...

func (r *CompanyGlobalRepository) Delete(id util.UUID) error {
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
func (r *CompanyGlobalRepository) FindAll(filters map[string][]string, page, pageSize int, useUnscoped bool) ([]model.CompanyGlobal, int64, error) {
	var companies []model.CompanyGlobal
	var totalItems int64
	query := r.db.Model(&model.CompanyGlobal{})

	if useUnscoped {
		query = query.Unscoped()
	}

	// Lista de colunas permitidas para filtragem para evitar injeção de SQL.
	allowedFilters := map[string]bool{
		"name":        true,
		"cgc":         true,
		"enabled":     true,
		"email":       true,
		"social_name": true,
	}

	for key, value := range filters {
		// Verifica se o filtro é permitido e se tem um valor.
		if allowed, ok := allowedFilters[key]; ok && allowed && len(value) > 0 {
			// Para campos de texto, usamos 'LIKE' para buscas parciais.
			// Para outros, usamos '=' para buscas exatas.
			if key == "name" {
				query = query.Where("name ILIKE ?", "%"+value[0]+"%") // ILIKE é case-insensitive (PostgreSQL)
			} else if key == "social_name" {
				query = query.Where("social_name ILIKE ?", "%"+value[0]+"%")
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
	if err := query.Preload("Address").Preload("Contacts").Offset(offset).Limit(pageSize).Find(&companies).Error; err != nil {
		return nil, 0, err
	}

	return companies, totalItems, nil
}

func (r *CompanyGlobalRepository) Restore(id util.UUID) error {
	result := r.db.Model(&model.CompanyGlobal{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
