package config

import "github.com/caarlos0/env/v11"

type Config struct {
	DbDsn      string `env:"DB_DSN"`      // Connection string для БД
	LogLvl     string `env:"LOG_LEVEL"`   // Уровень логирования
	EncryptKey string `env:"ENCRYPT_KEY"` // Ключ для шифрования данных
	JWTKey     string `env:"JWT_KEY"`
}

func New() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
