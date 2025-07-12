package error

import (
	"errors"
	domainError "todo-api-go/internal/domain/error"
)

func IsDuplicateEmailError(err error) bool {
	var userErr *domainError.UserDomainError
	if errors.As(err, &userErr) {
		return userErr.Code() == "DUPLICATE_EMAIL"
	}
	return false
}

func IsValidationError(err error) bool {

	var userErr *domainError.UserDomainError
	if errors.As(err, &userErr) {
		return userErr.Code() == "INVALID_USER_DATA"
	}
	return false
}

func IsAuthenticationError(err error) bool {
	var userErr *domainError.UserDomainError
	if errors.As(err, &userErr) {
		return userErr.Code() == "AUTHENTICATION_ERROR"
	}
	return false
}
