package handler

import (
	"errors"
	"fmt"
	"go-sales/internal/config"
	"go-sales/internal/dto"
	"go-sales/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	validatedDTO, exists := c.Get("validatedDTO")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Validated DTO not found"})
		return
	}

	createPermissionDTO, ok := validatedDTO.(*dto.CreatePermissionDTO)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid DTO type"})
		return
	}

	createdPermission, err := h.service.Create(*createPermissionDTO)
	if err != nil {
		if errors.Is(err, service.ErrPermissionNameInUse) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	c.JSON(http.StatusCreated, createdPermission)
}

// Update atualiza uma permissão existente.
func (h *PermissionHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	validatedDTO, exists := c.Get("validatedDTO")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Validated DTO not found"})
		return
	}

	updatePermissionDTO, ok := validatedDTO.(*dto.CreatePermissionDTO)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid DTO type"})
		return
	}

	updatedPermission, err := h.service.Update(*updatePermissionDTO, idStr)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrPermissionNameInUse) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	c.JSON(http.StatusOK, updatedPermission)
}

// Delete remove uma permissão.
func (h *PermissionHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	if err := h.service.Delete(idStr); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// FindAll retorna todas as permissões paginadas.
func (h *PermissionHandler) FindAll(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", fmt.Sprintf("%d", h.cfg.AppDefaultAPIPageSize)))
	if err != nil || pageSize < 1 {
		pageSize = h.cfg.AppDefaultAPIPageSize
	}
	filters := c.Request.URL.Query()
	paginatedResult, err := h.service.FindAll(filters, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusOK, paginatedResult)
}

// FindByID retorna uma permissão pelo ID.
func (h *PermissionHandler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	permission, err := h.service.FindByID(idStr)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusOK, permission)
}
