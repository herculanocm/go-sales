package service

import "errors"

var (
	// ErrEmailInUse é retornado quando uma tentativa de criar um usuário com um email que já existe.
	ErrEmailInUse = errors.New("email already in use")
	// ErrCGCInUse é retornado quando uma tentativa de criar uma empresa com um CGC que já existe.
	ErrCGCInUse = errors.New("cgc already in use")
	// ErrNotFound é retornado quando uma entidade não é encontrada.
	ErrNotFound = errors.New("entity not found")
	// ErrDuplicateKey é retornado quando uma violação de chave única ocorre.
	ErrDuplicateKey = errors.New("a record with this key already exists")
	// ErrForeignKeyConstraint é retornado ao tentar apagar um registro que é referenciado por outro.
	ErrForeignKeyConstraint = errors.New("cannot delete this record because it is referenced by other records")

	ErrCompanyGlobalNotFound = errors.New("company global not found")
	ErrRoleNotFound          = errors.New("role not found")
)
