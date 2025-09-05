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

	createPermissionDTO, utilError := GetValidatedDTO[*dto.CreatePermissionDTO](c, "validatedDTO")
	if utilError != nil {
		HandleError(utilError, "PermissionHandler.Create - Error getting validated DTO", c)
		return
	}

	createdPermission, err := h.service.Create(*createPermissionDTO)
	if err != nil {
		HandleError(err, "PermissionHandler.Create error", c)
		return
	}

	c.JSON(http.StatusCreated, createdPermission)
}

// Update atualiza uma permissão existente.
func (h *PermissionHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	log.Info().Str("permission_id", idStr).Msg("Updating permission")

	updatePermissionDTO, utilError := GetValidatedDTO[*dto.CreatePermissionDTO](c, "validatedDTO")
	if utilError != nil {
		HandleError(utilError, "PermissionHandler.Update - Error getting validated DTO", c)
		return
	}

	id, err := util.ParseSnowflake(idStr)
	if err != nil {
		customError := service.NewError(err.Error(), http.StatusBadRequest, "invalid_id_format")
		HandleError(customError, "PermissionHandler.Update - Error parsing ID", c)
		return
	}

	updatedPermission, errUpdate := h.service.Update(*updatePermissionDTO, id)
	if errUpdate != nil {
		HandleError(errUpdate, "PermissionHandler.Update error", c)
		return
	}

	c.JSON(http.StatusOK, updatedPermission)
}

// Delete remove uma permissão.
func (h *PermissionHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	log.Info().Str("permission_id", idStr).Msg("Deleting permission")

	id, err := util.ParseSnowflake(idStr)
	if err != nil {
		customError := service.NewError(err.Error(), http.StatusBadRequest, "invalid_id_format")
		HandleError(customError, "PermissionHandler.Delete - Error parsing ID", c)
		return
	}

	if err := h.service.Delete(id); err != nil {
		HandleError(err, "PermissionHandler.Delete error", c)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// FindAll retorna todas as permissões paginadas.
func (h *PermissionHandler) FindAll(c *gin.Context) {
	log.Info().Msg("Fetching all permissions")
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
		HandleError(customError, "PermissionHandler.FindAll - Error parsing company_global_id", c)
		return
	}
	filters := c.Request.URL.Query()
	paginatedResult, customErr := h.service.FindAll(filters, page, pageSize, companyGlobalId)
	if customErr != nil {

		HandleError(customErr, "PermissionHandler.FindAll error", c)
		return

	}

	c.JSON(http.StatusOK, paginatedResult)
}

// FindByID retorna uma permissão pelo ID.
func (h *PermissionHandler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	log.Info().Str("permission_id", idStr).Msg("Fetching permission by ID")

	id, err := util.ParseSnowflake(idStr)
	if err != nil {
		customError := service.NewError(err.Error(), http.StatusBadRequest, "invalid_id_format")
		HandleError(customError, "PermissionHandler.FindByID - Error parsing ID", c)
		return
	}

	permission, errFind := h.service.FindByID(id)
	if errFind != nil {
		HandleError(errFind, "PermissionHandler.FindByID error", c)
		return
	}
	c.JSON(http.StatusOK, permission)
}
