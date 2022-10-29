package usecase

import (
  "context"
  "good/internal/entity"
)

// GoodUseCase is an interface for model layer
type GoodUseCase interface {
  NewGood(ctx context.Context, good entity.Good) (string, error)
  UpdateGood(ctx context.Context, good entity.Good) error
  DeleteGood(ctx context.Context, id string) error
}

// GoodRepo is an interface for repo layer
type GoodRepo interface {
  GetByID(ctx context.Context, id string) (entity.Good, error)
  Create(ctx context.Context, good entity.Good) (string, error)
  Update(ctx context.Context, good entity.Good) error
  Delete(ctx context.Context, id string) error
}
