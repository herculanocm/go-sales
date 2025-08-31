package handler

import (
	"net/http"

	"go-sales/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func HandleError(err service.ErrorUtil, msg string, c *gin.Context) {
	log.Error().
		Err(err).
		Caller().
		Msg(msg)
	c.JSON(err.HTTPStatusCode(), gin.H{"error": err.Error(), "code": err.Code()})
}

func GetValidatedDTO[T any](c *gin.Context, key string) (T, service.ErrorUtil) {
	value, exists := c.Get(key)
	if !exists {
		err := service.NewError("Validated DTO not found", http.StatusBadRequest, "validated_dto_not_found")
		return *new(T), err
	}

	// 2. Fazer a asserção de tipo
	dto, ok := value.(T)
	if !ok {
		err := service.NewError("Type assertion failed for validatedDTO", http.StatusBadRequest, "type_assertion_failed")
		return *new(T), err
	}

	return dto, nil
}
