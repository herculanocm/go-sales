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
	"github.com/rs/zerolog/log"
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
	log.Info().Msg("Creating a new role")

	createDTO, utilError := GetValidatedDTO[*dto.CreateRoleDTO](c, "validatedDTO")
	if utilError != nil {
		HandleError(utilError, "RoleHandler.Create - Error getting validated DTO", c)
		return
	}

	createdRole, err := h.service.Create(*createDTO)
	if err != nil {
		HandleError(err, "RoleHandler.Create error", c)
		return
	}
	c.JSON(http.StatusCreated, createdRole)
}

func (h *RoleHandler) Update(c *gin.Context) {
	log.Info().Msg("Updating role")
	idStr := c.Param("id")

	updateRoleDTO, utilError := GetValidatedDTO[*dto.RoleDTO](c, "validatedDTO")
	if utilError != nil {
		HandleError(utilError, "RoleHandler.Update - Error getting validated DTO", c)
		return
	}

	id, err := util.ParseSnowflake(idStr)
	if err != nil {
		customError := service.NewError("invalid_id_format", http.StatusBadRequest, "invalid_id_format")
		HandleError(customError, "RoleHandler.Update - Error parsing ID", c)
		return
	}

	updatedRole, errUpdate := h.service.Update(*updateRoleDTO, id)
	if err != nil {
		HandleError(errUpdate, "RoleHandler.Update error", c)
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
	companyGlobalId, err := strconv.ParseInt(c.Query("company_global_id"), 10, 64)
	if err != nil {
		customError := service.NewError("invalid company_global_id format", http.StatusBadRequest, "invalid_company_global_id_format")
		HandleError(customError, "RoleHandler.FindAll - Error parsing company_global_id", c)
		return
	}
	filters := c.Request.URL.Query()
	paginatedResult, err := h.service.FindAll(filters, page, pageSize, companyGlobalId)
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
