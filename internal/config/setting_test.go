package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_LoadEnv(t *testing.T) {
	t.Run("環境変数テスト", func(t *testing.T) {
		cfg := LoadEnv()

		assert.Equal(t, "development", cfg.Env)
	})
}
