package error

import "fmt"

type SessionDomainError struct {
	code    string
	message string
	cause   error
}

func (e *SessionDomainError) Error() string {
	if e.cause != nil {
		// FYI: fmt.Sprintf
		// 		任意の型と文字列をまとめて文字列(string型に)
		// [More: https://qiita.com/Sekky0905/items/5a65602dce83551184b3]
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

func (e *SessionDomainError) Code() string {
	return e.code
}

func (e *SessionDomainError) Type() string {
	return "SESSION_DOMAIN_ERROR"
}

func (e *SessionDomainError) Unwrap() error {
	return e.cause
}
