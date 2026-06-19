package config

import (
    "errors"
    "fmt"
    "net/url"
)

func (c *Config) Validate() error {
    var errs []error

    // 데이터베이스 URL 검증
    if _, err := url.Parse(c.Database.URL); err != nil {
        errs = append(errs, fmt.Errorf("invalid DATABASE_URL: %w", err))
    }

    // 프로덕션 필수 검증
    if c.IsProduction() {
        if c.JWT.Secret == "dev-secret" || len(c.JWT.Secret) < 32 {
            errs = append(errs, errors.New("JWT_SECRET must be at least 32 characters in production"))
        }

        if c.App.Debug {
            errs = append(errs, errors.New("DEBUG must be false in production"))
        }
    }

    // 커넥션 풀 검증
    if c.Database.MaxIdleConns > c.Database.MaxOpenConns {
        errs = append(errs, errors.New("DB_MAX_IDLE_CONNS cannot exceed DB_MAX_OPEN_CONNS"))
    }

    if len(errs) > 0 {
        return errors.Join(errs...)
    }

    return nil
}
