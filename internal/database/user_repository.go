package database

import (
	"go-sales/internal/model"

	"gorm.io/gorm"
)

// UserRepositoryInterface define os métodos que nosso serviço pode usar para interagir com os dados do usuário.
type UserRepositoryInterface interface {
	FindByEmail(email string) (*model.User, error)
	Create(user *model.User) error
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
	if err := r.db.Where("email_address = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create salva um novo usuário no banco de dados.
func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}
