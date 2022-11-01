package order_uc

import (
	"context"
	"fmt"
	"history/internal/entity"
)

// OrderUseCase implements Order interface and keeps OrderRepo
type OrderUseCase struct {
	r OrderRepo
}

// New is a constructor for OrderUseCase
func New(r OrderRepo) *OrderUseCase {
	return &OrderUseCase{
		r: r,
	}
}

// CreateOrder create new order in db
func (uc *OrderUseCase) CreateOrder(ctx context.Context, order entity.Order) error {
	err := uc.r.Create(ctx, order)
	if err != nil {
		return fmt.Errorf("order_uc - createorder: %w", err)
	}
	return nil
}

// UpdateOrder update actual order info
func (uc *OrderUseCase) UpdateOrder(ctx context.Context, order entity.Order) error {
	err := uc.r.Update(ctx, order)
	if err != nil {
		return fmt.Errorf("order_uc - updateorder: %w", err)
	}
	return nil
}

// DeleteOrder remove order from active pool
func (uc *OrderUseCase) DeleteOrder(ctx context.Context, id int) error {
	err := uc.r.Archive(ctx, id)
	if err != nil {
		return fmt.Errorf("order_uc - deleteorder: %w", err)
	}
	return nil
}
