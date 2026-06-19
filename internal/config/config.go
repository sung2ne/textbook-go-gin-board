package config

import (
    "time"

    "github.com/caarlos0/env/v10"
)

type Config struct {
    DatabaseURL    string        `env:"DATABASE_URL,required"`
    JWTSecret      string        `env:"JWT_SECRET" envDefault:"dev-secret"`
    Port           int           `env:"PORT" envDefault:"8080"`
    ReadTimeout    time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
    WriteTimeout   time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
    AllowedOrigins []string      `env:"ALLOWED_ORIGINS" envSeparator:"," envDefault:"http://localhost:3000"`
    Debug          bool          `env:"DEBUG" envDefault:"false"`
}

func Load() (*Config, error) {
    cfg := &Config{}
    if err := env.Parse(cfg); err != nil {
        return nil, err
    }
    return cfg, nil
}
