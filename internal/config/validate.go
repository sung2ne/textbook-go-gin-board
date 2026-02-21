package config

import (
	"errors"
	"fmt"
)

func (c *Config) Validate() error {
	var errs []error

	if c.Database.Host == "" {
		errs = append(errs, errors.New("database host is required"))
	}

	if c.JWT.SecretKey == "" {
		errs = append(errs, errors.New("JWT secret key is required"))
	}

	if c.Server.Port < 1 || c.Server.Port > 65535 {
		errs = append(errs, fmt.Errorf("server port must be between 1 and 65535, got %d", c.Server.Port))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
