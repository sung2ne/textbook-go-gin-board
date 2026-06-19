package main

import (
    "fmt"
    "log"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type Post struct {
    gorm.Model
    Title string `gorm:"size:200"`
    Tags  []Tag  `gorm:"many2many:post_tags"`
}

type Tag struct {
    gorm.Model
    Name  string `gorm:"size:50;uniqueIndex"`
    Posts []Post `gorm:"many2many:post_tags"`
}

func main() {
    dsn := "host=localhost user=gouser password=gopassword dbname=godb port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // 마이그레이션
    db.AutoMigrate(&Post{}, &Tag{})

    // 태그 먼저 생성
    goTag := Tag{Name: "Go"}
    tutorialTag := Tag{Name: "Tutorial"}
    db.Create(&goTag)
    db.Create(&tutorialTag)

    // 게시글과 태그 연결
    post := Post{
        Title: "Go GORM 튜토리얼",
        Tags:  []Tag{goTag, tutorialTag},
    }
    db.Create(&post)

    fmt.Printf("게시글: %s\n", post.Title)
    for _, tag := range post.Tags {
        fmt.Printf("- 태그: %s\n", tag.Name)
    }
}
