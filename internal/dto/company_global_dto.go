package dto

type CreateCompanyGlobalDTO struct {
	Name        string `json:"name" binding:"required,max=255"`
	Description string `json:"description" binding:"max=4000"`
	CGC         string `json:"cgc" binding:"required,max=14"`
	Enabled     bool   `json:"enabled" binding:"required"`
}
