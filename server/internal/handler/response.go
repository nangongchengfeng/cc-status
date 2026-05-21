package handler

import "github.com/gin-gonic/gin"

func successData(payload gin.H) gin.H {
	return gin.H{"data": payload}
}

func errorData(code string, message string) gin.H {
	return gin.H{
		"code":    code,
		"message": message,
	}
}
