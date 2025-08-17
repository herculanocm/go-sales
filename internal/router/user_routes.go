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

// SetupUserRoutes encapsula a configuração das rotas de usuário.
func SetupUserRoutes(router *gin.RouterGroup, db *gorm.DB) {
	// 1. Inicializar o Repositório
	userRepo := database.NewUserRepository(db)

	// 2. Inicializar o Serviço
	userService := service.NewUserService(userRepo)

	// 3. Inicializar o Handler
	userHandler := handler.NewUserHandler(userService)

	// 4. Definir as rotas para /users
	router.POST("/users", middleware.ValidateDTO(&dto.CreateUserDTO{}), userHandler.Create)
	// router.GET("/users/:id", userHandler.GetByID) // Exemplo para futuras rotas
	// router.GET("/users", userHandler.GetAll)      // Exemplo para futuras rotas
}
