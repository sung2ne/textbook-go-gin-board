package middleware

import (
    "crypto/md5"
    "encoding/hex"
    "bytes"

    "github.com/gin-gonic/gin"
)

// ETag - ETag 미들웨어
func ETag() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 응답을 버퍼에 캡처
        writer := &bodyWriter{
            ResponseWriter: c.Writer,
            body:           &bytes.Buffer{},
        }
        c.Writer = writer

        c.Next()

        // 응답 본문의 해시 계산
        body := writer.body.Bytes()
        if len(body) > 0 {
            hash := md5.Sum(body)
            etag := `"` + hex.EncodeToString(hash[:]) + `"`

            // If-None-Match 헤더 확인
            ifNoneMatch := c.GetHeader("If-None-Match")
            if ifNoneMatch == etag {
                c.Status(304)  // Not Modified
                return
            }

            c.Header("ETag", etag)
        }

        // 실제 응답 전송
        c.Writer = writer.ResponseWriter
        c.Writer.Write(body)
    }
}

type bodyWriter struct {
    gin.ResponseWriter
    body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
    return w.body.Write(b)
}
