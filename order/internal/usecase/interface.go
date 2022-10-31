package usecase

import (
	"context"
	"order/internal/entity"
	"order/internal/usecase/kafka"
)

// OrderUseCase is a model layer interface
type OrderUseCase interface {
	CreateNewOrder(ctx context.Context, order entity.Order) (int, error)
	UpdateOrder(ctx context.Context, order entity.Order) error
	DeleteOrder(ctx context.Context, id int) error
}

// OrderRepository is a repo layer interface
type OrderRepository interface {
	GetByID(ctx context.Context, id int) (entity.Order, error)
	Create(ctx context.Context, order entity.Order) (int, error)
	Update(ctx context.Context, order entity.Order) error
	Delete(ctx context.Context, id int) error
}

// Notifier is a listener interface
type Notifier interface {
	Notify(msg kafka.Message) error
}
