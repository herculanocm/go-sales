package logger

import (
	"go-sales/internal/config"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init configura o logger global zerolog.
// Ele usa um formato legível para desenvolvimento e JSON para produção.
func Init(cfg *config.Config) {
	if cfg.AppEnv == "production" {
		// Em produção, usamos o formato JSON padrão.
		// O timestamp é no formato Unix.
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	} else {
		// Em desenvolvimento, usamos um formato mais legível no console.
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	// Define o nível de log global. Apenas logs com este nível ou superior serão exibidos.
	// Ex: zerolog.InfoLevel irá mostrar INFO, WARN, ERROR, FATAL, PANIC.
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Info().Msg("Logger initialized.")
}
