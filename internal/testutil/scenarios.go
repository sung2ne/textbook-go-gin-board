package testutil

type Scenario struct {
    db      *gorm.DB
    Users   []*User
    Posts   []*Post
    Comments []*Comment
}

func SetupBoardScenario(db *gorm.DB) *Scenario {
    s := &Scenario{db: db}

    // 사용자 생성
    admin := &User{Name: "관리자", Email: "admin@example.com", Role: "admin"}
    user1 := &User{Name: "사용자1", Email: "user1@example.com", Role: "user"}
    user2 := &User{Name: "사용자2", Email: "user2@example.com", Role: "user"}

    db.Create(admin)
    db.Create(user1)
    db.Create(user2)
    s.Users = []*User{admin, user1, user2}

    // 게시글 생성
    post1 := &Post{Title: "공지사항", AuthorID: admin.ID, Status: "published"}
    post2 := &Post{Title: "일반 게시글", AuthorID: user1.ID, Status: "published"}
    post3 := &Post{Title: "임시 저장", AuthorID: user2.ID, Status: "draft"}

    db.Create(post1)
    db.Create(post2)
    db.Create(post3)
    s.Posts = []*Post{post1, post2, post3}

    // 댓글 생성
    comment := &Comment{PostID: post2.ID, AuthorID: user2.ID, Content: "좋은 글이네요"}
    db.Create(comment)
    s.Comments = []*Comment{comment}

    return s
}

func (s *Scenario) GetAdmin() *User {
    for _, u := range s.Users {
        if u.Role == "admin" {
            return u
        }
    }
    return nil
}

func (s *Scenario) GetPublishedPosts() []*Post {
    var result []*Post
    for _, p := range s.Posts {
        if p.Status == "published" {
            result = append(result, p)
        }
    }
    return result
}
