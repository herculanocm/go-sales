package router

import (
	"go-sales/internal/database"
	"go-sales/internal/dto"
	"go-sales/internal/handler"
	"go-sales/internal/middleware"
	"go-sales/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupCompanyGlobalRoutes(router *gin.RouterGroup, db *gorm.DB) {
	repo := database.NewCompanyGlobalRepository(db)
	service := service.NewCompanyGlobalService(repo)
	handler := handler.NewCompanyGlobalHandler(service)

	router.POST("/company-globals", middleware.ValidateDTO(&dto.CreateCompanyGlobalDTO{}), handler.Create)
	router.PUT("/company-globals/:id", middleware.ValidateUUID("id"), middleware.ValidateDTO(&dto.CreateCompanyGlobalDTO{}), handler.Update)
	router.DELETE("/company-globals/:id", middleware.ValidateUUID("id"), handler.Delete)
}
