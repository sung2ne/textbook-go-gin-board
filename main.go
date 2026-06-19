package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // 정적 파일 서빙
    r.Static("/static", "./public")

    // 메인 페이지
    r.StaticFile("/", "./public/index.html")

    r.Run(":8080")
}
