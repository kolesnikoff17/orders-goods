package order_uc

import (
	"context"
	"history/internal/entity"
)

// Order is a model layer interface for CUD operations on orders
type Order interface {
	CreateOrder(ctx context.Context, order entity.Order) error
	UpdateOrder(ctx context.Context, order entity.Order) error
	DeleteOrder(ctx context.Context, id int) error
}

// OrderRepo is a repo layer interface for CUD operations on orders
type OrderRepo interface {
	Create(ctx context.Context, order entity.Order) error
	Update(ctx context.Context, order entity.Order) error
	Archive(ctx context.Context, id int) error
}
