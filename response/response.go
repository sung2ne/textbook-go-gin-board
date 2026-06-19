package response

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// 성공 응답
func Success(c *gin.Context, data any) {
    c.JSON(http.StatusOK, Response{
        Success: true,
        Data:    data,
    })
}

// 생성 성공 응답
func Created(c *gin.Context, data any) {
    c.JSON(http.StatusCreated, Response{
        Success: true,
        Data:    data,
    })
}

// 목록 응답
func List(c *gin.Context, data any, page, perPage, total int) {
    c.JSON(http.StatusOK, Response{
        Success: true,
        Data:    data,
        Meta: &Meta{
            Page:    page,
            PerPage: perPage,
            Total:   total,
        },
    })
}

// 에러 응답
func Fail(c *gin.Context, status int, code, message string) {
    c.AbortWithStatusJSON(status, Response{
        Success: false,
        Error: &Error{
            Code:    code,
            Message: message,
        },
    })
}

// 400 Bad Request
func BadRequest(c *gin.Context, message string) {
    Fail(c, http.StatusBadRequest, "BAD_REQUEST", message)
}

// 401 Unauthorized
func Unauthorized(c *gin.Context, message string) {
    Fail(c, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

// 403 Forbidden
func Forbidden(c *gin.Context, message string) {
    Fail(c, http.StatusForbidden, "FORBIDDEN", message)
}

// 404 Not Found
func NotFound(c *gin.Context, message string) {
    Fail(c, http.StatusNotFound, "NOT_FOUND", message)
}

// 500 Internal Server Error
func InternalError(c *gin.Context, message string) {
    Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", message)
}
