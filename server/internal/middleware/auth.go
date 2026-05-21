package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// StaticTokenAuth 校验静态 Bearer token。
func StaticTokenAuth(expectedToken string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedBody())
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != expectedToken {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedBody())
			return
		}

		ctx.Next()
	}
}

func unauthorizedBody() gin.H {
	return gin.H{
		"code":    "UNAUTHORIZED",
		"message": "缺少或无效的认证令牌",
	}
}
