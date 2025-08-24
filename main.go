package main

import (
	"go-sales/internal/config"
	"go-sales/internal/database"
	"go-sales/internal/logger"
	"go-sales/internal/router"
	"go-sales/internal/validator" // Importe o novo pacote validator
	dlog "log"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	dlog.Println("Starting application...")

	cfg, err := config.LoadConfig(".")
	if err != nil {
		// Se a configuração falhar, não temos nosso logger zerolog ainda.
		// Usamos o log padrão do Go como fallback para este erro crítico.
		dlog.Fatalf("Fatal Error: %v", err)
	}

	logger.Init(cfg)
	log.Info().Msg("Configuration loaded successfully with configs")
	log.Info().Msgf("App environment: %s", cfg.AppEnv)
	log.Info().Msgf("Configurations: %s", cfg.String())

	log.Info().Msgf("Trying to connect to database, ensure the schema %s exists...", cfg.DBSchema)
	err = database.Connect(cfg)
	if err != nil {
		log.Fatal().Msgf("Error connecting to database: %v", err)
	}
	log.Info().Msg("Database connection established successfully")

	log.Info().Msg("Setting up Gin server...")
	// Inicializa o validador customizado ANTES de criar a instância do Gin.
	validator.InitCustomValidator()

	// Define o modo do Gin com base na configuração (ex: "release" para produção)
	gin.SetMode(gin.ReleaseMode)
	if cfg.AppEnv != "production" {
		gin.SetMode(gin.DebugMode)
	}

	//server := gin.Default() // Cria um router com middlewares padrão (logger, recovery)

	// Use gin.New() para um controle explícito dos middlewares.
	server := gin.New()

	// Adicione seus middlewares customizados na ordem que devem executar.
	// 1. Middleware de log estruturado (usando zerolog).
	server.Use(logger.GinLogger()) // Você precisará criar esta função.

	// 2. Middleware de recuperação de pânico customizado.
	server.Use(logger.GinRecovery(true)) // Você precisará criar esta função.

	// Agrupar rotas sob um prefixo
	log.Info().Msgf("Setting up API routes with prefix: %s", cfg.AppAPIPrefix)
	api := server.Group(cfg.AppAPIPrefix)

	// --- REGISTRO DAS ROTAS MODULARES ---
	// Passe o grupo de rotas e a conexão com o DB para a função de setup.
	router.SetupUserRoutes(api, database.DB, cfg)
	router.SetupCompanyGlobalRoutes(api, database.DB, cfg)

	log.Info().Msgf("Server is starting on port %s...", cfg.AppAPIPort)
	// Inicia o servidor
	if err := server.Run(":" + cfg.AppAPIPort); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
