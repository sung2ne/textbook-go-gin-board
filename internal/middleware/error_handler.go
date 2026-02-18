package middleware

import (
	"github.com/gin-gonic/gin"
	"goboardapi/pkg/apperror"
	"goboardapi/pkg/problem"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		instance := c.Request.URL.Path

		var pd *problem.Detail

		if appErr, ok := apperror.AsAppError(err); ok {
			pd = problem.FromAppError(appErr, instance)
		} else {
			pd = problem.InternalError(instance)
		}

		c.Header("Content-Type", problem.ContentType)
		c.JSON(pd.Status, pd)
	}
}
