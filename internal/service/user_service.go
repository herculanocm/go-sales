package service

import (
	"errors"
	"go-sales/internal/database"
	"go-sales/internal/dto"
	"go-sales/internal/model"
	"go-sales/pkg/util"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Erros de negócio específicos do serviço de usuário.
var (
	ErrEmailAlreadyExists = errors.New("email already in use")
)

// UserService define a interface para a lógica de negócios do usuário.
type UserServiceInterface interface {
	Create(userDTO dto.CreateUserDTO) (*model.User, error)
}

// userService é a implementação concreta.
type userService struct {
	repo database.UserRepositoryInterface
}

// NewUserService cria uma nova instância do serviço de usuário.
// Ele recebe o repositório como uma dependência (Injeção de Dependência).
func NewUserService(repo database.UserRepositoryInterface) UserServiceInterface {
	return &userService{repo: repo}
}

// Create contém a lógica de negócios para criar um novo usuário.
func (s *userService) Create(userDTO dto.CreateUserDTO) (*model.User, error) {
	// 1. Verificar se o email já existe (lógica de negócio).
	existingUser, err := s.repo.FindByEmail(userDTO.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// Um erro inesperado do banco de dados.
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// 2. Hashear a senha (lógica de negócio de segurança).
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. Mapear o DTO para o modelo do banco de dados.
	newUser := &model.User{
		ID:       util.New(), // O GORM pode gerar, mas é bom ser explícito.
		Name:     userDTO.Name,
		Email:    userDTO.Email,
		Password: string(hashedPassword),
	}

	// 4. Chamar o repositório para persistir o usuário.
	if err := s.repo.Create(newUser); err != nil {
		return nil, err
	}

	// 5. Retornar o usuário criado (sem a senha).
	newUser.Password = "" // Nunca retorne o hash da senha.
	return newUser, nil
}
