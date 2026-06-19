package main

import (
    "myapp/response"

    "github.com/gin-gonic/gin"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    r := gin.Default()

    // 사용자 목록
    r.GET("/users", func(c *gin.Context) {
        users := []User{
            {1, "Alice", "alice@example.com"},
            {2, "Bob", "bob@example.com"},
        }

        response.List(c, users, 1, 10, 2)
    })

    // 사용자 상세
    r.GET("/users/:id", func(c *gin.Context) {
        id := c.Param("id")

        if id == "0" {
            response.NotFound(c, "사용자를 찾을 수 없습니다")
            return
        }

        user := User{1, "Alice", "alice@example.com"}
        response.Success(c, user)
    })

    // 사용자 생성
    r.POST("/users", func(c *gin.Context) {
        var req struct {
            Name  string `json:"name" binding:"required"`
            Email string `json:"email" binding:"required,email"`
        }

        if err := c.ShouldBindJSON(&req); err != nil {
            response.BadRequest(c, err.Error())
            return
        }

        user := User{1, req.Name, req.Email}
        response.Created(c, user)
    })

    r.Run(":8080")
}
