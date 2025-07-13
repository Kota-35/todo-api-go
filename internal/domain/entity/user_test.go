package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_VerifyPassword(t *testing.T) {
	// テスト用のユーザーを作成
	user, err := NewUser(
		"test@example.com",
		"testUser",
		"collectPassword123!",
	)

	require.NoError(t, err)

	tests := []struct {
		name           string
		inputPassword  string
		expectedResult bool
		description    string
	}{
		{
			name:           "正しいパスワード",
			inputPassword:  "collectPassword123!",
			expectedResult: true,
			description:    "作成時と同じパスワードを入力した場合",
		},
		{
			name:           "間違ったパスワード",
			inputPassword:  "wrongPassword",
			expectedResult: false,
			description:    "異なるパスワードを入力した場合",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 実際の検証
			result := user.VerifyPassword(tt.inputPassword)
			assert.Equal(t, tt.expectedResult, result, tt.description)
		})
	}

}
