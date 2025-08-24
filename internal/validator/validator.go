package validator

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// InitCustomValidator inicializa o validador do Gin para usar as tags 'json'
// nos nomes dos campos das mensagens de erro, em vez dos nomes dos campos da struct Go.
func InitCustomValidator() {
	// Acessa o motor de validação subjacente do Gin.
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Registra uma função customizada para obter o nome do campo.
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			// Pega o valor da tag 'json'.
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			// Se a tag for "-", o campo não deve ser incluído.
			if name == "-" {
				return ""
			}
			// Retorna o nome do campo JSON.
			return name
		})
	}
}
