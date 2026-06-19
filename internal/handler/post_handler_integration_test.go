package handler

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"
)

type PostHandlerTestSuite struct {
    suite.Suite
    router  *gin.Engine
    handler *PostHandler
}

func (s *PostHandlerTestSuite) SetupSuite() {
    gin.SetMode(gin.TestMode)
}

func (s *PostHandlerTestSuite) SetupTest() {
    s.router = gin.New()
    repo := NewMockPostRepository()
    s.handler = NewPostHandler(repo)

    s.router.GET("/posts", s.handler.List)
    s.router.GET("/posts/:id", s.handler.Get)
    s.router.POST("/posts", s.handler.Create)
    s.router.PUT("/posts/:id", s.handler.Update)
    s.router.DELETE("/posts/:id", s.handler.Delete)
}

func (s *PostHandlerTestSuite) TestList() {
    req := httptest.NewRequest("GET", "/posts?page=1&size=10", nil)
    rec := httptest.NewRecorder()

    s.router.ServeHTTP(rec, req)

    s.Equal(http.StatusOK, rec.Code)

    var response ListResponse
    err := json.Unmarshal(rec.Body.Bytes(), &response)
    s.NoError(err)
    s.NotEmpty(response.Data)
}

func (s *PostHandlerTestSuite) TestCreate() {
    body := map[string]interface{}{
        "title":   "새 게시글",
        "content": "게시글 내용입니다.",
    }
    jsonBody, _ := json.Marshal(body)

    req := httptest.NewRequest("POST", "/posts", bytes.NewReader(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    rec := httptest.NewRecorder()

    s.router.ServeHTTP(rec, req)

    s.Equal(http.StatusCreated, rec.Code)

    var response PostResponse
    err := json.Unmarshal(rec.Body.Bytes(), &response)
    s.NoError(err)
    s.Equal("새 게시글", response.Title)
}

func (s *PostHandlerTestSuite) TestCreate_ValidationError() {
    tests := []struct {
        name string
        body map[string]interface{}
    }{
        {"empty title", map[string]interface{}{"title": "", "content": "내용"}},
        {"empty content", map[string]interface{}{"title": "제목", "content": ""}},
        {"both empty", map[string]interface{}{"title": "", "content": ""}},
    }

    for _, tt := range tests {
        s.Run(tt.name, func() {
            jsonBody, _ := json.Marshal(tt.body)
            req := httptest.NewRequest("POST", "/posts", bytes.NewReader(jsonBody))
            req.Header.Set("Content-Type", "application/json")
            rec := httptest.NewRecorder()

            s.router.ServeHTTP(rec, req)

            s.Equal(http.StatusBadRequest, rec.Code)
        })
    }
}

func TestPostHandlerTestSuite(t *testing.T) {
    suite.Run(t, new(PostHandlerTestSuite))
}
