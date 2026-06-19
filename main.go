package main

import (
    "net/http"
    "path/filepath"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // 업로드 파일 크기 제한 (8MB)
    r.MaxMultipartMemory = 8 << 20

    r.POST("/upload", func(c *gin.Context) {
        // 폼에서 파일 가져오기
        file, err := c.FormFile("file")
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // 파일 저장
        filename := filepath.Base(file.Filename)
        dst := "./uploads/" + filename

        if err := c.SaveUploadedFile(file, dst); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "filename": filename,
            "size":     file.Size,
        })
    })

    r.Run(":8080")
}
