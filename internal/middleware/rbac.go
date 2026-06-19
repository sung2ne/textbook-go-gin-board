
// RequireAnyPermission 주어진 권한 중 하나라도 있으면 통과
func RequireAnyPermission(permissions ...domain.Permission) gin.HandlerFunc {
    return func(c *gin.Context) {
        claims, ok := GetCurrentUser(c)
        if !ok {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "인증이 필요합니다",
            })
            return
        }

        userRole := domain.Role(claims.Role)
        for _, p := range permissions {
            if domain.HasPermission(userRole, p) {
                c.Next()
                return
            }
        }

        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
            "error": "권한이 없습니다",
        })
    }
}

// RequireAllPermissions 주어진 권한을 모두 가져야 통과
func RequireAllPermissions(permissions ...domain.Permission) gin.HandlerFunc {
    return func(c *gin.Context) {
        claims, ok := GetCurrentUser(c)
        if !ok {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "인증이 필요합니다",
            })
            return
        }

        userRole := domain.Role(claims.Role)
        for _, p := range permissions {
            if !domain.HasPermission(userRole, p) {
                c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                    "error":               "권한이 없습니다",
                    "missing_permission": string(p),
                })
                return
            }
        }

        c.Next()
    }
}
