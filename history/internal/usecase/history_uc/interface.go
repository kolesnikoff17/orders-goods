package history_uc

import (
	"context"
	"history/internal/entity"
)

// History is a model layer interface for Read operations on orders
type History interface {
	GetOrderByID(ctx context.Context, id int) (entity.Order, error)
	GetOrderHistory(ctx context.Context, id int) ([]entity.Order, error)
}

// HistoryRepo is a repo layer interface for Read operations on orders
type HistoryRepo interface {
	GetByID(ctx context.Context, id int) (entity.Order, error)
	GetHistory(ctx context.Context, id int) ([]entity.Order, error)
}
