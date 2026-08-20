package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DbDsn         string        `env:"DB_DSN"`    // Connection string для БД
	LogLvl        string        `env:"LOG_LEVEL"` // Уровень логирования
	MaxBodySize   int64         `env:"MAX_BODY_SIZE" envDefault:"2048"`
	ServerAddress string        `env:"SERVER_ADDRESS"` // Адрес запуска HTTP-сервера
	ReadTimeout   time.Duration `env:"READ_TIMEOUT" envDefault:"5s"`
	WriteTimeout  time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout   time.Duration `env:"IDLE_TIMEOUT" envDefault:"30s"`
	JWTKey        string        `env:"JWT_KEY"`
	ExpToken      time.Duration `env:"TOKEN_EXPIRE" envDefault:"1h"`
}

func New() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
