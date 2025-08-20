package handler

import (
	"errors"
	"go-sales/internal/dto"
	"go-sales/internal/mapper"
	"go-sales/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CompanyGlobalHandler struct {
	service service.CompanyGlobalServiceInterface
}

func NewCompanyGlobalHandler(service service.CompanyGlobalServiceInterface) *CompanyGlobalHandler {
	return &CompanyGlobalHandler{
		service: service,
	}
}

// Create é o método do handler para criar uma nova empresa global.
func (h *CompanyGlobalHandler) Create(c *gin.Context) {
	validatedDTO, exists := c.Get("validatedDTO")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Validated DTO not found"})
		return
	}

	// Faz a asserção de tipo
	createCompanyDTO, ok := validatedDTO.(*dto.CreateCompanyGlobalDTO)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid DTO type"})
		return
	}

	// 2. Chamar a Camada de Serviço
	createdCompany, err := h.service.Create(*createCompanyDTO)
	if err != nil {
		// 3. Tratar Erros da Camada de Serviço
		// Verifica se o erro é um erro de negócio conhecido (CGC duplicado).
		if errors.Is(err, service.ErrCGCAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		// Para qualquer outro erro, consideramos um erro interno do servidor.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	// 4. Retornar Resposta de Sucesso
	c.JSON(http.StatusCreated, mapper.MapToCompanyGlobalDTO(createdCompany))
}

// Update é o método do handler para atualizar uma empresa global.
func (h *CompanyGlobalHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	validatedDTO, exists := c.Get("validatedDTO")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Validated DTO not found"})
		return
	}

	// Faz a asserção de tipo
	updateCompanyDTO, ok := validatedDTO.(*dto.CreateCompanyGlobalDTO)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid DTO type"})
		return
	}

	// 2. Chamar a Camada de Serviço
	updatedCompany, err := h.service.Update(*updateCompanyDTO, idStr)
	if err != nil {
		// 3. Tratar Erros da Camada de Serviço
		if errors.Is(err, service.EntityNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrCGCAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		// Para qualquer outro erro não esperado, retorne um erro genérico.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	// 4. Retornar Resposta de Sucesso
	c.JSON(http.StatusOK, mapper.MapToCompanyGlobalDTO(updatedCompany))
}

// Delete é o método do handler para apagar uma empresa global.
func (h *CompanyGlobalHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")

	// 2. Chamar a Camada de Serviço
	if err := h.service.Delete(idStr); err != nil {
		// 3. Tratar Erros da Camada de Serviço
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

func (h *CompanyGlobalHandler) FindByCGC(c *gin.Context) {
	cgcStr := c.Param("cgc")

	company, err := h.service.FindByCGC(cgcStr)
	if err != nil {
		// 3. Tratar Erros da Camada de Serviço
		if errors.Is(err, service.EntityNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		// Para qualquer outro erro, consideramos um erro interno do servidor.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	// 4. Retornar Resposta de Sucesso
	c.JSON(http.StatusOK, mapper.MapToCompanyGlobalDTO(company))
}

func (h *CompanyGlobalHandler) FindByID(c *gin.Context) {
	idStr := c.Param("id")

	company, err := h.service.FindByID(idStr)
	if err != nil {
		// 3. Tratar Erros da Camada de Serviço
		if errors.Is(err, service.EntityNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		// Para qualquer outro erro, consideramos um erro interno do servidor.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	// 4. Retornar Resposta de Sucesso
	c.JSON(http.StatusOK, mapper.MapToCompanyGlobalDTO(company))
}

func (h *CompanyGlobalHandler) FindAll(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// 2. Extrai os filtros.
	filters := c.Request.URL.Query()

	// 2. Chamar a Camada de Serviço, passando os filtros.
	// 3. Chama o serviço.
	paginatedResult, err := h.service.FindAll(filters, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		return
	}

	// 4. Retorna o resultado paginado.
	c.JSON(http.StatusOK, paginatedResult)
}
