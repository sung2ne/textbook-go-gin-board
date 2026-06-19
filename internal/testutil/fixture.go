package testutil

import (
    "encoding/json"
    "os"
    "path/filepath"
    "testing"

    "gorm.io/gorm"
)

type FixtureLoader struct {
    db      *gorm.DB
    baseDir string
}

func NewFixtureLoader(db *gorm.DB) *FixtureLoader {
    return &FixtureLoader{
        db:      db,
        baseDir: "testdata",
    }
}

func (l *FixtureLoader) LoadUsers(t *testing.T) []*User {
    var users []*User
    l.loadJSON(t, "users.json", &users)

    for _, u := range users {
        l.db.Create(u)
    }

    return users
}

func (l *FixtureLoader) LoadPosts(t *testing.T) []*Post {
    var posts []*Post
    l.loadJSON(t, "posts.json", &posts)

    for _, p := range posts {
        l.db.Create(p)
    }

    return posts
}

func (l *FixtureLoader) loadJSON(t *testing.T, filename string, v interface{}) {
    path := filepath.Join(l.baseDir, filename)
    data, err := os.ReadFile(path)
    if err != nil {
        t.Fatalf("failed to read fixture file: %v", err)
    }

    if err := json.Unmarshal(data, v); err != nil {
        t.Fatalf("failed to parse fixture file: %v", err)
    }
}
