package good_uc

import (
	"context"
	"history/internal/entity"
)

// Good is a model layer interface for CUD operations on goods
type Good interface {
	CreateGood(ctx context.Context, good entity.GoodInOrder) error
	UpdateGood(ctx context.Context, good entity.GoodInOrder) error
	DeleteGood(ctx context.Context, id string) error
}

// GoodRepo is a repo layer interface for CUD operations on goods
type GoodRepo interface {
	Create(ctx context.Context, good entity.GoodInOrder) error
	Update(ctx context.Context, good entity.GoodInOrder) error
	Archive(ctx context.Context, id string) error
}
