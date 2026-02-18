package middleware

import (
	"errors"
	"net/http"
	"strings"

	"goboardapi/internal/auth"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	AuthorizationType   = "Bearer"
	ContextUserKey      = "user"
)

var (
	ErrMissingToken  = errors.New("missing authorization token")
	ErrInvalidFormat = errors.New("invalid authorization format")
)

// AuthMiddleware는 JWT 인증 미들웨어입니다.
func AuthMiddleware(tokenService *auth.TokenService, tokenStore auth.TokenStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "인증 토큰이 필요합니다",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != AuthorizationType {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "잘못된 인증 형식입니다",
			})
			return
		}

		tokenString := parts[1]

		claims, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			handleTokenError(c, err)
			return
		}

		// 블랙리스트 확인
		if tokenStore != nil {
			tokenID := claims.RegisteredClaims.ID
			if tokenID == "" {
				tokenID = auth.HashToken(tokenString)
			}

			blacklisted, err := tokenStore.IsBlacklisted(c.Request.Context(), tokenID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "토큰 검증 중 오류가 발생했습니다",
				})
				return
			}

			if blacklisted {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "토큰이 무효화되었습니다",
					"code":  "TOKEN_REVOKED",
				})
				return
			}
		}

		// 컨텍스트에 사용자 정보 저장
		c.Set(ContextUserKey, claims)

		ctx := SetUserToContext(c.Request.Context(), claims)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// OptionalAuthMiddleware는 선택적 인증 미들웨어입니다.
func OptionalAuthMiddleware(tokenService *auth.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)

		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != AuthorizationType {
			c.Next()
			return
		}

		tokenString := parts[1]

		claims, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			if errors.Is(err, auth.ErrExpiredToken) {
				c.Header("X-Token-Expired", "true")
			}
			c.Next()
			return
		}

		c.Set(ContextUserKey, claims)
		ctx := SetUserToContext(c.Request.Context(), claims)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func handleTokenError(c *gin.Context, err error) {
	if errors.Is(err, auth.ErrExpiredToken) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "토큰이 만료되었습니다",
			"code":  "TOKEN_EXPIRED",
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": "유효하지 않은 토큰입니다",
	})
}
