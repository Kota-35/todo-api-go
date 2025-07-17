package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) (*string, error) {
	userID, exits := c.Get("user_id")

	if !exits {
		return nil, errors.New("ユーザーIDが見つかりません")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return nil, errors.New("ユーザーIDの型が正しくありません")
	}

	return &userIDStr, nil
}
