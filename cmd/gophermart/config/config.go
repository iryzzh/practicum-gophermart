package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"iryzzh/practicum-gophermart/internal/app/model"
	"sync"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS" envDefault:"127.0.0.1:8181"`
	DatabaseURI          string `env:"DATABASE_URI" envDefault:"postgresql://postgres:postgres@localhost/praktikum?sslmode=disable"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://127.0.0.1:8282"`
	MinPasswordLength    int    `env:"MIN_PASSWORD_LENGTH" envDefault:"6"`
	SessionKey           string `env:"SESSION_KEY" envDefault:"secret-key"`
}

var once sync.Once

func New() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	once.Do(func() {
		flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "run address")
		flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "database uri")
		flag.StringVar(&cfg.AccrualSystemAddress, "r", cfg.AccrualSystemAddress, "accrual system address")
		flag.IntVar(&cfg.MinPasswordLength, "pl", cfg.MinPasswordLength, "minimum password length")
		flag.StringVar(&cfg.SessionKey, "s", cfg.SessionKey, "session key")

		flag.Parse()
	})

	model.DefaultMinPasswordLength = cfg.MinPasswordLength

	return cfg, nil
}
