package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JWTGenerator interface {
	GenerateAccessToken(uid string, ssID string, expiresAt time.Time) (string, error)
	VerifyAccessToken(tokenString string) (*JWTClaims, error)
}

type JWTClaims struct {
	UserID         string    `json:"user_id"`
	RefreshTokenID string    `json:"refresh_token_id"`
	ExpiresAt      time.Time `json:"expires_at"`
	Iat            time.Time `json:"iat"` // issued at - 発行時刻
	Jti            uuid.UUID `json:"jti"` // jwt id - トークン固有のID
}

type jwtGenerator struct {
	jwtSecret string
}

func NewJWTGenerator(jwtSecret string) JWTGenerator {
	return &jwtGenerator{
		jwtSecret: jwtSecret,
	}
}

func (g *jwtGenerator) GenerateAccessToken(uid string, rtID string, expiresAt time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":          uid,
		"refresh_token_id": rtID,
		"expires_at":       expiresAt.Unix(),
		"iat":              time.Now().Unix(),
		"jti":              uuid.New().String(),
	})

	tokenString, err := token.SignedString([]byte(g.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("[jwtGenerator]JWTの生成に失敗しました: %w", err)
	}

	return tokenString, nil
}

/*
@NOTE: 以下を参考にしました

- https://okkun-sh.hatenablog.com/entry/2023/07/24/011338
*/
func (g *jwtGenerator) VerifyAccessToken(tokenString string) (*JWTClaims, error) {
	// 署名の検証
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("[jwtGenerator]Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(g.jwtSecret), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("[jwtGenerator]トークンの解析に失敗しました: %w", err)
	}

	// クレームの検証
	if !token.Valid {
		return nil, fmt.Errorf("[jwtGenerator]トークンが無効です")
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("[jwtGenerator]クレームの取得に失敗しました")
	}

	// クレームの詳細検証
	userID, ok := (*claims)["user_id"].(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("[jwtGenerator]user_idクレームが無効です")
	}

	refreshTokenID, ok := (*claims)["refresh_token_id"].(string)
	if !ok || refreshTokenID == "" {
		return nil, fmt.Errorf("[jwtGenerator]refresh_token_idクレームが無効です")
	}

	expiresAtFloat, ok := (*claims)["expires_at"].(float64)
	if !ok {
		return nil, fmt.Errorf("[jwtGenerator]expires_atクレームが無効です")
	}
	expiresAt := int64(expiresAtFloat)

	iatFloat, ok := (*claims)["iat"].(float64)
	if !ok {
		return nil, fmt.Errorf("[jwtGenerator]iatクレームが無効です")
	}
	iat := int64(iatFloat)

	jtiStr, ok := (*claims)["jti"].(string)
	if !ok {
		return nil, fmt.Errorf("[jwtGenerator]jtiクレームが無効です")
	}

	jtiUUID, err := uuid.Parse(jtiStr)
	if err != nil {
		return nil, fmt.Errorf("[jwtGenerator]jtiクレームのUUID変換に失敗しました: %w", err)
	}

	// 期限・失効確認チェック
	if time.Now().Unix() > expiresAt {
		return nil, fmt.Errorf("[jwtGenerator]トークンの有効期限が切れています")
	}

	return &JWTClaims{
		UserID:         userID,
		RefreshTokenID: refreshTokenID,
		ExpiresAt:      time.Unix(expiresAt, 0),
		Iat:            time.Unix(iat, 0),
		Jti:            jtiUUID,
	}, nil
}
