package middleware

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"goboardapi/pkg/sanitize"
)

func SanitizeInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.ContentType() != "application/json" {
			c.Next()
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Next()
			return
		}

		var data map[string]any
		if err := json.Unmarshal(body, &data); err != nil {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			c.Next()
			return
		}

		sanitizeMap(data)

		sanitized, _ := json.Marshal(data)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(sanitized))

		c.Next()
	}
}

func sanitizeMap(m map[string]any) {
	for k, v := range m {
		switch val := v.(type) {
		case string:
			m[k] = sanitize.Strict(val)
		case map[string]any:
			sanitizeMap(val)
		case []any:
			sanitizeSlice(val)
		}
	}
}

func sanitizeSlice(s []any) {
	for i, v := range s {
		switch val := v.(type) {
		case string:
			s[i] = sanitize.Strict(val)
		case map[string]any:
			sanitizeMap(val)
		}
	}
}
