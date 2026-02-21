package handler

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func setupRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    return r
}

func TestGetPost(t *testing.T) {
    // 라우터 설정
    r := setupRouter()
    handler := NewPostHandler(mockRepo)
    r.GET("/posts/:id", handler.Get)

    // 요청 생성
    req := httptest.NewRequest("GET", "/posts/1", nil)
    rec := httptest.NewRecorder()

    // 실행
    r.ServeHTTP(rec, req)

    // 검증
    assert.Equal(t, http.StatusOK, rec.Code)

    var response PostResponse
    err := json.Unmarshal(rec.Body.Bytes(), &response)
    require.NoError(t, err)
    assert.Equal(t, uint(1), response.ID)
}
