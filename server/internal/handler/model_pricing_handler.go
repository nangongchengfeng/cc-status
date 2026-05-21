package handler

import (
	"context"
	"net/http"

	"cc-status/server/internal/model/dto"
	"cc-status/server/internal/service"

	"github.com/gin-gonic/gin"
)

// modelPricingLister 抽象定价列表接口所需的最小业务能力。
type modelPricingLister interface {
	List(context.Context) ([]dto.ModelPricingResponse, error)
	Create(context.Context, dto.ModelPricingUpsertRequest) (dto.ModelPricingResponse, error)
	Update(context.Context, uint, dto.ModelPricingUpsertRequest) (dto.ModelPricingResponse, error)
}

// ModelPricingHandler 负责模型定价管理接口的参数绑定与响应序列化。
type ModelPricingHandler struct {
	service modelPricingLister
}

// NewModelPricingHandler 创建模型定价管理 handler。
func NewModelPricingHandler(modelPricingService *service.ModelPricingService) *ModelPricingHandler {
	return &ModelPricingHandler{service: modelPricingService}
}

// HandleListModelPricings 返回全部模型定价记录。
func (handler *ModelPricingHandler) HandleListModelPricings(ctx *gin.Context) {
	pricings, err := handler.service.List(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorData("INTERNAL_ERROR", "读取模型定价失败"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": pricings})
}

// HandleCreateModelPricing 创建一条新的模型定价。
func (handler *ModelPricingHandler) HandleCreateModelPricing(ctx *gin.Context) {
	var request dto.ModelPricingUpsertRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorData("INVALID_REQUEST", "请求体格式不合法"))
		return
	}

	pricing, err := handler.service.Create(ctx, request)
	if err != nil {
		if service.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, errorData("INVALID_REQUEST", err.Error()))
			return
		}
		if service.IsConflictError(err) {
			ctx.JSON(http.StatusConflict, errorData("CONFLICT", err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorData("INTERNAL_ERROR", "创建模型定价失败"))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": pricing})
}

// HandleUpdateModelPricing 按全量更新语义覆盖已有模型定价。
func (handler *ModelPricingHandler) HandleUpdateModelPricing(ctx *gin.Context) {
	var uri dto.ModelPricingURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorData("INVALID_REQUEST", "路径参数不合法"))
		return
	}

	var request dto.ModelPricingUpsertRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorData("INVALID_REQUEST", "请求体格式不合法"))
		return
	}

	pricing, err := handler.service.Update(ctx, uri.ID, request)
	if err != nil {
		if service.IsValidationError(err) {
			ctx.JSON(http.StatusBadRequest, errorData("INVALID_REQUEST", err.Error()))
			return
		}
		if service.IsConflictError(err) {
			ctx.JSON(http.StatusConflict, errorData("CONFLICT", err.Error()))
			return
		}
		if service.IsNotFoundError(err) {
			ctx.JSON(http.StatusNotFound, errorData("NOT_FOUND", err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorData("INTERNAL_ERROR", "更新模型定价失败"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": pricing})
}
