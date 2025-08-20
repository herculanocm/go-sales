package dto

// PageInfo contém os metadados da paginação.
type PageInfo struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalItems int64 `json:"totalItems"`
	TotalPages int   `json:"totalPages"`
}

// PaginatedResponse é uma estrutura de resposta genérica para qualquer lista paginada.
// [T any] significa que T é um parâmetro de tipo que pode ser qualquer tipo (any).
type PaginatedResponse[T any] struct {
	Items    []T      `json:"items"`
	PageInfo PageInfo `json:"pageInfo"`
}
