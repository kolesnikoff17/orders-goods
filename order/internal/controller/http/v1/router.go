package v1

import (
	"common.local/pkg/logger"
	"github.com/gin-gonic/gin"
	"order/internal/usecase"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	// Swagger docs
	_ "order/docs"
)

// NewRouter is an entry point to controller layer: it sets up middleware for "/" route
// and groups routers by version
func NewRouter(handler *gin.Engine, b usecase.OrderUseCase, l logger.Interface) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/order/swagger/*any", swaggerHandler)

	h := handler.Group("/v1")
	{
		newBalanceRoutes(h, b, l)
	}
}
