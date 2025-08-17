package main

import (
	"go-sales/internal/config"
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

	log.Info().Msg("Application started successfully.")
}
