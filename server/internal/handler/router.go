package handler

import (
	"net/http"

	"cc-status/server/internal/middleware"

	"github.com/gin-gonic/gin"
)

// NewRouter 构建首版 server 路由骨架。
func NewRouter(
	authToken string,
	syncHandler gin.HandlerFunc,
	modelPricingHandler *ModelPricingHandler,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, successData(gin.H{"status": "ok"}))
	})

	protected := router.Group("/api/v1")
	protected.Use(middleware.StaticTokenAuth(authToken))
	if syncHandler != nil {
		protected.POST("/sync", syncHandler)
	}
	if modelPricingHandler != nil {
		protected.GET("/model-pricings", modelPricingHandler.HandleListModelPricings)
		protected.POST("/model-pricings", modelPricingHandler.HandleCreateModelPricing)
		protected.PUT("/model-pricings/:id", modelPricingHandler.HandleUpdateModelPricing)
	}
	protected.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, successData(gin.H{"message": "pong"}))
	})

	return router
}
