package service

import (
	"errors"
	"net/http"
	"reflect"

	"gorm.io/gorm"
)

type ErrorUtil interface {
	HTTPStatusCode() int
	Code() string
	Error() string
}

type AbstractError struct {
	error          string
	httpStatusCode int
	code           string
}

// NewError cria um novo erro customizado de forma segura.
// Parâmetros:
// - msg: mensagem de erro.
// - statusCode: código HTTP (ex: http.StatusBadRequest).
// - code: código único do erro (ex: "validated_dto_not_found").
// Retorna: ErrorUtil (interface implementada por AbstractError).
func NewError(msg string, statusCode int, code string) ErrorUtil {
	return &AbstractError{
		error:          msg,
		httpStatusCode: statusCode,
		code:           code,
	}
}

func (e *AbstractError) Error() string { return e.error }

func (e *AbstractError) HTTPStatusCode() int {
	return e.httpStatusCode
}
func (e *AbstractError) Code() string {
	return e.code
}

var (
	// ErrEmailInUse é retornado quando uma tentativa de criar um usuário com um email que já existe.
	ErrEmailInUse = &AbstractError{
		error:          "email already in use",
		httpStatusCode: http.StatusConflict,
		code:           "email_in_use",
	}
	// ErrCGCInUse é retornado quando uma tentativa de criar uma empresa com um CGC que já existe.
	ErrCGCInUse = &AbstractError{
		error:          "cgc already in use",
		httpStatusCode: http.StatusConflict,
		code:           "cgc_in_use",
	}

	// ErrNotFound é retornado quando uma entidade não é encontrada.
	ErrNotFound = &AbstractError{
		error:          "entity not found",
		httpStatusCode: http.StatusNotFound,
		code:           "entity_not_found",
	}
	// ErrDuplicateKey é retornado quando uma violação de chave única ocorre.
	ErrDuplicateKey = &AbstractError{
		error:          "a record with this key already exists",
		httpStatusCode: http.StatusConflict,
		code:           "duplicate_key",
	}
	// ErrForeignKeyConstraint é retornado ao tentar apagar um registro que é referenciado por outro.
	ErrForeignKeyConstraint = &AbstractError{
		error:          "cannot delete this record because it is referenced by other records",
		httpStatusCode: http.StatusConflict,
		code:           "foreign_key_constraint",
	}

	ErrForeignKeyViolated = &AbstractError{
		error:          "foreign key violated",
		httpStatusCode: http.StatusConflict,
		code:           "foreign_key_violated",
	}

	ErrInvalidData = &AbstractError{
		error:          "invalid data",
		httpStatusCode: http.StatusBadRequest,
		code:           "invalid_data",
	}

	ErrCompanyGlobalNotFound = &AbstractError{
		error:          "company global not found",
		httpStatusCode: http.StatusNotFound,
		code:           "company_global_not_found",
	}
	ErrRoleNotFound = &AbstractError{
		error:          "role not found",
		httpStatusCode: http.StatusNotFound,
		code:           "role_not_found",
	}
	ErrPermissionNameInUse = &AbstractError{
		error:          "permission name already in use",
		httpStatusCode: http.StatusConflict,
		code:           "permission_name_in_use",
	}
	ErrRoleNameInUse = &AbstractError{
		error:          "role name already in use",
		httpStatusCode: http.StatusConflict,
		code:           "role_name_in_use",
	}
	ErrPermissionNotFound = &AbstractError{
		error:          "permission not found",
		httpStatusCode: http.StatusNotFound,
		code:           "permission_not_found",
	}
	ErrInternalServer = &AbstractError{
		error:          "internal server error",
		httpStatusCode: http.StatusInternalServerError,
		code:           "internal_server_error",
	}
	ErrDatabase = &AbstractError{
		error:          "database error",
		httpStatusCode: http.StatusInternalServerError,
		code:           "database_error",
	}

	ErrCompanyGlobalContactsFieldValidationUnique = &AbstractError{
		error:          "the properties cgc, phone and email must be unique in contacts",
		httpStatusCode: http.StatusConflict,
		code:           "company_global_contacts_field_validation_unique",
	}

	ErrRoleMustHavePermissions = &AbstractError{
		error:          "role must have at least one permission",
		httpStatusCode: http.StatusBadRequest,
		code:           "role_must_have_permissions",
	}

	ErrPermissionsNotFound = &AbstractError{
		error:          "one or more permissions not found or do not belong to the specified company global",
		httpStatusCode: http.StatusBadRequest,
		code:           "permissions_not_found",
	}
)

func GormDefaultError(err error) ErrorUtil {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return ErrForeignKeyViolated
	}
	if errors.Is(err, gorm.ErrInvalidData) {
		return ErrInvalidData
	}
	if errors.Is(err, gorm.ErrInvalidField) {
		return ErrInvalidData
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrDuplicateKey
	}

	// Tenta acessar o campo Code via reflexão
	if err != nil {
		val := reflect.ValueOf(err)
		// Se for ponteiro, pega o elemento
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		codeField := val.FieldByName("Code")
		if codeField.IsValid() && codeField.Kind() == reflect.String {
			code := codeField.String()
			switch code {
			case "23505": // unique_violation
				return ErrDuplicateKey
			case "23503": // foreign_key_violation
				return ErrForeignKeyViolated
			}
		}
	}
	return ErrDatabase
}
