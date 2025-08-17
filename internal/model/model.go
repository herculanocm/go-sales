package model

// RegisteredModels é uma slice que contém todas as structs de modelo
// que devem ser incluídas na migração automática do GORM.
var RegisteredModels []interface{}

// A função init é executada automaticamente pelo Go quando o pacote 'model' é importado.
// Nós a usamos para registrar todos os nossos modelos em um único lugar.
func init() {
	RegisteredModels = append(
		RegisteredModels,
		&User{},
		// Adicione novos modelos aqui no futuro.
		// Por exemplo: &Product{}, &Order{},
	)
}

// Tabler é uma interface que os modelos podem implementar
// para fornecer um nome de tabela personalizado.
// É bom manter esta definição aqui também.
type Tabler interface {
	TableName() string
}
