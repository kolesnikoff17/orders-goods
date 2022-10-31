package mw

import (
	"common.local/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
)

type response struct {
  Msg string `json:"error"`
}

// GoodRequestBody -.
type GoodRequestBody struct {
  Name       string                 `json:"name" binding:"required" example:"Ice cream"`
  Category   string                 `json:"category" binding:"required" example:"Food"`
  Price      string                 `json:"price" binding:"required" example:"200"`
  Additional map[string]interface{} `json:"additional" binding:"omitempty"`
}

// ValidateJSONBody binds request body to GoodRequestBody
func ValidateJSONBody(l logger.Interface) gin.HandlerFunc {
  return func(c *gin.Context) {
    var body GoodRequestBody
    err := c.ShouldBindJSON(&body)
    if err != nil {
      l.Infof("validation err: %s", err)
      c.AbortWithStatusJSON(http.StatusBadRequest, response{Msg: "Invalid request body format"})
      return
    }
    m, err := decimal.NewFromString(body.Price)
    if err != nil || !m.IsPositive() {
      l.Infof("validation err: %s", err)
      c.AbortWithStatusJSON(http.StatusBadRequest, response{Msg: "Wrong money format"})
      return
    }
    c.Set("jsonBody", body)
    c.Next()
  }
}

// GetJSONBody returns request body as a GoodRequestBody
func GetJSONBody(c *gin.Context) GoodRequestBody {
  return c.MustGet("jsonBody").(GoodRequestBody)
}
