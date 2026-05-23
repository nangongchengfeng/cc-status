package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 中文请求日志中间件，记录请求方法、路径、状态码、耗时和客户端 IP。
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		log.Printf("[请求] %s %s | %d | %v | %s",
			c.Request.Method,
			path,
			statusCode,
			latency,
			clientIP,
		)
	}
}