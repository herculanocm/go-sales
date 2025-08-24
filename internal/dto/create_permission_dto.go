package dto

// CreatePermissionDTO representa os dados necessários para criar/atualizar uma permissão.
type CreatePermissionDTO struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}
