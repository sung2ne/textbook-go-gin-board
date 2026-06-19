package main

import (
    "context"
    "fmt"
    "sync"
    "time"

    "golang.org/x/sync/errgroup"
)

type UserData struct {
    Profile  string
    Posts    []string
    Followers int
}

func fetchProfile(ctx context.Context, userID int) (string, error) {
    time.Sleep(100 * time.Millisecond)
    return fmt.Sprintf("User %d Profile", userID), nil
}

func fetchPosts(ctx context.Context, userID int) ([]string, error) {
    time.Sleep(150 * time.Millisecond)
    return []string{"Post 1", "Post 2"}, nil
}

func fetchFollowers(ctx context.Context, userID int) (int, error) {
    time.Sleep(80 * time.Millisecond)
    return 1234, nil
}

func fetchUserData(ctx context.Context, userID int) (*UserData, error) {
    g, ctx := errgroup.WithContext(ctx)

    var data UserData
    var mu sync.Mutex

    g.Go(func() error {
        profile, err := fetchProfile(ctx, userID)
        if err != nil {
            return err
        }
        mu.Lock()
        data.Profile = profile
        mu.Unlock()
        return nil
    })

    g.Go(func() error {
        posts, err := fetchPosts(ctx, userID)
        if err != nil {
            return err
        }
        mu.Lock()
        data.Posts = posts
        mu.Unlock()
        return nil
    })

    g.Go(func() error {
        followers, err := fetchFollowers(ctx, userID)
        if err != nil {
            return err
        }
        mu.Lock()
        data.Followers = followers
        mu.Unlock()
        return nil
    })

    if err := g.Wait(); err != nil {
        return nil, err
    }

    return &data, nil
}

func main() {
    ctx := context.Background()

    start := time.Now()
    data, err := fetchUserData(ctx, 123)
    elapsed := time.Since(start)

    if err != nil {
        fmt.Println("에러:", err)
        return
    }

    fmt.Printf("데이터: %+v\n", data)
    fmt.Printf("소요 시간: %v\n", elapsed) // 약 150ms (가장 느린 작업)
}
