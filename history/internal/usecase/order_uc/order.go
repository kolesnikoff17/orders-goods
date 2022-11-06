package order_uc

import (
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
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
	order, err := uc.r.GetPrices(ctx, order)
	switch {
	case errors.Is(err, entity.ErrNoID):
		return err
	case err != nil:
		return fmt.Errorf("order_uc - createorder: %w", err)
	}
	order.Sum = calculateSum(order)
	err = uc.r.Create(ctx, order)
	if err != nil {
		return fmt.Errorf("order_uc - createorder: %w", err)
	}
	return nil
}

// UpdateOrder update actual order info
func (uc *OrderUseCase) UpdateOrder(ctx context.Context, order entity.Order) error {
	order, err := uc.r.GetPrices(ctx, order)
	switch {
	case errors.Is(err, entity.ErrNoID):
		return err
	case err != nil:
		return fmt.Errorf("order_uc - createorder: %w", err)
	}
	order.Sum = calculateSum(order)
	err = uc.r.Update(ctx, order)
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

func calculateSum(order entity.Order) string {
	sum := decimal.NewFromInt(0)
	for _, g := range order.Goods {
		price, _ := decimal.NewFromString(g.Price)
		price = price.Mul(decimal.NewFromInt(int64(g.Amount)))
		sum = sum.Add(price)
	}
	return sum.String()
}
