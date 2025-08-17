package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ValidateDTO valida o corpo da requisição com base na struct DTO fornecida.
func ValidateDTO(dto interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// É importante passar um ponteiro para a struct para que ShouldBindJSON possa preenchê-la.
		// No entanto, a validação funciona com o tipo, então não precisamos de um novo ponteiro a cada chamada.
		if err := c.ShouldBindJSON(dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		// Coloca o DTO validado no contexto para que o handler possa usá-lo.
		c.Set("validatedDTO", dto)
		c.Next()
	}
}
