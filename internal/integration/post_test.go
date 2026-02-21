package integration

import (
    "testing"

    "github.com/stretchr/testify/suite"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "myapp/internal/handler"
    "myapp/internal/repository"
    "myapp/internal/service"
)

type PostIntegrationSuite struct {
    suite.Suite
    db      *gorm.DB
    tx      *gorm.DB
    handler *handler.PostHandler
}

func (s *PostIntegrationSuite) SetupSuite() {
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&repository.Post{}, &repository.User{})
    s.db = db
}

func (s *PostIntegrationSuite) SetupTest() {
    s.tx = s.db.Begin()

    repo := repository.NewPostRepository(s.tx)
    svc := service.NewPostService(repo)
    s.handler = handler.NewPostHandler(svc)
}

func (s *PostIntegrationSuite) TearDownTest() {
    s.tx.Rollback()
}

func (s *PostIntegrationSuite) TearDownSuite() {
    sqlDB, _ := s.db.DB()
    sqlDB.Close()
}

func (s *PostIntegrationSuite) TestCreateAndGet() {
    // Service를 통해 생성
    created, err := s.handler.Service.Create("제목", "내용", 1)
    s.Require().NoError(err)

    // 조회
    found, err := s.handler.Service.FindByID(created.ID)
    s.Require().NoError(err)

    s.Equal(created.ID, found.ID)
    s.Equal("제목", found.Title)
}

func (s *PostIntegrationSuite) TestList() {
    // 데이터 준비
    for i := 0; i < 15; i++ {
        s.tx.Create(&repository.Post{
            Title:    fmt.Sprintf("게시글 %d", i),
            Content:  "내용",
            AuthorID: 1,
        })
    }

    // 첫 페이지
    page1, _ := s.handler.Service.List(1, 10)
    s.Len(page1, 10)

    // 두 번째 페이지
    page2, _ := s.handler.Service.List(2, 10)
    s.Len(page2, 5)
}

func TestPostIntegrationSuite(t *testing.T) {
    suite.Run(t, new(PostIntegrationSuite))
}
