package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

type Email struct {
	value string
}

func NewEmail(email string) (Email, error) {
	if email == "" {
		return Email{}, errors.New("メールアドレスは必須です")
	}

	// 正規化(小文字化, 前後の空白削除
	// NOTE: メールアドレスは小文字・大文字を区別しないため
	normalized := strings.ToLower(strings.TrimSpace(email))

	// Validation
	if !isValidEmail(normalized) {
		return Email{}, errors.New(("無効なメールアドレスの形式です"))
	}

	return Email{value: normalized}, nil
}

// 値の取得
func (e Email) String() string {
	return e.value
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
