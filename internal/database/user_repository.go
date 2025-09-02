package database

import (
	"go-sales/internal/model"

	"gorm.io/gorm"
)

// UserRepositoryInterface define os métodos que nosso serviço pode usar para interagir com os dados do usuário.
type UserRepositoryInterface interface {
	FindByEmail(email string) (*model.User, error)
	FindByID(id string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id string) error
	AssociateRoles(user *model.User, roles []*model.Role) error
	FindAll(filters map[string][]string, page, pageSize int) ([]model.User, int64, error)
	EmailExists(email string, company_global_id int64) (bool, error)
}

// userRepository é a implementação concreta que usa o GORM.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository cria uma nova instância do nosso repositório de usuário.
func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &userRepository{db: db}
}

// FindByEmail busca um usuário pelo seu email.
func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	// Usamos First para retornar um erro gorm.ErrRecordNotFound se o usuário não for encontrado.
	if err := r.db.Preload("CompanyGlobal").Preload("Roles").Where("email_address = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) EmailExists(email string, company_global_id int64) (bool, error) {
	var count int64
	if err := r.db.Model(&model.User{}).Where("email_address = ? AND company_global_id = ?", email, company_global_id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Create salva um novo usuário no banco de dados.
func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// Update user
func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	if err := r.db.Preload("CompanyGlobal").Preload("Roles").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Delete(id string) error {
	// Executa a operação de delete e armazena o resultado.
	result := r.db.Where("id = ?", id).Delete(&model.User{})

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

func (r *userRepository) AssociateRoles(user *model.User, roles []*model.Role) error {
	// O método Association("Roles") usa o nome do campo na struct User.
	// Append adiciona as associações na tabela de junção.
	return r.db.Model(user).Association("Roles").Append(roles)
}

func (r *userRepository) FindAll(filters map[string][]string, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var totalItems int64
	query := r.db.Model(&model.User{}).Preload("CompanyGlobal").Preload("Roles")

	allowedFilters := map[string]bool{
		"name":    true,
		"email":   true,
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

	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, totalItems, nil
}
