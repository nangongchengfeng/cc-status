package handler

import (
	"context"
	"net/http"

	"cc-status/server/internal/model/dto"
	"cc-status/server/internal/service"

	"github.com/gin-gonic/gin"
)

// statsOverviewer 抽象总览统计接口所需的最小业务能力。
type statsOverviewer interface {
	Overview(context.Context) (dto.StatsOverviewResponse, error)
}

// StatsHandler 负责统计接口的参数绑定与响应序列化。
type StatsHandler struct {
	service statsOverviewer
}

// NewStatsHandler 创建统计 handler。
func NewStatsHandler(statsService *service.StatsService) *StatsHandler {
	return &StatsHandler{service: statsService}
}

// HandleOverview 返回总览统计结果。
func (handler *StatsHandler) HandleOverview(ctx *gin.Context) {
	overview, err := handler.service.Overview(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorData("INTERNAL_ERROR", "读取总览统计失败"))
		return
	}

	ctx.JSON(http.StatusOK, successData(gin.H{
		"total_tokens":   overview.TotalTokens,
		"total_cost_usd": overview.TotalCostUSD,
		"total_requests": overview.TotalRequests,
		"active_clients": overview.ActiveClients,
		"top_models":     overview.TopModels,
		"top_clients":    overview.TopClients,
	}))
}
