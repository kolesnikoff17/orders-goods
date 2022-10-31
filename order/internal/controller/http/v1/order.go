package v1

import (
	"common.local/pkg/logger"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"order/internal/controller/http/v1/middleware"
	"order/internal/entity"
	"order/internal/usecase"
	"strconv"
)

type orderRouter struct {
	uc usecase.OrderUseCase
	l  logger.Interface
}

type emptyJSONResponse struct {
}

func newBalanceRoutes(handler *gin.RouterGroup, uc usecase.OrderUseCase, l logger.Interface) {
	r := &orderRouter{
		uc: uc,
		l:  l,
	}

	handler.POST("/order", mw.ValidateJSONBody(r.l), r.createOrder)
	handler.PUT("/order/:id", mw.ValidateJSONBody(r.l), r.updateOrder)
	handler.DELETE("/order/:id", r.deleteOrder)
}

type orderPostResponse struct {
	ID int `json:"id"`
}

// @Summary     createOrder
// @Description create new order in repo
// @Tags  	    order
// @Accept      json
// @Produce     json
// @Param       request body mw.OrderRequestBody true "order data"
// @Success     201 {object} orderPostResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /order [post]
func (r *orderRouter) createOrder(c *gin.Context) {
	b := mw.GetJSONBody(c)
	l := make([]entity.GoodsInOrder, len(b.Goods))
	for i, v := range b.Goods {
		l[i].Amount = v.Amount
		l[i].GoodID = v.GoodID
	}
	id, err := r.uc.CreateNewOrder(c.Request.Context(), entity.Order{
		UserID: b.UserID,
		Goods:  l,
	})
	if err != nil {
		r.l.Warnf("error %s with body %v", err, b)
		errorResponse(c, http.StatusInternalServerError, "Database error")
		return
	}
	c.JSON(http.StatusCreated, orderPostResponse{ID: id})
}

// @Summary     updateOrder
// @Description update order data
// @Tags  	    order
// @Accept      json
// @Produce     json
// @Param       id path int true "order id"
// @Param       request body mw.OrderRequestBody true "order data"
// @Success     200 {object} emptyJSONResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /order/{id} [put]
func (r *orderRouter) updateOrder(c *gin.Context) {
	b := mw.GetJSONBody(c)
	l := make([]entity.GoodsInOrder, len(b.Goods))
	for i, v := range b.Goods {
		l[i].Amount = v.Amount
		l[i].GoodID = v.GoodID
	}
	idPath := c.Param("id")
	id, err := strconv.Atoi(idPath)
	if err != nil {
		r.l.Infof("error %s with id %v", err, id)
		errorResponse(c, http.StatusBadRequest, "Wrong id")
		return
	}
	err = r.uc.UpdateOrder(c.Request.Context(), entity.Order{
		ID:     id,
		UserID: b.UserID,
		Goods:  l,
	})
	switch {
	case errors.Is(err, entity.ErrNoID):
		r.l.Infof("error %s with id %v", err, id)
		errorResponse(c, http.StatusBadRequest, "No such id")
		return
	case err != nil:
		r.l.Warnf("error %s with id %v", err, id)
		errorResponse(c, http.StatusInternalServerError, "Database error")
		return
	}
	c.JSON(http.StatusOK, emptyJSONResponse{})
}

// @Summary     deleteOrder
// @Description delete order from db
// @Tags  	    order
// @Accept      json
// @Produce     json
// @Param       id path int true "order id"
// @Success     200 {object} emptyJSONResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /order/{id} [delete]
func (r *orderRouter) deleteOrder(c *gin.Context) {
	idPath := c.Param("id")
	id, err := strconv.Atoi(idPath)
	if err != nil {
		r.l.Infof("error %s with id %v", err, id)
		errorResponse(c, http.StatusBadRequest, "Wrong id")
		return
	}
	err = r.uc.DeleteOrder(c.Request.Context(), id)
	switch {
	case errors.Is(err, entity.ErrNoID):
		r.l.Infof("error %s with id %v", err, id)
		errorResponse(c, http.StatusBadRequest, "No such id")
		return
	case err != nil:
		r.l.Warnf("error %s with id %v", err, id)
		errorResponse(c, http.StatusInternalServerError, "Database error")
		return
	}
	c.JSON(http.StatusOK, emptyJSONResponse{})
}
