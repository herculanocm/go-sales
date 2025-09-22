package service

import (
	"errors"
	"fmt"
	"go-sales/internal/database"
	"go-sales/internal/dto"
	"go-sales/internal/mapper"
	"go-sales/internal/model"
	"math"
	"strconv"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceInterface interface {
	Create(userDTO dto.CreateUserDTO) (*dto.UserDTO, ErrorUtil)
	Update(userDTO dto.CreateUserDTO, userID string) (*dto.UserDTO, ErrorUtil)
	Delete(userID string) ErrorUtil
	FindByID(userID string) (*dto.UserDTO, ErrorUtil)
	FindAll(filters map[string][]string, page, pageSize int, companyGlobalID int64) (*dto.PaginatedResponse[dto.UserDTO], ErrorUtil)
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
func (s *userService) Create(userDTO dto.CreateUserDTO) (*dto.UserDTO, ErrorUtil) {
	log.Debug().Msgf("Creating user: %+v", userDTO)

	companyGlobalExists, errCompanyGlobalExists := CheckCompanyGlobalExists(s.repoCompany, userDTO.CompanyGlobalID, false)
	if errCompanyGlobalExists != nil {
		log.Error().
			Err(errCompanyGlobalExists).Str("company_global_id", strconv.FormatInt(userDTO.CompanyGlobalID, 10)).
			Msg("failed to check if company global exists")
		return nil, errCompanyGlobalExists
	}
	if !companyGlobalExists {
		log.Error().
			Err(errCompanyGlobalExists).
			Caller().
			Str("company_global_id", strconv.FormatInt(userDTO.CompanyGlobalID, 10)).
			Msg("failed to find existing company global")
		return nil, ErrCompanyGlobalNotFound
	}

	userEmailExists, errUserEmailExists := CheckUserEmailExists(s.repo, userDTO.Email, userDTO.CompanyGlobalID, false)
	if errUserEmailExists != nil {
		log.Error().
			Err(errUserEmailExists).Str("email", userDTO.Email).
			Str("company_global_id", strconv.FormatInt(userDTO.CompanyGlobalID, 10)).
			Msg("failed to check if user email exists")
		return nil, errUserEmailExists
	}
	if userEmailExists {
		log.Error().
			Err(ErrEmailInUse).
			Caller().
			Str("email", userDTO.Email).
			Str("company_global_id", strconv.FormatInt(userDTO.CompanyGlobalID, 10)).
			Msg("email already in use")
		return nil, ErrEmailInUse
	}

	hashedPassword, errBcryptPass := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if errBcryptPass != nil {
		log.Error().
			Err(errBcryptPass).
			Msg("failed to hash password for new user, email: " + userDTO.Email)
		return nil, ErrInternalServer
	}

	existingCompanyGlobal, err := s.repoCompany.FindByID(userDTO.CompanyGlobalID, false)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().
			Err(err).
			Str("company_global_id", strconv.FormatInt(userDTO.CompanyGlobalID, 10)).
			Msg("failed to find existing company global")
		return nil, ErrDatabase
	}

	roles, err := s.repoRole.FindAllByIDs(userDTO.RoleIDs)
	if err != nil {
		log.Error().
			Err(err).
			Str("role_ids", fmt.Sprintf("%v", userDTO.RoleIDs)).
			Msg("failed to find roles by IDs")
		return nil, ErrDatabase // Erro de banco de dados
	}

	if len(roles) != len(userDTO.RoleIDs) {
		log.Error().
			Str("requested_role_ids", fmt.Sprintf("%v", userDTO.RoleIDs)).
			Str("found_role_count", fmt.Sprintf("%d", len(roles))).
			Msg("some roles not found")
		return nil, ErrRoleNotFound // Um novo erro de serviço que você pode criar
	}

	newUser := mapper.MapCreateUserDTOToUser(&userDTO, string(hashedPassword), existingCompanyGlobal, roles)

	if err := s.repo.Create(newUser); err != nil {
		log.Error().
			Err(err).
			Caller().
			Msg("failed to create user")
		return nil, ErrDatabase
	}

	newUser.Password = "" // Nunca retorne o hash da senha.
	log.Info().Msgf("User created successfully: %+v", newUser)
	return mapper.MapToUserDTO(newUser), nil
}

func (s *userService) Update(userDTO dto.CreateUserDTO, userID string) (*dto.UserDTO, ErrorUtil) {
	// 1. Buscar o usuário existente.

	companyGlobalExists, errCompanyGlobalExists := CheckCompanyGlobalExists(s.repoCompany, userDTO.CompanyGlobalID, false)
	if errCompanyGlobalExists != nil {
		log.Error().
			Err(errCompanyGlobalExists).Str("company_global_id", strconv.FormatInt(userDTO.CompanyGlobalID, 10)).
			Msg("failed to check if company global exists")
		return nil, errCompanyGlobalExists
	}
	if !companyGlobalExists {
		log.Error().
			Err(errCompanyGlobalExists).
			Caller().
			Str("company_global_id", strconv.FormatInt(userDTO.CompanyGlobalID, 10)).
			Msg("failed to find existing company global")
		return nil, ErrCompanyGlobalNotFound
	}

	existingUser, err := s.repo.FindByID(userID)
	if err != nil {
		// AQUI ESTÁ A TRADUÇÃO DO ERRO!
		// Se o repositório retornou "record not found", o serviço retorna "ErrNotFound".
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		// Para qualquer outro erro do banco de dados, apenas repasse.
		return nil, ErrDatabase
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
			return nil, ErrDatabase // Erro de banco de dados
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
		return nil, ErrDatabase
	}

	// 5. Retornar o usuário atualizado.
	existingUser.Password = ""
	return mapper.MapToUserDTO(existingUser), nil
}

func (s *userService) Delete(id string) ErrorUtil {

	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound // Retorne um erro específico se o usuário não for encontrado.
		}
		return ErrDatabase // Retorne outros erros do banco de dados.
	}
	return nil // Retorno nil indica sucesso na exclusão.
}

func (s *userService) FindByID(userID string) (*dto.UserDTO, ErrorUtil) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, ErrDatabase
	}
	return mapper.MapToUserDTO(user), nil
}

func (s *userService) FindAll(filters map[string][]string, page, pageSize int, companyGlobalID int64) (*dto.PaginatedResponse[dto.UserDTO], ErrorUtil) {

	companyExists, errCompanyExists := CheckCompanyGlobalExists(s.repoCompany, companyGlobalID, false)
	if errCompanyExists != nil {
		log.Error().
			Err(errCompanyExists).
			Caller().
			Str("company_global_id", strconv.FormatInt(companyGlobalID, 10)).
			Msg("failed to check if company global exists")
		return nil, errCompanyExists
	}
	if !companyExists {
		return nil, ErrCompanyGlobalNotFound
	}

	// 1. Chamar o repositório para buscar os usuários.
	users, totalItems, err := s.repo.FindAll(filters, page, pageSize)
	if err != nil {
		return nil, ErrDatabase
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
