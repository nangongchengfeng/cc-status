package handler

import (
	"context"
	"net/http"

	"cc-status/server/internal/model/dto"
	"cc-status/server/internal/service"

	"github.com/gin-gonic/gin"
)

// statsReader 抽象统计接口所需的最小业务能力。
type statsReader interface {
	Overview(context.Context) (dto.StatsOverviewResponse, error)
	Trend(context.Context, dto.StatsTrendQuery) ([]dto.StatsTrendPoint, error)
	Dashboard(context.Context, dto.StatsDashboardQuery) (dto.StatsDashboardResponse, error)
}

// StatsHandler 负责统计接口的参数绑定与响应序列化。
type StatsHandler struct {
	service statsReader
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

// HandleTrend 返回趋势统计结果。
func (handler *StatsHandler) HandleTrend(ctx *gin.Context) {
	var query dto.StatsTrendQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorData("INVALID_REQUEST", "查询参数不合法"))
		return
	}

	trend, err := handler.service.Trend(ctx, query)
	if err != nil {
		if service.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, errorData("INVALID_REQUEST", err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorData("INTERNAL_ERROR", "读取趋势统计失败"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": trend})
}

// HandleDashboard 返回仪表盘统计结果。
func (handler *StatsHandler) HandleDashboard(ctx *gin.Context) {
	var query dto.StatsDashboardQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, errorData("INVALID_REQUEST", "查询参数不合法"))
		return
	}

	dashboard, err := handler.service.Dashboard(ctx, query)
	if err != nil {
		if service.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, errorData("INVALID_REQUEST", err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorData("INTERNAL_ERROR", "读取仪表盘统计失败"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": dashboard})
}
