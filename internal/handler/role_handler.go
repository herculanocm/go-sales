package handler

import (
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
	if errUpdate != nil {
		HandleError(errUpdate, "RoleHandler.Update error", c)
		return
	}
	c.JSON(http.StatusOK, updatedRole)
}

func (h *RoleHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	log.Info().Str("role_id", idStr).Msg("Deleting role")

	id, err := util.ParseSnowflake(idStr)
	if err != nil {
		customError := service.NewError(err.Error(), http.StatusBadRequest, "invalid_id_format")
		HandleError(customError, "RoleHandler.Delete - Error parsing ID", c)
		return
	}
	if err := h.service.Delete(id); err != nil {
		HandleError(err, "RoleHandler.Delete error", c)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *RoleHandler) FindAll(c *gin.Context) {
	log.Info().Msg("Fetching all roles")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", fmt.Sprintf("%d", h.cfg.AppDefaultAPIPageSize)))
	if err != nil || pageSize < 1 {
		pageSize = h.cfg.AppDefaultAPIPageSize
	}

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

	filters := c.Request.URL.Query()
	paginatedResult, errFindAll := h.service.FindAll(filters, page, pageSize, companyGlobalId)
	if errFindAll != nil {
		HandleError(errFindAll, "RoleHandler.FindAll error", c)
		return
	}
	c.JSON(http.StatusOK, paginatedResult)
}

func (h *RoleHandler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	log.Info().Str("role_id", idStr).Msg("Fetching role by ID")
	id, err := util.ParseSnowflake(idStr)
	if err != nil {
		customError := service.NewError(err.Error(), http.StatusBadRequest, "invalid_id_format")
		HandleError(customError, "RoleHandler.FindByID - Error parsing ID", c)
		return
	}
	role, errFind := h.service.FindByID(id)
	if errFind != nil {
		HandleError(errFind, "RoleHandler.FindByID error", c)
		return
	}
	c.JSON(http.StatusOK, role)
}
