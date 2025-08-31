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

// SetupUserRoutes encapsula a configuração das rotas de usuário.
func SetupUserRoutes(router *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	// 1. Inicializar o Repositório
	userRepo := database.NewUserRepository(db)
	userCompanyRepo := database.NewCompanyGlobalRepository(db)
	userRoleRepo := database.NewRoleRepository(db)

	// 2. Inicializar o Serviço
	userService := service.NewUserService(userRepo, userCompanyRepo, userRoleRepo)

	// 3. Inicializar o Handler
	userHandler := handler.NewUserHandler(userService, cfg)

	// 4. Definir as rotas para /users
	router.POST("/users", middleware.ValidateDTO(reflect.TypeOf(dto.CreateUserDTO{})), userHandler.Create)
	router.PUT("/users/:id", middleware.ValidateUUID("id"), middleware.ValidateDTO(reflect.TypeOf(dto.CreateUserDTO{})), userHandler.Update)
	router.DELETE("/users/:id", middleware.ValidateUUID("id"), userHandler.Delete)
	router.GET("/users/:id", middleware.ValidateUUID("id"), userHandler.FindByID)
	router.GET("/users", userHandler.FindAll)
}
