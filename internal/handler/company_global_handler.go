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

type CompanyGlobalHandler struct {
	service service.CompanyGlobalServiceInterface
	cfg     *config.Config
}

func NewCompanyGlobalHandler(service service.CompanyGlobalServiceInterface, cfg *config.Config) *CompanyGlobalHandler {
	return &CompanyGlobalHandler{
		service: service,
		cfg:     cfg,
	}
}

// Create é o método do handler para criar uma nova empresa global.
func (h *CompanyGlobalHandler) Create(c *gin.Context) {
	log.Info().Msg("Creating a new company global")

	createCompanyDTO, utilError := GetValidatedDTO[*dto.CreateCompanyGlobalDTO](c, "validatedDTO")
	if utilError != nil {
		HandleError(utilError, "CompanyGlobalHandler.Create - Error getting validated DTO", c)
		return
	}

	createdCompany, err := h.service.Create(*createCompanyDTO)
	if err != nil {
		HandleError(err, "CompanyGlobalHandler.Create error", c)
		return
	}

	// 4. Retornar Resposta de Sucesso
	c.JSON(http.StatusCreated, createdCompany)
}

// Update é o método do handler para atualizar uma empresa global.
func (h *CompanyGlobalHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	log.Info().Str("id", idStr).Msg("Updating company global")
	updateCompanyDTO, utilError := GetValidatedDTO[*dto.CreateCompanyGlobalDTO](c, "validatedDTO")
	if utilError != nil {
		HandleError(utilError, "CompanyGlobalHandler.Update - Error getting validated DTO", c)
		return
	}

	id, err := util.Parse(idStr)
	if err != nil {
		customError := service.NewError(err.Error(), http.StatusBadRequest, "invalid_id_format")
		HandleError(customError, "CompanyGlobalHandler.Update - Error parsing ID", c)
		return
	}

	// 2. Chamar a Camada de Serviço
	updatedCompany, errUpdate := h.service.Update(*updateCompanyDTO, id)
	if errUpdate != nil {
		HandleError(errUpdate, "CompanyGlobalHandler.Update error", c)
		return
	}

	// 4. Retornar Resposta de Sucesso
	c.JSON(http.StatusOK, updatedCompany)
}

// Delete é o método do handler para apagar uma empresa global.
func (h *CompanyGlobalHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	log.Info().Str("id", idStr).Msg("Deleting company global")
	id, err := util.Parse(idStr)
	if err != nil {
		customError := service.NewError(err.Error(), http.StatusBadRequest, "invalid_id_format")
		HandleError(customError, "CompanyGlobalHandler.Delete - Error parsing ID", c)
		return
	}

	// 2. Chamar a Camada de Serviço
	if err := h.service.Delete(id); err != nil {
		HandleError(err, "CompanyGlobalHandler.Delete error", c)
		return
	}

	// 4. Retornar Resposta de Sucesso
	c.JSON(http.StatusNoContent, nil)
}

func (h *CompanyGlobalHandler) FindByCGC(c *gin.Context) {
	cgcStr := c.Param("cgc")
	log.Info().Str("cgc", cgcStr).Msg("Finding company global by CGC")
	company, err := h.service.FindByCGC(cgcStr)

	if err != nil {
		HandleError(err, "CompanyGlobalHandler.FindByCGC error", c)
		return
	}

	// 4. Retornar Resposta de Sucesso
	c.JSON(http.StatusOK, company)
}

func (h *CompanyGlobalHandler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	log.Info().Str("id", idStr).Msg("Finding company global by ID")
	id, err := util.Parse(idStr)
	if err != nil {
		customError := service.NewError(err.Error(), http.StatusBadRequest, "invalid_id_format")
		HandleError(customError, "CompanyGlobalHandler.FindByID - Error parsing ID", c)
		return
	}

	company, errFind := h.service.FindByID(id)
	if errFind != nil {
		HandleError(errFind, "CompanyGlobalHandler.FindByID error", c)
		return
	}

	// 4. Retornar Resposta de Sucesso
	c.JSON(http.StatusOK, company)
}

func (h *CompanyGlobalHandler) FindAll(c *gin.Context) {
	log.Info().Msg("Fetching all company globals")
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
	paginatedResult, customErr := h.service.FindAll(filters, page, pageSize)
	if customErr != nil {

		HandleError(customErr, "CompanyGlobalHandler.FindAll error", c)
		return

	}

	// 4. Retorna o resultado paginado.
	c.JSON(http.StatusOK, paginatedResult)
}
