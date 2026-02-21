package repository

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type PostRepositoryTestSuite struct {
    suite.Suite
    db   *gorm.DB
    repo *PostRepository
}

func (s *PostRepositoryTestSuite) SetupSuite() {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    s.Require().NoError(err)

    err = db.AutoMigrate(&Post{})
    s.Require().NoError(err)

    s.db = db
}

func (s *PostRepositoryTestSuite) SetupTest() {
    // 각 테스트 전에 테이블 비우기
    s.db.Exec("DELETE FROM posts")
    s.repo = NewPostRepository(s.db)
}

func (s *PostRepositoryTestSuite) TearDownSuite() {
    sqlDB, _ := s.db.DB()
    sqlDB.Close()
}

func (s *PostRepositoryTestSuite) TestCreate() {
    post := &Post{Title: "제목", Content: "내용"}

    err := s.repo.Create(post)

    s.NoError(err)
    s.NotZero(post.ID)
    s.NotZero(post.CreatedAt)
}

func (s *PostRepositoryTestSuite) TestFindByID() {
    // 준비
    post := &Post{Title: "제목", Content: "내용"}
    s.db.Create(post)

    // 실행
    found, err := s.repo.FindByID(post.ID)

    // 검증
    s.NoError(err)
    s.Equal(post.ID, found.ID)
    s.Equal("제목", found.Title)
}

func (s *PostRepositoryTestSuite) TestFindByID_NotFound() {
    _, err := s.repo.FindByID(999)

    s.Error(err)
    s.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (s *PostRepositoryTestSuite) TestFindAll() {
    // 10개 생성
    for i := 1; i <= 10; i++ {
        s.db.Create(&Post{
            Title:   fmt.Sprintf("게시글 %d", i),
            Content: "내용",
        })
    }

    // 첫 페이지 조회
    posts, err := s.repo.FindAll(1, 5)

    s.NoError(err)
    s.Len(posts, 5)
}

func (s *PostRepositoryTestSuite) TestUpdate() {
    post := &Post{Title: "원래 제목", Content: "내용"}
    s.db.Create(post)

    post.Title = "변경된 제목"
    err := s.repo.Update(post)

    s.NoError(err)

    var found Post
    s.db.First(&found, post.ID)
    s.Equal("변경된 제목", found.Title)
}

func (s *PostRepositoryTestSuite) TestDelete() {
    post := &Post{Title: "삭제할 게시글", Content: "내용"}
    s.db.Create(post)

    err := s.repo.Delete(post.ID)

    s.NoError(err)

    var count int64
    s.db.Model(&Post{}).Count(&count)
    s.Equal(int64(0), count)
}

func TestPostRepositoryTestSuite(t *testing.T) {
    suite.Run(t, new(PostRepositoryTestSuite))
}
