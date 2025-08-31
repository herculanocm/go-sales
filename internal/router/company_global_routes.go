package router

import (
	"go-sales/internal/config"
	"go-sales/internal/database"
	"go-sales/internal/dto"
	"go-sales/internal/handler"
	"go-sales/internal/middleware"
	"go-sales/internal/service"
	"reflect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupCompanyGlobalRoutes(router *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	repo := database.NewCompanyGlobalRepository(db)
	service := service.NewCompanyGlobalService(repo)
	handler := handler.NewCompanyGlobalHandler(service, cfg)

	router.POST("/company-globals", middleware.ValidateDTO(reflect.TypeOf(dto.CreateCompanyGlobalDTO{})), handler.Create)
	router.PUT("/company-globals/:id", middleware.ValidateUUID("id"), middleware.ValidateDTO(reflect.TypeOf(dto.CreateCompanyGlobalDTO{})), handler.Update)
	router.DELETE("/company-globals/:id", middleware.ValidateUUID("id"), handler.Delete)
	router.GET("/company-globals/:id", middleware.ValidateUUID("id"), handler.FindByID)
	router.GET("/company-globals", handler.FindAll)

	// adding restore point for soft delete
	router.POST("/company-globals/:id/restore", middleware.ValidateUUID("id"), handler.Restore)
}
