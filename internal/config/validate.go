package config

import (
    "errors"
    "fmt"
)

func (c *Config) Validate() error {
    var errs []error

    if c.DatabaseURL == "" {
        errs = append(errs, errors.New("DATABASE_URL is required"))
    }

    if c.JWTSecret == "" {
        errs = append(errs, errors.New("JWT_SECRET is required"))
    }

    if c.JWTSecret == "dev-secret" && !c.Debug {
        errs = append(errs, errors.New("JWT_SECRET must be changed in production"))
    }

    if c.Port < 1 || c.Port > 65535 {
        errs = append(errs, fmt.Errorf("PORT must be between 1 and 65535, got %d", c.Port))
    }

    if len(errs) > 0 {
        return errors.Join(errs...)
    }

    return nil
}
