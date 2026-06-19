package main

import (
    "gin-tutorial/handler"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.GET("/", handler.Hello)

    r.Run(":8080")
}
