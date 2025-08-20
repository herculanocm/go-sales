package handler

import (
	"errors"
	"go-sales/internal/config"
	"go-sales/internal/dto"
	"go-sales/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler encapsula a dependência do serviço de usuário.
type UserHandler struct {
	service service.UserServiceInterface
	cfg     *config.Config
}

// NewUserHandler cria uma nova instância do nosso handler de usuário.
func NewUserHandler(s service.UserServiceInterface, cfg *config.Config) *UserHandler {
	return &UserHandler{
		service: s,
		cfg:     cfg,
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

func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	validatedDTO, exists := c.Get("validatedDTO")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Validated DTO not found"})
		return
	}

	// Faz a asserção de tipo
	updateUserDTO, ok := validatedDTO.(*dto.CreateUserDTO)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid DTO type"})
		return
	}

	// 2. Chamar a Camada de Serviço
	updatedUser, err := h.service.Update(*updateUserDTO, idStr)
	if err != nil {
		// 3. Tratar Erros da Camada de Serviço
		if errors.Is(err, service.EntityNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		// Para qualquer outro erro não esperado, retorne um erro genérico.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	// 4. Retornar Resposta de Sucesso
	c.JSON(http.StatusOK, updatedUser)
}

func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")

	// 2. Chamar a Camada de Serviço
	if err := h.service.Delete(idStr); err != nil {
		// 3. Tratar Erros da Camada de Serviço
		// Verifica se o erro é o nosso erro de "não encontrado".
		if errors.Is(err, service.EntityNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		// Para qualquer outro erro, consideramos um erro interno do servidor.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	// 4. Retornar Resposta de Sucesso
	c.JSON(http.StatusNoContent, nil)
}
