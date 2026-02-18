package middleware

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
	"goboardapi/pkg/logger"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func BodyLogging(maxBodySize int) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.FromGin(c)

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		if len(requestBody) > 0 && len(requestBody) <= maxBodySize {
			log.Debug().
				RawJSON("request_body", requestBody).
				Msg("요청 본문")
		}

		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = blw

		c.Next()

		if blw.body.Len() > 0 && blw.body.Len() <= maxBodySize {
			log.Debug().
				RawJSON("response_body", blw.body.Bytes()).
				Msg("응답 본문")
		}
	}
}
