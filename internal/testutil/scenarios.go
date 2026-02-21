package testutil

import (
	"goboardapi/internal/domain"
	"gorm.io/gorm"
)

type Scenario struct {
	db       *gorm.DB
	Users    []*domain.User
	Posts    []*domain.Post
	Comments []*domain.Comment
}

func SetupBoardScenario(db *gorm.DB) *Scenario {
	s := &Scenario{db: db}

	admin := &domain.User{Username: "관리자", Email: "admin@example.com", Role: "admin"}
	user1 := &domain.User{Username: "사용자1", Email: "user1@example.com", Role: "user"}
	user2 := &domain.User{Username: "사용자2", Email: "user2@example.com", Role: "user"}

	db.Create(admin)
	db.Create(user1)
	db.Create(user2)
	s.Users = []*domain.User{admin, user1, user2}

	post1 := &domain.Post{Title: "공지사항", AuthorID: admin.ID}
	post2 := &domain.Post{Title: "일반 게시글", AuthorID: user1.ID}
	post3 := &domain.Post{Title: "임시 저장", AuthorID: user2.ID}

	db.Create(post1)
	db.Create(post2)
	db.Create(post3)
	s.Posts = []*domain.Post{post1, post2, post3}

	comment := &domain.Comment{PostID: post2.ID, AuthorID: user2.ID, Content: "좋은 글이네요"}
	db.Create(comment)
	s.Comments = []*domain.Comment{comment}

	return s
}

func (s *Scenario) GetAdmin() *domain.User {
	for _, u := range s.Users {
		if u.Role == "admin" {
			return u
		}
	}
	return nil
}

func (s *Scenario) GetPosts() []*domain.Post {
	return s.Posts
}
