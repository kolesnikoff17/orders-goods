package v1

import (
  "github.com/gin-gonic/gin"
  "good/internal/usecase"
  "good/pkg/logger"

  swaggerFiles "github.com/swaggo/files"
  ginSwagger "github.com/swaggo/gin-swagger"
  // Swagger docs
  _ "good/docs"
)

// NewRouter is an entry point to controller layer: it sets up middleware for "/" route
// and groups routers by version
func NewRouter(handler *gin.Engine, b usecase.GoodUseCase, l logger.Interface) {
  handler.Use(gin.Logger())
  handler.Use(gin.Recovery())

  swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
  handler.GET("/good/swagger/*any", swaggerHandler)

  h := handler.Group("/v1")
  {
    newBalanceRoutes(h, b, l)
  }
}
