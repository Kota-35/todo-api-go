package valueobject

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/hkdf"
)

type Pepper []byte

func (p Pepper) deriveKey() []byte {
	key := make([]byte, 32)
	kdf := hkdf.New(sha256.New, p, nil, []byte("refresh_token"))
	kdf.Read(key)

	return key
}

type RefreshToken struct {
	hashedValue string
}

func NewRefreshToken(plainRefreshToken string, pepper Pepper) (RefreshToken, error) {
	if err := validateRefreshToken(plainRefreshToken); err != nil {
		return RefreshToken{}, err
	}

	mac := hmac.New(sha256.New, pepper.deriveKey())
	mac.Write([]byte(plainRefreshToken))
	sum := mac.Sum(nil)

	return RefreshToken{
		hashedValue: hex.EncodeToString(sum),
	}, nil

}

func validateRefreshToken(refreshToken string) error {
	if refreshToken == "" {
		return errors.New("トークンが空文字です")
	}

	return nil
}

// 検証：平文トークン→HMAC→一定時間比較
func (t RefreshToken) Match(raw string, pepper Pepper) bool {
	mac := hmac.New(sha256.New, pepper.deriveKey())
	mac.Write([]byte(raw))
	want := mac.Sum(nil)
	got, _ := hex.DecodeString(t.hashedValue)

	return subtle.ConstantTimeCompare(want, got) == 1
}

func (t RefreshToken) Hash() string {
	return t.hashedValue
}
