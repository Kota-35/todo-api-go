package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// RefreshTokenGenerator はリフレッシュトークンの生成を担当するドメインサービス
type RefreshTokenGenerator interface {
	Generate() (string, error)
}

type refreshTokenGenerator struct{}

func NewRefreshTokenGenerator() RefreshTokenGenerator {
	return &refreshTokenGenerator{}
}

// Generate は暗号学的に安全な乱数を使用してリフレッシュトークンを生成
func (g *refreshTokenGenerator) Generate() (string, error) {
	// 32バイトの乱数を生成
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("リフレッシュトークンの生成に失敗しました: %w", err)
	}

	// base64エンコードして返す
	return base64.URLEncoding.EncodeToString(bytes), nil
}
