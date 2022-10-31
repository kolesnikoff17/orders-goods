package usecase

import (
	"context"
	"errors"
	"fmt"
	"order/internal/entity"
	"order/internal/usecase/kafka"
)

// Order keeps repo and notifier interfaces and implements business-logic
type Order struct {
	r OrderRepository
	k Notifier
}

// New is a constructor for Order
func New(r OrderRepository, k Notifier) *Order {
	return &Order{
		r: r,
		k: k,
	}
}

// CreateNewOrder create new order and return its id
func (uc *Order) CreateNewOrder(ctx context.Context, order entity.Order) (int, error) {
	id, err := uc.r.Create(ctx, order)
	if err != nil {
		return 0, fmt.Errorf("usecase - createorder: %w", err)
	}
	order, err = uc.r.GetByID(ctx, order.ID)
	if err != nil {
		return 0, fmt.Errorf("usecase - createorder: %w", err)
	}
	err = uc.k.Notify(kafka.Message{Action: kafka.Create, Data: kafka.Order{ID: id, Data: order}})
	if err != nil {
		return 0, fmt.Errorf("usecase - createorder: %w", err)
	}
	return id, nil
}

// UpdateOrder updates order with given id
func (uc *Order) UpdateOrder(ctx context.Context, order entity.Order) error {
	_, err := uc.r.GetByID(ctx, order.ID)
	switch {
	case errors.Is(err, entity.ErrNoID):
		return err
	case err != nil:
		return fmt.Errorf("usecase - updateorder: %w", err)
	}
	err = uc.r.Update(ctx, order)
	if err != nil {
		return fmt.Errorf("usecase - updateorder: %w", err)
	}
	order, err = uc.r.GetByID(ctx, order.ID)
	if err != nil {
		return fmt.Errorf("usecase - updateorder: %w", err)
	}
	err = uc.k.Notify(kafka.Message{Action: kafka.Update, Data: kafka.Order{ID: order.ID, Data: order}})
	return nil
}

// DeleteOrder deletes order with given id
func (uc *Order) DeleteOrder(ctx context.Context, id int) error {
	_, err := uc.r.GetByID(ctx, id)
	switch {
	case errors.Is(err, entity.ErrNoID):
		return err
	case err != nil:
		return fmt.Errorf("usecase - deleteorder: %w", err)
	}
	err = uc.r.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("usecase - deleteorder: %w", err)
	}
	err = uc.k.Notify(kafka.Message{Action: kafka.Delete, Data: kafka.Order{ID: id}})
	return nil
}
