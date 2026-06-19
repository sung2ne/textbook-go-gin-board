package handler_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "goboardapi/internal/config"
    "goboardapi/internal/database"
    "goboardapi/internal/domain"
    "goboardapi/internal/handler"
    "goboardapi/internal/repository"
    "goboardapi/internal/router"
    "goboardapi/internal/service"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// PostHandlerSuite 테스트 스위트
type PostHandlerSuite struct {
    suite.Suite
    db     *gorm.DB
    router *gin.Engine
}

// SetupSuite 테스트 시작 전 1회 실행
func (s *PostHandlerSuite) SetupSuite() {
    gin.SetMode(gin.TestMode)

    // 테스트 DB 연결
    dsn := "host=localhost user=gouser password=gopassword dbname=godb_test port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    s.Require().NoError(err)
    s.db = db

    // 마이그레이션
    db.AutoMigrate(&domain.Post{}, &domain.Comment{})

    // 의존성 주입
    cfg := &config.Config{
        Pagination: config.PaginationConfig{
            DefaultSize: 10,
            MaxSize:     100,
        },
    }

    postRepo := repository.NewPostRepository(db)
    commentRepo := repository.NewCommentRepository(db)
    postService := service.NewPostService(postRepo, cfg)
    commentService := service.NewCommentService(commentRepo, postRepo)
    postHandler := handler.NewPostHandler(postService)
    commentHandler := handler.NewCommentHandler(commentService)

    r := router.NewRouter(postHandler, commentHandler)
    s.router = r.Setup()
}

// TearDownSuite 테스트 종료 후 1회 실행
func (s *PostHandlerSuite) TearDownSuite() {
    // 테스트 테이블 삭제
    s.db.Migrator().DropTable(&domain.Comment{}, &domain.Post{})
}

// SetupTest 각 테스트 전 실행
func (s *PostHandlerSuite) SetupTest() {
    // 테이블 초기화
    s.db.Exec("TRUNCATE TABLE comments RESTART IDENTITY CASCADE")
    s.db.Exec("TRUNCATE TABLE posts RESTART IDENTITY CASCADE")
}

func TestPostHandlerSuite(t *testing.T) {
    suite.Run(t, new(PostHandlerSuite))
}
