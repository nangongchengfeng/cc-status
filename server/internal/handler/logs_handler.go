package handler

import (
	"context"
	"net/http"

	"cc-status/server/internal/model/dto"
	"cc-status/server/internal/service"

	"github.com/gin-gonic/gin"
)

// logsLister 抽象日志查询接口所需的最小业务能力。
type logsLister interface {
	List(context.Context, dto.LogsQuery) (dto.LogsResponse, error)
}

// LogsHandler 负责原始日志查询接口的参数绑定与响应序列化。
type LogsHandler struct {
	service logsLister
}

// NewLogsHandler 创建日志查询 handler。
func NewLogsHandler(logsService *service.LogsService) *LogsHandler {
	return &LogsHandler{service: logsService}
}

// HandleListLogs 返回原始日志分页结果。
func (handler *LogsHandler) HandleListLogs(ctx *gin.Context) {
	var query dto.LogsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorData("INVALID_REQUEST", "查询参数不合法"))
		return
	}

	result, err := handler.service.List(ctx, query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorData("INTERNAL_ERROR", "读取原始日志失败"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":   result.Data,
		"total":  result.Total,
		"offset": result.Offset,
		"limit":  result.Limit,
		"page":   result.Page,
	})
}
