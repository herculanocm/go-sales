package router

import (
	"go-sales/internal/config"
	"go-sales/internal/database"
	"go-sales/internal/dto"
	"go-sales/internal/handler"
	"go-sales/internal/middleware"
	"go-sales/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoleRoutes encapsula a configuração das rotas de roles.
func SetupRoleRoutes(router *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	roleRepo := database.NewRoleRepository(db)
	permRepo := database.NewPermissionRepository(db)
	companyRepo := database.NewCompanyGlobalRepository(db)
	roleService := service.NewRoleService(roleRepo, permRepo, companyRepo)
	roleHandler := handler.NewRoleHandler(roleService, cfg)

	router.POST("/roles", middleware.ValidateDTO(&dto.CreateRoleDTO{}), roleHandler.Create)
	router.PUT("/roles/:id", middleware.ValidateUUID("id"), middleware.ValidateDTO(&dto.CreateRoleDTO{}), roleHandler.Update)
	router.DELETE("/roles/:id", middleware.ValidateUUID("id"), roleHandler.Delete)
	router.GET("/roles/:id", middleware.ValidateUUID("id"), roleHandler.FindByID)
	router.GET("/roles", roleHandler.FindAll)
}
