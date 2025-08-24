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
	if err := r.db.Preload("CompanyGlobal").Where("email_address = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
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
	if err := r.db.Preload("CompanyGlobal").Where("id = ?", id).First(&user).Error; err != nil {
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
