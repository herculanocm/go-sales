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
	log.Info().Msg("Creating a new user")

	create, utilError := GetValidatedDTO[*dto.CreateUserDTO](c, "validatedDTO")
	if utilError != nil {
		log.Error().Err(utilError).Caller().Msg("UserHandler.Create - Error getting validated DTO")
		HandleError(utilError, "UserHandler.Create - Error getting validated DTO", c)
		return
	}

	createdUser, err := h.service.Create(*create)
	if err != nil {
		HandleError(err, "UserHandler.Create error", c)
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
	log.Debug().Msg("Finding all users")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", fmt.Sprintf("%d", h.cfg.AppDefaultAPIPageSize)))
	if err != nil || pageSize < 1 {
		pageSize = h.cfg.AppDefaultAPIPageSize
	}
	filters := c.Request.URL.Query()

	companyGlobalId, err := strconv.ParseInt(c.Query("companyGlobalId"), 10, 64)
	if err != nil {
		customError := service.NewError("invalid companyGlobalId format", http.StatusBadRequest, "invalid_company_global_id_format")
		HandleError(customError, "PermissionHandler.FindAll - Error parsing companyGlobalId", c)
		return
	}
	if companyGlobalId == 0 {
		customError := service.NewError("companyGlobalId is required", http.StatusBadRequest, "company_global_id_required")
		HandleError(customError, "PermissionHandler.FindAll - companyGlobalId is required", c)
		return
	}

	paginatedResult, customErr := h.service.FindAll(filters, page, pageSize, companyGlobalId)
	if customErr != nil {

		HandleError(customErr, "UserHandler.FindAll error", c)
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
