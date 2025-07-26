package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseUrl string `env:"DATABASE_URL"`
	JwtSecret   string `env:"JWT_SECRET"`
	Env         string `env:"ENV"`
	UseHTTPS    bool   `env:"USE_HTTPS"     envDefault:"false"`
	SslCertPath string `env:"SSL_CERT_PATH"`
	SslKeyPath  string `env:"SSL_KEY_PATH"`
}

func LoadEnv() Config {
	cfg := Config{}
	env.Parse(&cfg)

	return cfg
}

// HTTPSかつ__Host-プレフィックスクッキーが使用可能かどうかを判定
func (c Config) ShouldUseSecureCookies() bool {
	return c.UseHTTPS
}
