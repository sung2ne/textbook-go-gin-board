package handler

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
)

func BenchmarkGetPost(b *testing.B) {
    gin.SetMode(gin.ReleaseMode)
    r := setupRouter()

    // 테스트 데이터 준비
    setupTestData()

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        req := httptest.NewRequest("GET", "/posts/1", nil)
        rec := httptest.NewRecorder()
        r.ServeHTTP(rec, req)
    }
}

func BenchmarkListPosts(b *testing.B) {
    gin.SetMode(gin.ReleaseMode)
    r := setupRouter()

    setupTestData()

    cases := []struct {
        name string
        url  string
    }{
        {"page1", "/posts?page=1&size=10"},
        {"page10", "/posts?page=10&size=10"},
        {"size100", "/posts?page=1&size=100"},
    }

    for _, tc := range cases {
        b.Run(tc.name, func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                req := httptest.NewRequest("GET", tc.url, nil)
                rec := httptest.NewRecorder()
                r.ServeHTTP(rec, req)
            }
        })
    }
}
