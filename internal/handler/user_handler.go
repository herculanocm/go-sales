package handler

import (
	"fmt"
	"go-sales/internal/config"
	"go-sales/internal/dto"
	"go-sales/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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
		log.Error().
			Str("error_code", err.Code()). // imprime o código do erro, que você define em cada erro
			Err(err).
			Msg("UserHandler.Create error")
		c.JSON(err.HTTPStatusCode(), gin.H{"error": err.Error(), "code": err.Code()})
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
		c.JSON(err.HTTPStatusCode(), gin.H{"error": err.Error(), "code": err.Code()})
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
		if err != nil {
			c.JSON(err.HTTPStatusCode(), gin.H{"error": err.Error(), "code": err.Code()})
			return
		}

		// Para qualquer outro erro, consideramos um erro interno do servidor.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	// 4. Retornar Resposta de Sucesso
	c.JSON(http.StatusNoContent, nil)
}

func (h *UserHandler) FindAll(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", fmt.Sprintf("%d", h.cfg.AppDefaultAPIPageSize)))
	if err != nil || pageSize < 1 {
		pageSize = h.cfg.AppDefaultAPIPageSize
	}

	// 2. Extrai os filtros.
	filters := c.Request.URL.Query()

	// 2. Chamar a Camada de Serviço, passando os filtros.
	// 3. Chama o serviço.
	paginatedResult, err := h.service.FindAll(filters, page, pageSize)
	if customErr, ok := err.(service.ErrorUtil); ok {
		c.JSON(customErr.HTTPStatusCode(), gin.H{"error": customErr.Error(), "code": customErr.Code()})
		return
	}

	// 4. Retorna o resultado paginado.
	c.JSON(http.StatusOK, paginatedResult)
}

func (h *UserHandler) FindByID(c *gin.Context) {
	idStr := c.Param("id")

	// 2. Chamar a Camada de Serviço
	user, err := h.service.FindByID(idStr)
	if err != nil {
		// 3. Tratar Erros da Camada de Serviço
		if err != nil {
			c.JSON(err.HTTPStatusCode(), gin.H{"error": err.Error(), "code": err.Code()})
			return
		}

		// Para qualquer outro erro, consideramos um erro interno do servidor.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	// 4. Retornar Resposta de Sucesso
	c.JSON(http.StatusOK, user)
}
