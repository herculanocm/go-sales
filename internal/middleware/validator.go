package middleware

import (
	"errors"
	"go-sales/pkg/util"
	"net/http"

	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateDTO(dtoType reflect.Type) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Cria uma nova instância do DTO para cada requisição
		dto := reflect.New(dtoType).Interface()

		if err := c.ShouldBindJSON(dto); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				out := make([]map[string]string, len(ve))
				typ := dtoType
				if typ.Kind() == reflect.Ptr {
					typ = typ.Elem()
				}
				for i, fe := range ve {
					field, ok := typ.FieldByName(fe.StructField())
					jsonTag := ""
					if ok {
						jsonTag = field.Tag.Get("json")
					}
					jsonName := strings.Split(jsonTag, ",")[0]
					if jsonName == "" {
						jsonName = fe.Field()
					}
					out[i] = map[string]string{
						"field":   jsonName,
						"message": fe.Error(),
					}
				}
				c.JSON(http.StatusBadRequest, gin.H{"errors": out, "code": "validation_error"})
				c.Abort()
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "validation_error"})
			c.Abort()
			return
		}
		c.Set("validatedDTO", dto)
		c.Next()
	}
}

// ValidateUUID verifica se um parâmetro da URL é um UUID válido.
func ValidateUUID(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param(paramName)
		if _, err := util.Parse(idStr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid " + paramName + " format. Must be a valid UUID.", "code": "invalid_id_format"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ValidateSnowflake verifica se um parâmetro da URL é um Snowflake válido.
func ValidateSnowflake(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param(paramName)
		if _, err := util.ParseSnowflake(idStr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid " + paramName + " format. Must be a valid Snowflake ID.",
				"code":  "invalid_id_format",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func ValidateID(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param(paramName)
		if _, err := util.ParseSnowflake(idStr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid " + paramName + " format. Must be a valid Snowflake ID.",
				"code":  "invalid_id_format",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func ValidateCGC(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cgcStr := c.Param(paramName)
		if !util.IsValidCGC(cgcStr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid " + paramName + " format. Must be a valid CGC.", "code": "invalid_cgc_format"})
			c.Abort()
			return
		}
		c.Next()
	}
}
