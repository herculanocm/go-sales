package service

import (
	"errors"
	"go-sales/internal/database"
	"go-sales/internal/dto"
	"go-sales/internal/mapper"
	"go-sales/internal/model"
	"go-sales/pkg/util"
	"math"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceInterface interface {
	Create(userDTO dto.CreateUserDTO) (*dto.UserDTO, error)
	Update(userDTO dto.CreateUserDTO, userID string) (*dto.UserDTO, error)
	Delete(userID string) error
	FindByID(userID string) (*dto.UserDTO, error)
	FindAll(filters map[string][]string, page, pageSize int) (*dto.PaginatedResponse[dto.UserDTO], error)
}

// userService é a implementação concreta.
type userService struct {
	repo        database.UserRepositoryInterface
	repoCompany database.CompanyGlobalRepositoryInterface
	repoRole    database.RoleRepositoryInterface
}

// NewUserService cria uma nova instância do serviço de usuário.
// Ele recebe o repositório como uma dependência (Injeção de Dependência).
func NewUserService(repo database.UserRepositoryInterface, repoCompany database.CompanyGlobalRepositoryInterface, repoRole database.RoleRepositoryInterface) UserServiceInterface {
	return &userService{repo: repo, repoCompany: repoCompany, repoRole: repoRole}
}

// Create contém a lógica de negócios para criar um novo usuário.
func (s *userService) Create(userDTO dto.CreateUserDTO) (*dto.UserDTO, error) {
	// 1. Verificar se o email já existe (lógica de negócio).
	existingUser, err := s.repo.FindByEmail(userDTO.Email)
	if err != nil {
		// Um erro inesperado do banco de dados.
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailInUse
	}

	// 2. Hashear a senha (lógica de negócio de segurança).
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	existingCompanyGlobal, err := s.repoCompany.FindByID(userDTO.CompanyGlobalID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrCompanyGlobalNotFound
	}
	if existingCompanyGlobal == nil {
		return nil, ErrCompanyGlobalNotFound
	}

	// 1. Buscar as roles que serão associadas.
	roles, err := s.repoRole.FindAllByIDs(userDTO.RoleIDs)
	if err != nil {
		return nil, err // Erro de banco de dados
	}
	// Validação importante: garantir que todas as roles solicitadas foram encontradas.
	if len(roles) != len(userDTO.RoleIDs) {
		return nil, ErrRoleNotFound // Um novo erro de serviço que você pode criar
	}

	// 3. Mapear o DTO para o modelo do banco de dados.
	newUser := &model.User{
		ID:              util.New(),
		Name:            userDTO.Name,
		Email:           userDTO.Email,
		Password:        string(hashedPassword),
		CompanyGlobalID: userDTO.CompanyGlobalID,
		CompanyGlobal:   *existingCompanyGlobal,
	}

	// 4. Chamar o repositório para persistir o usuário.
	if err := s.repo.Create(newUser); err != nil {
		return nil, err
	}

	// 4. Associar as Roles na tabela de junção (user_roles).
	// Esta é a forma correta de fazer a associação.
	if err := s.repo.AssociateRoles(newUser, roles); err != nil {
		// Em um cenário real, você poderia querer deletar o usuário recém-criado
		// para não deixar dados inconsistentes (transação).
		return nil, err
	}

	// 5. Retornar o usuário criado (sem a senha).
	newUser.Password = "" // Nunca retorne o hash da senha.
	return mapper.MapToUserDTO(newUser), nil
}

func (s *userService) Update(userDTO dto.CreateUserDTO, userID string) (*dto.UserDTO, error) {
	// 1. Buscar o usuário existente.
	existingUser, err := s.repo.FindByID(userID)
	if err != nil {
		// AQUI ESTÁ A TRADUÇÃO DO ERRO!
		// Se o repositório retornou "record not found", o serviço retorna "ErrNotFound".
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		// Para qualquer outro erro do banco de dados, apenas repasse.
		return nil, err
	}

	// A verificação de 'existingUser == nil' se torna redundante se o repositório
	// já retorna gorm.ErrRecordNotFound, mas não custa manter por segurança.
	if existingUser == nil {
		return nil, ErrNotFound
	}

	// 2. Verificar se o novo email já está em uso por OUTRO usuário.
	if userDTO.Email != existingUser.Email {
		userWithNewEmail, err := s.repo.FindByEmail(userDTO.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err // Erro de banco de dados
		}
		if userWithNewEmail != nil {
			return nil, ErrEmailInUse
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
	return mapper.MapToUserDTO(existingUser), nil
}

func (s *userService) Delete(id string) error {

	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound // Retorne um erro específico se o usuário não for encontrado.
		}
		return err // Retorne outros erros do banco de dados.
	}
	return nil // Retorno nil indica sucesso na exclusão.
}

func (s *userService) FindByID(userID string) (*dto.UserDTO, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	return mapper.MapToUserDTO(user), nil
}

func (s *userService) FindAll(filters map[string][]string, page, pageSize int) (*dto.PaginatedResponse[dto.UserDTO], error) {
	// 1. Chamar o repositório para buscar os usuários.
	users, totalItems, err := s.repo.FindAll(filters, page, pageSize)
	if err != nil {
		return nil, err
	}

	userPtrs := make([]*model.User, len(users))
	for i := range users {
		userPtrs[i] = &users[i]
	}
	userDTOs := mapper.MapToUserDTOs(userPtrs)
	totalPages := 0
	if pageSize > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(pageSize)))
	}

	return &dto.PaginatedResponse[dto.UserDTO]{
		Items: *userDTOs,
		PageInfo: dto.PageInfo{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}, nil
}
