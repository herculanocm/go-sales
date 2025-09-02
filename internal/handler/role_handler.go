package handler

import (
	"errors"
	"fmt"
	"go-sales/internal/config"
	"go-sales/internal/dto"
	"go-sales/internal/service"
	"go-sales/pkg/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RoleHandler encapsula a dependência do serviço de role.
type RoleHandler struct {
	service service.RoleServiceInterface
	cfg     *config.Config
}

// NewRoleHandler cria uma nova instância do handler de role.
func NewRoleHandler(s service.RoleServiceInterface, cfg *config.Config) *RoleHandler {
	return &RoleHandler{
		service: s,
		cfg:     cfg,
	}
}

func (h *RoleHandler) Create(c *gin.Context) {
	validatedDTO, exists := c.Get("validatedDTO")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Validated DTO not found"})
		return
	}
	createRoleDTO, ok := validatedDTO.(*dto.CreateRoleDTO)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid DTO type"})
		return
	}
	createdRole, err := h.service.Create(*createRoleDTO)
	if err != nil {
		if errors.Is(err, service.ErrRoleNameInUse) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusCreated, createdRole)
}

func (h *RoleHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	validatedDTO, exists := c.Get("validatedDTO")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Validated DTO not found"})
		return
	}
	updateRoleDTO, ok := validatedDTO.(*dto.CreateRoleDTO)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid DTO type"})
		return
	}

	id, err := util.ParseSnowflake(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	updatedRole, err := h.service.Update(*updateRoleDTO, id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrRoleNameInUse) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusOK, updatedRole)
}

func (h *RoleHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := util.ParseSnowflake(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *RoleHandler) FindAll(c *gin.Context) {
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

func (h *RoleHandler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := util.ParseSnowflake(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	role, err := h.service.FindByID(id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}
	c.JSON(http.StatusOK, role)
}
