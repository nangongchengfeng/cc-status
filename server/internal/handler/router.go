package handler

import (
	"net/http"

	"cc-status/server/internal/middleware"

	"github.com/gin-gonic/gin"
)

// NewRouter 构建首版 server 路由骨架。
func NewRouter(authToken string) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, successData(gin.H{"status": "ok"}))
	})

	protected := router.Group("/api/v1")
	protected.Use(middleware.StaticTokenAuth(authToken))
	protected.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, successData(gin.H{"message": "pong"}))
	})

	return router
}
