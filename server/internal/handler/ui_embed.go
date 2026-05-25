//go:build embed
// +build embed

package handler

import (
	"embed"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// 嵌入的静态资源文件系统：相对于 ui_embed.go 所在目录
// 构建时 Makefile 会把 web/dist 复制到这里的 ui/dist 目录
//
//go:embed ui/dist/*
var embeddedFS embed.FS

// RegisterUIRoutes 注册 UI 相关路由
func RegisterUIRoutes(router *gin.Engine) {
	staticFS, ok := getEmbeddedFS()
	if !ok {
		return
	}

	fileServer := http.FileServer(http.FS(staticFS))

	// 先处理静态文件
	router.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") || c.Request.URL.Path == "/healthz" {
			c.Next()
			return
		}

		cleanPath := strings.TrimPrefix(c.Request.URL.Path, "/")
		if cleanPath == "" {
			cleanPath = "index.html"
		}

		if _, err := staticFS.Open(cleanPath); err == nil {
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		c.Next()
	})

	// NoRoute 处理 SPA fallback
	router.NoRoute(func(c *gin.Context) {
		if _, err := staticFS.Open("index.html"); err == nil {
			c.Request.URL.Path = "/index.html"
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})
}

// getEmbeddedFS 获取嵌入的文件系统
func getEmbeddedFS() (fs.FS, bool) {
	subFS, err := fs.Sub(embeddedFS, "ui/dist")
	if err != nil {
		return nil, false
	}

	// 检查是否有内容
	entries, err := subFS.ReadDir(".")
	if err != nil || len(entries) == 0 {
		return nil, false
	}

	return subFS, true
}
