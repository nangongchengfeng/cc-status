package handler

import (
	"context"
	"net/http"

	"cc-status/server/internal/model/dto"
	"cc-status/server/internal/service"

	"github.com/gin-gonic/gin"
)

// syncIngestor 抽象同步 handler 所需的最小业务能力。
type syncIngestor interface {
	Ingest(context.Context, dto.SyncRequest) (dto.SyncResult, error)
}

// SyncHandler 负责同步接口的参数绑定与响应序列化。
type SyncHandler struct {
	service syncIngestor
}

// NewSyncHandler 创建同步 handler。
func NewSyncHandler(syncService *service.SyncService) *SyncHandler {
	return &SyncHandler{service: syncService}
}

// HandleSync 兼容现有 client 的响应协议，不使用通用 data 包装。
func (handler *SyncHandler) HandleSync(ctx *gin.Context) {
	var request dto.SyncRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorData("INVALID_REQUEST", "请求体格式不合法"))
		return
	}

	result, err := handler.service.Ingest(ctx, request)
	if err != nil {
		if service.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, errorData("INVALID_REQUEST", err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorData("INTERNAL_ERROR", "同步写入失败"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":            0,
		"message":         "success",
		"accepted_count":  result.AcceptedCount,
		"duplicate_count": result.DuplicateCount,
	})
}
