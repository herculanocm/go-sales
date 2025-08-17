package main

import (
	"go-sales/internal/config"
	"go-sales/internal/database"
	"go-sales/internal/logger"
	dlog "log"

	"github.com/rs/zerolog/log"
)

func main() {
	dlog.Println("Starting application...")
	// 1. Carregue a configuração PRIMEIRO.
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

	log.Info().Msg("Trying to migrate database...")
	err = database.Migrate(database.DB)
	if err != nil {
		log.Fatal().Msgf("Error migrating database: %v", err)
	}
	log.Info().Msg("Database migrated successfully.")

	log.Info().Msg("Application started successfully.")
}
