package repository

import (
    "context"
    "testing"

    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "github.com/testcontainers/testcontainers-go/wait"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func setupPostgresContainer(t *testing.T) (*gorm.DB, func()) {
    ctx := context.Background()

    pgContainer, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:16-alpine"),
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2)),
    )
    require.NoError(t, err)

    connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
    require.NoError(t, err)

    db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
    require.NoError(t, err)

    err = db.AutoMigrate(&Post{}, &Comment{})
    require.NoError(t, err)

    cleanup := func() {
        pgContainer.Terminate(ctx)
    }

    return db, cleanup
}

func TestPostRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    db, cleanup := setupPostgresContainer(t)
    defer cleanup()

    repo := NewPostRepository(db)

    t.Run("Create and Find", func(t *testing.T) {
        post := &Post{Title: "제목", Content: "내용"}
        err := repo.Create(post)
        require.NoError(t, err)

        found, err := repo.FindByID(post.ID)
        require.NoError(t, err)
        assert.Equal(t, post.Title, found.Title)
    })

    t.Run("Update", func(t *testing.T) {
        post := &Post{Title: "원래 제목", Content: "내용"}
        repo.Create(post)

        post.Title = "변경된 제목"
        err := repo.Update(post)
        require.NoError(t, err)

        found, _ := repo.FindByID(post.ID)
        assert.Equal(t, "변경된 제목", found.Title)
    })

    t.Run("Delete", func(t *testing.T) {
        post := &Post{Title: "삭제할 게시글", Content: "내용"}
        repo.Create(post)

        err := repo.Delete(post.ID)
        require.NoError(t, err)

        _, err = repo.FindByID(post.ID)
        assert.Error(t, err)
    })
}
