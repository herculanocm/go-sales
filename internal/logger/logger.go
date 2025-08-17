package logger

import (
	"go-sales/internal/config"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init configura o logger global zerolog.
// Ele usa um formato legível para desenvolvimento e JSON para produção.
func Init(cfg *config.Config) {
	var logLevel zerolog.Level = zerolog.InfoLevel // Nível padrão é Info

	if cfg.AppEnv == "production" {
		// Em produção, usamos o formato JSON padrão.
		// O timestamp é no formato Unix.
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		logLevel = zerolog.WarnLevel
	} else {
		// Em desenvolvimento, usamos um formato mais legível no console.
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	// Define o nível de log global. Apenas logs com este nível ou superior serão exibidos.
	// Ex: zerolog.InfoLevel irá mostrar INFO, WARN, ERROR, FATAL, PANIC.
	zerolog.SetGlobalLevel(logLevel)

	log.Info().Msg("Logger initialized.")
}

// GinLogger é um middleware para logar requisições HTTP usando zerolog.
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Processa a requisição
		c.Next()

		// Após a requisição, loga os detalhes
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		event := log.Info()
		if statusCode >= http.StatusInternalServerError {
			event = log.Error().Str("error", errorMessage)
		}

		event.
			Str("method", method).
			Str("path", path).
			Int("status_code", statusCode).
			Int("body_size", bodySize).
			Str("client_ip", clientIP).
			Dur("latency", latency).
			Msg("Request processed")
	}
}

// GinRecovery é um middleware que recupera de qualquer pânico e loga usando zerolog.
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Loga o pânico com stack trace
				log.Error().
					Interface("error", err).
					Bool("stack", stack).
					Msg("Panic recovered")

				// Retorna uma resposta de erro genérica para o cliente
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "An internal server error occurred",
				})
			}
		}()
		c.Next()
	}
}
