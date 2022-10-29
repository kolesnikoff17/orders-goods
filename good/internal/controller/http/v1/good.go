package v1

import (
  "errors"
  "github.com/gin-gonic/gin"
  mw "good/internal/controller/http/v1/middleware"
  "good/internal/entity"
  "good/internal/usecase"
  "good/pkg/logger"
  "net/http"
)

type goodRouter struct {
  uc usecase.GoodUseCase
  l  logger.Interface
}

type emptyJSONResponse struct {
}

func newBalanceRoutes(handler *gin.RouterGroup, uc usecase.GoodUseCase, l logger.Interface) {
  r := &goodRouter{
    uc: uc,
    l:  l,
  }

  handler.POST("/good", mw.ValidateJSONBody(r.l), r.createGood)
  handler.PUT("/good/:id", mw.ValidateJSONBody(r.l), r.updateGood)
  handler.DELETE("/good/:id", r.deleteGood)
}

type goodPostRequest struct {
  Name     string `json:"name" binding:"required" example:"Ice cream"`
  Category string `json:"category" binding:"required" example:"Food"`
}

type goodPostResponse struct {
  ID string `json:"id"`
}

// @Summary     createGood
// @Description creates new good in repo
// @Tags  	    good
// @Accept      json
// @Produce     json
// @Param       request body goodPostRequest true "name and category"
// @Success     201 {object} goodPostResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /order [post]
func (r *goodRouter) createGood(c *gin.Context) {
  b := mw.GetJSONBody[goodPostRequest](c)
  id, err := r.uc.NewGood(c.Request.Context(), entity.Good{Data: b})
  switch {
  case errors.Is(err, entity.GoodExists):
    r.l.Infof("error %s with body %v", err, b)
    errorResponse(c, http.StatusBadRequest, "Good already exists")
    return
  case err != nil:
    r.l.Warnf("error %s with body %v", err, b)
    errorResponse(c, http.StatusInternalServerError, "Database error")
    return
  }
  c.JSON(http.StatusCreated, goodPostResponse{ID: id})
}

type goodPutRequest struct {
  Name     string `json:"name" binding:"required" example:"Bicycle"`
  Category string `json:"category" binding:"required" example:"Transport"`
}

// @Summary     createGood
// @Description creates new good in repo
// @Tags  	    good
// @Accept      json
// @Produce     json
// @Param       id path string true "good id"
// @Param       request body goodPutRequest true "new good data"
// @Success     200 {object} emptyJSONResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /order/{id} [put]
func (r *goodRouter) updateGood(c *gin.Context) {
  b := mw.GetJSONBody[goodPutRequest](c)
  id := c.Param("id")
  err := r.uc.UpdateGood(c.Request.Context(), entity.Good{ID: id, Data: b})
  switch {
  case errors.Is(err, entity.ErrNoID):
    r.l.Infof("error %s with id %s", err, id)
    errorResponse(c, http.StatusBadRequest, "No such good")
    return
  case err != nil:
    r.l.Warnf("error %s with id %s", err, id)
    errorResponse(c, http.StatusInternalServerError, "Database error")
    return
  }
  c.JSON(http.StatusOK, emptyJSONResponse{})
}

// @Summary     deleteGood
// @Description deletes good
// @Tags  	    good
// @Accept      json
// @Produce     json
// @Param       id path string true "good id"
// @Success     200 {object} emptyJSONResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /order/{id} [put]
func (r *goodRouter) deleteGood(c *gin.Context) {
  id := c.Param("id")
  err := r.uc.DeleteGood(c.Request.Context(), id)
  switch {
  case errors.Is(err, entity.ErrNoID):
    r.l.Infof("error %s with id %s", err, id)
    errorResponse(c, http.StatusBadRequest, "No such good")
    return
  case err != nil:
    r.l.Warnf("error %s with id %s", err, id)
    errorResponse(c, http.StatusInternalServerError, "Database error")
    return
  }
  c.JSON(http.StatusOK, emptyJSONResponse{})
}
