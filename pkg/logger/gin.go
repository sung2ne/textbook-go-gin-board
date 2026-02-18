package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func FromGin(c *gin.Context) zerolog.Logger {
	if l, exists := c.Get("logger"); exists {
		if logger, ok := l.(zerolog.Logger); ok {
			return logger
		}
	}
	return Logger
}
