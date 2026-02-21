package middleware

import (
    "crypto/md5"
    "encoding/hex"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

// ConditionalGet - 조건부 GET 지원
func ConditionalGet(getModTime func(c *gin.Context) time.Time) gin.HandlerFunc {
    return func(c *gin.Context) {
        modTime := getModTime(c)

        // ETag 생성
        etag := generateETag(c.Request.URL.Path, modTime)
        c.Header("ETag", etag)
        c.Header("Last-Modified", modTime.Format(http.TimeFormat))

        // If-None-Match 확인
        if match := c.GetHeader("If-None-Match"); match == etag {
            c.AbortWithStatus(http.StatusNotModified)
            return
        }

        // If-Modified-Since 확인
        if since := c.GetHeader("If-Modified-Since"); since != "" {
            t, err := time.Parse(http.TimeFormat, since)
            if err == nil && !modTime.After(t.Add(time.Second)) {
                c.AbortWithStatus(http.StatusNotModified)
                return
            }
        }

        c.Next()
    }
}

func generateETag(path string, modTime time.Time) string {
    data := path + modTime.String()
    hash := md5.Sum([]byte(data))
    return `"` + hex.EncodeToString(hash[:]) + `"`
}
