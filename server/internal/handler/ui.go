//go:build !embed
// +build !embed

package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// 静态资源目录（相对于 server 目录）
const uiDistDir = "internal/handler/ui/dist"

// RegisterUIRoutes 注册 UI 相关路由
// 如果没有静态资源可提供，则不注册任何路由
func RegisterUIRoutes(router *gin.Engine) {
	// 检查本地文件系统是否有内容
	if !hasLocalFiles() {
		return
	}

	// 使用本地文件系统
	fileServer := http.FileServer(http.Dir(uiDistDir))

	// 先处理静态文件
	router.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") || c.Request.URL.Path == "/healthz" {
			c.Next()
			return
		}

		filePath := filepath.Join(uiDistDir, strings.TrimPrefix(c.Request.URL.Path, "/"))
		if _, err := os.Stat(filePath); err == nil {
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		c.Next()
	})

	// NoRoute 处理 SPA fallback
	router.NoRoute(func(c *gin.Context) {
		indexPath := filepath.Join(uiDistDir, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			http.ServeFile(c.Writer, c.Request, indexPath)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})
}

// hasLocalFiles 检查本地文件系统是否有内容
func hasLocalFiles() bool {
	info, err := os.Stat(uiDistDir)
	if err != nil || !info.IsDir() {
		return false
	}

	entries, err := os.ReadDir(uiDistDir)
	if err != nil || len(entries) == 0 {
		return false
	}

	return true
}
