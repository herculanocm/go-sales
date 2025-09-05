package validator

import (
	"go-sales/pkg/util"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Valida se o campo é um Snowflake válido (int64 > 0)
func SnowflakeValidator(fl validator.FieldLevel) bool {
	switch fl.Field().Kind() {
	case reflect.Int, reflect.Int64:
		return fl.Field().Int() > 0
	case reflect.String:
		_, err := util.ParseSnowflake(fl.Field().String())
		return err == nil
	default:
		return false
	}
}

func InitSnowflakeValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("snowflake", SnowflakeValidator)
	}
}
