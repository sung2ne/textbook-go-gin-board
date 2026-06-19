package handler

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
    // 요청 생성
    req := httptest.NewRequest("GET", "/health", nil)

    // 응답 레코더 생성
    rec := httptest.NewRecorder()

    // 핸들러 실행
    HealthHandler(rec, req)

    // 응답 검증
    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Contains(t, rec.Body.String(), "ok")
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"ok"}`))
}
