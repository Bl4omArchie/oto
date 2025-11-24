package oto

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	PostgresqlDsn string `env:"POSTGRESQL_DSN,required"`
	TemporalHost  string `env:"TEMPORAL_HOST" envDefault:"localhost:7233"`
	TemporalWebUI string `env:"TEMPORAL_WEBUI_HOST" envDefault:"localhost:8080"`
}

func LoadOptionsFromEnv() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
