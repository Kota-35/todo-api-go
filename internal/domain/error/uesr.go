package error

import "fmt"

type UserDomainError struct {
	code    string
	message string
	cause   error
}

func (e *UserDomainError) Error() string {
	if e.cause != nil {
		// FYI: fmt.Sprintf
		// 		任意の型と文字列をまとめて文字列(string型に)
		// [More: https://qiita.com/Sekky0905/items/5a65602dce83551184b3]
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

func (e *UserDomainError) Code() string {
	return e.code
}

func (e *UserDomainError) Type() string {
	return "USER_DOMAIN_ERROR"
}

func (e *UserDomainError) Unwrap() error {
	return e.cause
}

func NewDuplicateEmailError(email string) *UserDomainError {
	return &UserDomainError{
		code:    "DUPLICATE_EMAIL",
		message: fmt.Sprintf("このメールアドレスはすでに使用されています: %s", email),
	}
}

func NewInvalidUserDataError(message string) *UserDomainError {
	return &UserDomainError{
		code:    "INVALID_USER_DATA",
		message: message,
	}
}

func NewAuthenticationError(message string) *UserDomainError {
	return &UserDomainError{
		code:    "AUTHENTICATION_ERROR",
		message: message,
	}
}
