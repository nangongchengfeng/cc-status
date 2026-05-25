package handler

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// 静态资源目录（相对于 server 目录）
const uiDistDir = "internal/ui/dist"

// RegisterUIRoutes 注册 UI 相关路由
// 如果没有静态资源可提供，则不注册任何路由
func RegisterUIRoutes(router *gin.Engine) {
	// 检查静态资源目录是否存在且有内容
	if !hasStaticFiles() {
		return
	}

	// 创建文件服务器
	fs := http.FileServer(http.Dir(uiDistDir))

	// 先处理静态文件：所有以 /assets/、/static/ 开头或带扩展名的请求
	router.Use(func(c *gin.Context) {
		// 先检查是否是 API 或健康检查请求，直接跳过
		if strings.HasPrefix(c.Request.URL.Path, "/api/") || c.Request.URL.Path == "/healthz" {
			c.Next()
			return
		}

		// 尝试作为静态文件处理
		filePath := filepath.Join(uiDistDir, c.Request.URL.Path)
		if _, err := os.Stat(filePath); err == nil {
			// 文件存在，直接 serve
			fs.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		// 继续下一个中间件
		c.Next()
	})

	// NoRoute 处理 SPA fallback：所有未匹配的路由都返回 index.html
	router.NoRoute(func(c *gin.Context) {
		indexPath := filepath.Join(uiDistDir, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			http.ServeFile(c.Writer, c.Request, indexPath)
			return
		}

		// index.html 也不存在，返回 404
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})
}

// hasStaticFiles 检查是否有静态资源文件
func hasStaticFiles() bool {
	// 检查目录是否存在
	info, err := os.Stat(uiDistDir)
	if err != nil || !info.IsDir() {
		return false
	}

	// 检查目录是否为空
	entries, err := os.ReadDir(uiDistDir)
	if err != nil || len(entries) == 0 {
		return false
	}

	return true
}

// 供后续 embed 使用的变量占位
var (
	embeddedFS fs.FS
)

