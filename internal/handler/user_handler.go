package handler

import (
	"errors"
	"go-sales/internal/dto"
	"go-sales/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler encapsula a dependência do serviço de usuário.
type UserHandler struct {
	service service.UserServiceInterface
}

// NewUserHandler cria uma nova instância do nosso handler de usuário.
func NewUserHandler(s service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		service: s,
	}
}

// Create é o método do handler para criar um novo usuário.
func (h *UserHandler) Create(c *gin.Context) {
	validatedDTO, exists := c.Get("validatedDTO")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Validated DTO not found"})
		return
	}

	// Faz a asserção de tipo
	createUserDTO, ok := validatedDTO.(*dto.CreateUserDTO)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid DTO type"})
		return
	}

	// 2. Chamar a Camada de Serviço
	createdUser, err := h.service.Create(*createUserDTO)
	if err != nil {
		// 3. Tratar Erros da Camada de Serviço
		// Verifica se o erro é um erro de negócio conhecido (email duplicado).
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			// Retorna um erro 409 (Conflict).
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		// Para qualquer outro erro, consideramos um erro interno do servidor.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	// 4. Retornar Resposta de Sucesso
	// Retorna o usuário criado com o status 201 (Created).
	c.JSON(http.StatusCreated, createdUser)
}
