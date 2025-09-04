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

// SetupPermissionRoutes encapsula a configuração das rotas de permissões.
func SetupPermissionRoutes(router *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	// 1. Inicializar o Repositório
	permissionRepo := database.NewPermissionRepository(db)
	companyGlobalRepo := database.NewCompanyGlobalRepository(db)

	// 2. Inicializar o Serviço
	permissionService := service.NewPermissionService(permissionRepo, companyGlobalRepo)

	// 3. Inicializar o Handler
	permissionHandler := handler.NewPermissionHandler(permissionService, cfg)

	// 4. Definir as rotas para /permissions
	router.POST("/permissions", middleware.ValidateDTO(reflect.TypeOf(dto.CreatePermissionDTO{})), permissionHandler.Create)
	router.PUT("/permissions/:id", middleware.ValidateID("id"), middleware.ValidateDTO(reflect.TypeOf(dto.CreatePermissionDTO{})), permissionHandler.Update)
	router.DELETE("/permissions/:id", middleware.ValidateID("id"), permissionHandler.Delete)
	router.GET("/permissions/:id", middleware.ValidateID("id"), permissionHandler.FindByID)
	router.GET("/permissions", permissionHandler.FindAll)

	// adding restore point for soft delete
}
