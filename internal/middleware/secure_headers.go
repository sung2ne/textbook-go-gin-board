package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type SecureConfig struct {
	HSTSMaxAge            int
	HSTSIncludeSubdomains bool
	FrameOption           string
	ContentSecurityPolicy string
	ReferrerPolicy        string
	Debug                 bool
}

func DefaultSecureConfig() SecureConfig {
	return SecureConfig{
		HSTSMaxAge:            31536000,
		HSTSIncludeSubdomains: true,
		FrameOption:           "DENY",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		Debug:                 false,
	}
}

func SecureHeaders(cfg SecureConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("X-Frame-Options", cfg.FrameOption)
		c.Header("Referrer-Policy", cfg.ReferrerPolicy)
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		if !cfg.Debug {
			hsts := fmt.Sprintf("max-age=%d", cfg.HSTSMaxAge)
			if cfg.HSTSIncludeSubdomains {
				hsts += "; includeSubDomains"
			}
			c.Header("Strict-Transport-Security", hsts)
		}

		if cfg.ContentSecurityPolicy != "" {
			c.Header("Content-Security-Policy", cfg.ContentSecurityPolicy)
		}

		c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
		c.Header("Pragma", "no-cache")

		c.Next()
	}
}
