package service

import (
	"errors"
	"go-sales/internal/database"
	"strconv"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func CheckCompanyGlobalExists(repo database.CompanyGlobalRepositoryInterface, companyGlobalID int64, useUnscoped bool) (bool, ErrorUtil) {
	companyGlobalExists, err := repo.Exists(companyGlobalID, useUnscoped)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().
			Err(err).
			Str("company_global_id", strconv.FormatInt(companyGlobalID, 10)).
			Msg("failed to check if company global exists")
		return false, GormDefaultError(err)
	}
	if !companyGlobalExists {
		log.Error().
			Err(err).
			Caller().
			Str("company_global_id", strconv.FormatInt(companyGlobalID, 10)).
			Msg("failed to find existing company global")
		return false, ErrCompanyGlobalNotFound
	}
	return true, nil
}

func CheckUserEmailExists(repo database.UserRepositoryInterface, email string, companyGlobalID int64, useUnscoped bool) (bool, ErrorUtil) {
	existingUser, err := repo.EmailExists(email, companyGlobalID, useUnscoped)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().
			Err(err).
			Caller().
			Str("email", email).
			Str("company_global_id", strconv.FormatInt(companyGlobalID, 10)).
			Msg("failed to check if user email exists")
		return false, GormDefaultError(err)
	}
	return existingUser, nil
}
