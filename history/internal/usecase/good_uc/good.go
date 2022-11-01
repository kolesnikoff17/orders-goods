package good_uc

import (
	"context"
	"fmt"
	"history/internal/entity"
)

// GoodUseCase implements Order interface and keeps OrderRepo
type GoodUseCase struct {
	r GoodRepo
}

// New is a constructor for GoodUseCase
func New(r GoodRepo) *GoodUseCase {
	return &GoodUseCase{
		r: r,
	}
}

// CreateGood create new good in db
func (uc *GoodUseCase) CreateGood(ctx context.Context, good entity.GoodInOrder) error {
	err := uc.r.Create(ctx, good)
	if err != nil {
		return fmt.Errorf("good_uc - creategood: %w", err)
	}
	return nil
}

// UpdateGood update actual good info
func (uc *GoodUseCase) UpdateGood(ctx context.Context, good entity.GoodInOrder) error {
	err := uc.r.Update(ctx, good)
	if err != nil {
		return fmt.Errorf("good_uc - updategood: %w", err)
	}
	return nil
}

// DeleteGood remove good from active pool
func (uc *GoodUseCase) DeleteGood(ctx context.Context, id string) error {
	err := uc.r.Archive(ctx, id)
	if err != nil {
		return fmt.Errorf("good_uc - deletegood: %w", err)
	}
	return nil
}
