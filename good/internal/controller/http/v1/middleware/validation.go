package mw

import (
  "github.com/gin-gonic/gin"
  "good/pkg/logger"
  "net/http"
)

type response struct {
  Msg string `json:"error"`
}

// ValidateJSONBody binds request body to map
func ValidateJSONBody(l logger.Interface) gin.HandlerFunc {
  return func(c *gin.Context) {
    var body map[string]interface{}
    err := c.ShouldBindJSON(&body)
    if err != nil {
      l.Infof("validation err: %s", err)
      c.AbortWithStatusJSON(http.StatusBadRequest, response{Msg: "Invalid request body format"})
      return
    }
    _, ok := body["name"]
    if !ok {
      l.Infof("no name field in request body", err)
      c.AbortWithStatusJSON(http.StatusBadRequest, response{Msg: "Name is a required field"})
      return
    }
    _, ok = body["category"]
    if !ok {
      l.Infof("no name field in request body", err)
      c.AbortWithStatusJSON(http.StatusBadRequest, response{Msg: "Category is a required field"})
      return
    }
    c.Set("jsonBody", body)
    c.Next()
  }
}

// GetJSONBody returns request body as a map
func GetJSONBody(c *gin.Context) map[string]interface{} {
  return c.MustGet("jsonBody").(map[string]interface{})
}
