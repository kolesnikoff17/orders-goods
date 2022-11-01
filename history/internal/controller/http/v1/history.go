package v1

import (
	"common.local/pkg/logger"
	"errors"
	"github.com/gin-gonic/gin"
	"history/internal/entity"
	"history/internal/usecase/history_uc"
	"net/http"
	"strconv"
)

type goodRouter struct {
	uc history_uc.History
	l  logger.Interface
}

type emptyJSONResponse struct {
}

func newBalanceRoutes(handler *gin.RouterGroup, uc history_uc.History, l logger.Interface) {
	r := &goodRouter{
		uc: uc,
		l:  l,
	}

	handler.GET("/orders/:id", r.getByID)
	handler.GET("/orders/history/:id", r.getOrderHistory)
}

// @Summary     getByID
// @Description return order with given id
// @Tags  	    order
// @Accept      json
// @Produce     json
// @Param       id path string true "order id"
// @Success     200 {object} entity.Order
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /orders/{id} [get]
func (r *goodRouter) getByID(c *gin.Context) {
	idPath := c.Param("id")
	id, err := strconv.Atoi(idPath)
	if err != nil {
		r.l.Infof("error %s with id %d", err, id)
		errorResponse(c, http.StatusBadRequest, "Wrong id format")
		return
	}
	o, err := r.uc.GetOrderByID(c.Request.Context(), id)
	switch {
	case errors.Is(err, entity.ErrNoID):
		r.l.Infof("error %s with id %d", err, id)
		errorResponse(c, http.StatusBadRequest, "No such order")
		return
	case err != nil:
		r.l.Warnf("error %s with id %d", err, id)
		errorResponse(c, http.StatusInternalServerError, "Database error")
		return
	}
	c.JSON(http.StatusOK, o)
}

// @Summary     getOrderHistory
// @Description return order's state history with given id
// @Tags  	    order
// @Accept      json
// @Produce     json
// @Param       id path string true "order id"
// @Success     200 {array} entity.Order
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /orders/history/{id} [get]
func (r *goodRouter) getOrderHistory(c *gin.Context) {
	idPath := c.Param("id")
	id, err := strconv.Atoi(idPath)
	if err != nil {
		r.l.Infof("error %s with id %d", err, id)
		errorResponse(c, http.StatusBadRequest, "Wrong id format")
		return
	}
	h, err := r.uc.GetOrderHistory(c.Request.Context(), id)
	switch {
	case errors.Is(err, entity.ErrNoID):
		r.l.Infof("error %s with id %d", err, id)
		errorResponse(c, http.StatusBadRequest, "No such order")
		return
	case err != nil:
		r.l.Warnf("error %s with id %d", err, id)
		errorResponse(c, http.StatusInternalServerError, "Database error")
		return
	}
	c.JSON(http.StatusOK, h)
}
