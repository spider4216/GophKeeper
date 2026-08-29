package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DbDsn               string        `env:"DB_DSN"`                       // Connection string для БД
	LogLvl              string        `env:"LOG_LEVEL" envDefault:"debug"` // Уровень логирования
	EncryptKey          string        `env:"ENCRYPT_KEY"`                  // Ключ для шифрования данных
	JWTKey              string        `env:"JWT_KEY"`
	SrvHost             string        `env:"SERVER_HOST"` // Хост сервера
	TLSHandshakeTimeout time.Duration `env:"TLS_HANDSHAKE_TIMEOUT" envDefault:"3s"`
	RespHeaderTimeout   time.Duration `env:"RESPONSE_HEADER_TIMEOUT" envDefault:"3s"`
	DialerTimeout       time.Duration `env:"DEALER_TIMEOUT" envDefault:"5s"`
	CliTimeout          time.Duration `env:"CLIENT_TIMEOUT" envDefault:"10s"`
	CtxTimeout          time.Duration `env:"CTX_TIMEOUT" envDefault:"3s"`
	SyncChankSize       int           `env:"SYNC_CHANK_SIZE" envDefault:"100"`
	SyncLimit           int           `env:"SYNC_LIMIT" envDefault:"100"`
}

func New() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
