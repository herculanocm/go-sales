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

// PermissionHandler encapsula a dependência do serviço de permissões.
type PermissionHandler struct {
	service service.PermissionServiceInterface
	cfg     *config.Config
}

// NewPermissionHandler cria uma nova instância do handler de permissões.
func NewPermissionHandler(s service.PermissionServiceInterface, cfg *config.Config) *PermissionHandler {
	return &PermissionHandler{
		service: s,
		cfg:     cfg,
	}
}

// Create cria uma nova permissão.
func (h *PermissionHandler) Create(c *gin.Context) {
	log.Info().Msg("Creating a new permission")
	validatedDTO, exists := c.Get("validatedDTO")
	if !exists {
		log.Error().Msg("Permission.Create - Validated DTO not found")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Validated DTO not found"})
		return
	}

	createPermissionDTO, ok := validatedDTO.(*dto.CreatePermissionDTO)
	if !ok {
		log.Error().Msg("Permission.Create - Invalid DTO type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid DTO type"})
		return
	}

	createdPermission, err := h.service.Create(*createPermissionDTO)
	if err != nil {
		log.Error().
			Str("error_code", err.Code()). // imprime o código do erro, que você define em cada erro
			Err(err).
			Msg("UserHandler.Create error")
		c.JSON(err.HTTPStatusCode(), gin.H{"error": err.Error(), "code": err.Code()})
		return
	}

	c.JSON(http.StatusCreated, createdPermission)
}

// Update atualiza uma permissão existente.
func (h *PermissionHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	log.Info().Str("permission_id", idStr).Msg("Updating permission")
	validatedDTO, exists := c.Get("validatedDTO")
	if !exists {
		log.Error().Msg("Permission.Update - Validated DTO not found")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Validated DTO not found"})
		return
	}

	updatePermissionDTO, ok := validatedDTO.(*dto.CreatePermissionDTO)
	if !ok {
		log.Error().Msg("Permission.Update - Invalid DTO type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid DTO type"})
		return
	}

	updatedPermission, err := h.service.Update(*updatePermissionDTO, idStr)
	if err != nil {
		log.Error().
			Str("error_code", err.Code()). // imprime o código do erro, que você define em cada erro
			Err(err).
			Msg("UserHandler.Update error")
		c.JSON(err.HTTPStatusCode(), gin.H{"error": err.Error(), "code": err.Code()})
		return
	}

	c.JSON(http.StatusOK, updatedPermission)
}

// Delete remove uma permissão.
func (h *PermissionHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	log.Info().Str("permission_id", idStr).Msg("Deleting permission")
	if err := h.service.Delete(idStr); err != nil {
		log.Error().
			Str("error_code", err.Code()). // imprime o código do erro, que você define em cada erro
			Err(err).
			Msg("UserHandler.Delete error")
		c.JSON(err.HTTPStatusCode(), gin.H{"error": err.Error(), "code": err.Code()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// FindAll retorna todas as permissões paginadas.
func (h *PermissionHandler) FindAll(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	log.Info().Msg("Fetching all permissions")
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", fmt.Sprintf("%d", h.cfg.AppDefaultAPIPageSize)))
	if err != nil || pageSize < 1 {
		pageSize = h.cfg.AppDefaultAPIPageSize
	}
	filters := c.Request.URL.Query()
	paginatedResult, err := h.service.FindAll(filters, page, pageSize)
	if customErr, ok := err.(service.ErrorUtil); ok {
		log.Error().
			Str("error_code", customErr.Code()). // imprime o código do erro, que você define em cada erro
			Err(customErr).
			Msg("UserHandler.FindAll error")
		c.JSON(customErr.HTTPStatusCode(), gin.H{"error": customErr.Error(), "code": customErr.Code()})
		return
	}

	c.JSON(http.StatusOK, paginatedResult)
}

// FindByID retorna uma permissão pelo ID.
func (h *PermissionHandler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	log.Info().Str("permission_id", idStr).Msg("Fetching permission by ID")
	permission, err := h.service.FindByID(idStr)
	if err != nil {
		log.Error().
			Str("error_code", err.Code()). // imprime o código do erro, que você define em cada erro
			Err(err).
			Msg("UserHandler.Update error")
		c.JSON(err.HTTPStatusCode(), gin.H{"error": err.Error(), "code": err.Code()})
		return
	}
	c.JSON(http.StatusOK, permission)
}
