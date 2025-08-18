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

// UserService define a interface para a lógica de negócios do usuário.
type UserServiceInterface interface {
	Create(userDTO dto.CreateUserDTO) (*model.User, error)
	Update(userDTO dto.CreateUserDTO, userID string) (*model.User, error)
	Delete(userID string) error
	FindByID(userID string) (*model.User, error)
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

func (s *userService) Update(userDTO dto.CreateUserDTO, userID string) (*model.User, error) {
	// 1. Buscar o usuário existente.
	existingUser, err := s.repo.FindByID(userID)
	if err != nil {
		// AQUI ESTÁ A TRADUÇÃO DO ERRO!
		// Se o repositório retornou "record not found", o serviço retorna "EntityNotFound".
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, EntityNotFound
		}
		// Para qualquer outro erro do banco de dados, apenas repasse.
		return nil, err
	}

	// A verificação de 'existingUser == nil' se torna redundante se o repositório
	// já retorna gorm.ErrRecordNotFound, mas não custa manter por segurança.
	if existingUser == nil {
		return nil, EntityNotFound
	}

	// 2. Verificar se o novo email já está em uso por OUTRO usuário.
	if userDTO.Email != existingUser.Email {
		userWithNewEmail, err := s.repo.FindByEmail(userDTO.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err // Erro de banco de dados
		}
		if userWithNewEmail != nil {
			return nil, ErrEmailAlreadyExists
		}
	}

	// 3. Atualizar os campos do usuário.
	existingUser.Name = userDTO.Name
	existingUser.Email = userDTO.Email // Atualiza o email
	if userDTO.Password != "" {
		// ... (lógica de hashear a senha) ...
	}

	// 4. Chamar o repositório para persistir as alterações.
	if err := s.repo.Update(existingUser); err != nil {
		return nil, err
	}

	// 5. Retornar o usuário atualizado.
	existingUser.Password = ""
	return existingUser, nil
}

func (s *userService) Delete(id string) error {

	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return EntityNotFound // Retorne um erro específico se o usuário não for encontrado.
		}
		return err // Retorne outros erros do banco de dados.
	}
	return nil // Retorno nil indica sucesso na exclusão.
}

func (s *userService) FindByID(userID string) (*model.User, error) {
	return s.repo.FindByID(userID)
}
