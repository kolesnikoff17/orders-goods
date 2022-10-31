package mw

import (
	"common.local/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type response struct {
	Msg string `json:"error"`
}

// OrderRequestBody -.
type OrderRequestBody struct {
	UserID int         `json:"user_id" binding:"required,gte=1" example:"1"`
	Goods  []GoodsList `json:"goods" binding:"required"`
}

type GoodsList struct {
	GoodID string `json:"good_id" binding:"required"`
	Amount int    `json:"amount" binding:"required,gte=1"`
}

// ValidateJSONBody binds request body to OrderRequestBody
func ValidateJSONBody(l logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body OrderRequestBody
		err := c.ShouldBindJSON(&body)
		if err != nil {
			l.Infof("validation err: %s", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, response{Msg: "Invalid request body format"})
			return
		}
		c.Set("jsonBody", body)
		c.Next()
	}
}

// GetJSONBody returns request body as a OrderRequestBody
func GetJSONBody(c *gin.Context) OrderRequestBody {
	return c.MustGet("jsonBody").(OrderRequestBody)
}
